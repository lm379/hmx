# 腾讯云数据万象转码回调配置说明

## 概述

系统已集成腾讯云数据万象（Cloud Infinite）工作流回调功能，用于自动处理视频转码、截图和字幕提取。

## 流程概览

- 前端通过预签名 URL 将文件上传到 `tmp/` 目录。
- 前端调用 `POST /api/v1/operas` 创建数据库记录（`video_path` 初始为 tmp 路径）。
- 腾讯云工作流执行转码、截图、字幕提取，并在完成后向本服务回调。
- 服务根据回调结果更新对应 `Opera` 记录；若回调先到，则将结果暂存到 Redis，待任务创建时拾取。

## 输出路径约定

- 视频（HLS）：public/video/{InputName}/{InputName}_hls.m3u8
- 封面（PNG）：public/cover/{InputName}/{InputName}.png
- 字幕（SRT）：public/video/{InputName}/{InputName}.srt

## 回调接口

端点：

```
POST /api/v1/callback/transcode
```

要点：

- 对外无需认证（仅供腾讯云回调）；
- 回调包含原始 `Object`（例如 `tmp/...`）以及各任务的 `Output.Object`；
- 仅更新转码成功的子任务字段，允许部分任务失败。

请求体（示例）简化为关键信息：

```json
{
  "EventName": "WorkflowFinish",
  "WorkflowExecution": {
    "Object": "tmp/video_20260103_123456.mp4",
    "State": "Success",
    "Tasks": {
      "Transcode": {"Code":"Success","Output":{"Object":"public/video/黄梅戏经典选段/黄梅戏经典选段_hls.m3u8"}},
      "Snapshot": {"Code":"Success","Output":{"Object":"public/cover/黄梅戏经典选段/黄梅戏经典选段.png"}},
      "Subtitle": {"Code":"Success","Output":{"Object":"public/video/黄梅戏经典选段/黄梅戏经典选段.srt"}}
    }
  }
}
```

响应示例：

成功：

```json
{
  "message":"Transcode callback processed successfully",
  "opera_id":123,
  "updates":{
    "video_path":"public/video/黄梅戏经典选段/黄梅戏经典选段_hls.m3u8",
    "avatar":"public/cover/黄梅戏经典选段/黄梅戏经典选段.png",
    "srt_path":"public/video/黄梅戏经典选段/黄梅戏经典选段.srt"
  }
}
```

失败（找不到记录）：

```json
{ "error":"Opera not found with original path: tmp/video_20260103_123456.mp4" }
```

## 腾讯云工作流配置建议

- 转码节点（Transcode）：输出 HLS，路径 `public/video/${InputName}/${InputName}_hls.m3u8`；
- 截图节点（Snapshot）：PNG，路径 `public/cover/${InputName}/${InputName}.png`，建议在第 5 秒或 10% 处截取；
- 字幕（Subtitle，可选）：SRT，路径 `public/video/${InputName}/${InputName}.srt`（若要浏览器直接使用请转为 VTT）。

> 浏览器对 SRT 支持有限，建议在 CDN 边缘将 SRT 转为 WebVTT（示例见下）。

> 本项目使用的EdgeOne的边缘函数如下：

```javascript
const VTT_HEADER = 'WEBVTT\n\n';

async function handleRequest(request) {
  const url = new URL(request.url);

  if (!url.pathname.endsWith('.srt')) {
    return fetch(request);
  }

  const originalResponse = await fetch(request);

  if (!originalResponse.ok) {
    return originalResponse;
  }

  const srtContent = await originalResponse.text();
  const vttContent = convertToWebvtt(srtContent);

  const response = new Response(vttContent);

  copyHeaders(originalResponse, response);
  
  response.headers.set('Content-Type', 'text/vtt; charset=utf-8');
  response.headers.set('Access-Control-Allow-Origin', '*');
  return response;
}

function convertToWebvtt(srtText) {
  const regex = /(\d{2}:\d{2}:\d{2}),(\d{3})\s+-->\s+(\d{2}:\d{2}:\d{2}),(\d{3})/g;
  const convertedBody = srtText.replace(regex, '$1.$2 --> $3.$4');
  return VTT_HEADER + convertedBody;
}

function copyHeaders(sourceRes, targetRes) {
  for (const [key, value] of sourceRes.headers) {
    if (key.toLowerCase() !== 'content-length') {
      targetRes.headers.set(key, value);
    }
  }
}

addEventListener('fetch', (event) => {
  event.respondWith(handleRequest(event.request));
});
```

## 回调时序与 Redis 临时结果（关键）

问题：回调可能在后端创建记录之前到达，导致找不到对应 Task/Opera。解决方案：使用 Redis 做“暂存箱”。

流程：

- 回调到达（消费者）：检查 `task_id`/原始 `Object` 对应的记录；若存在则直接更新；若不存在则把回调结果写入 `temp_result:{TaskID}`，设置短期过期（例如 5 分钟）；返回 200 表示已接收。
- 任务创建（生产者）：创建记录前先查询 `temp_result:{TaskID}`，若存在则用缓存结果直接将记录标记为完成并写入对应字段；若不存在则创建 PENDING 任务，等待回调。

优点：容错时序颠倒，解耦上传与转码回调；缺点：需要额外的 Redis 查询与过期策略。

实现要点：

- Redis Key 设计：`temp_result:{TaskID}`，Value 存原始回调 JSON 或结构化字段；
- 过期时间：建议 5 分钟，可根据系统吞吐调整；
- 回调接口收到临时写入时应记录日志方便排查；
- 在创建任务时读取缓存后应立即删除临时 Key，保证幂等性。

## 数据库字段说明

- `video_path` varchar(255)：回调后保存 HLS 路径，初始为 `tmp/...`；
- `avatar` varchar(255)：回调后保存封面路径；
- `srt_path` varchar(255)：回调后保存字幕路径（SRT），如需 VTT 则由 CDN 转换。

