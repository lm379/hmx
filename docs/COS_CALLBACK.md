# 腾讯云数据万象转码回调配置说明

## 概述

系统已集成腾讯云数据万象（Cloud Infinite）工作流回调功能，用于自动处理视频转码、截图和字幕提取。

## 工作流程

1. **视频上传**：用户通过预签名 URL 上传视频到 `/tmp` 目录
2. **创建记录**：前端调用 `POST /api/v1/operas` 创建 Opera 记录，`video_path` 为 tmp 路径
3. **自动转码**：腾讯云数据万象工作流自动触发，执行以下任务：
   - 转码为 HLS 格式 (m3u8)
   - 生成封面截图 (png)
   - 提取字幕 (srt)
4. **回调通知**：转码完成后，腾讯云调用回调接口更新 Opera 记录
5. **路径更新**：系统自动将临时路径更新为正式路径

## 转码后的文件路径规范

- **视频**：`public/video/{视频标题}/{视频标题}_hls.m3u8`
- **封面**：`public/cover/{视频标题}/{视频标题}.png`
- **字幕**：`public/video/{视频标题}/{视频标题}.srt`

## 回调接口

### 端点
```
POST /api/v1/callback/transcode
```

### 特性
- ✅ 无需认证（供腾讯云服务调用）
- ✅ 自动查找对应的 Opera 记录
- ✅ 批量更新视频、封面、字幕路径
- ✅ 支持部分任务失败（仅更新成功的任务）

### 请求体示例

```json
{
  "EventName": "WorkflowFinish",
  "WorkflowExecution": {
    "RunId": "i166ee19017dcf11eda8a5525400c******",
    "BucketId": "examplebucket-1250000000",
    "Object": "tmp/video_20260103_123456.mp4",
    "State": "Success",
    "CreateTime": "2026-01-03T12:00:00+08:00",
    "Tasks": {
      "Transcode": {
        "Code": "Success",
        "Message": "success",
        "JobId": "j29a82fea08ba11edb9335254008618d9",
        "TemplateId": "t182c0ca7d91ca40969a3fc97c5559091a",
        "TemplateName": "HLS-转码模板",
        "Output": {
          "Region": "ap-guangzhou",
          "Bucket": "examplebucket-1250000000",
          "Object": "public/video/黄梅戏经典选段/黄梅戏经典选段_hls.m3u8"
        }
      },
      "Snapshot": {
        "Code": "Success",
        "Message": "success",
        "JobId": "j29a82fea08ba11edb9335254008618d9",
        "TemplateId": "t182c0ca7d91ca40969a3fc97c5559091a",
        "TemplateName": "封面截图",
        "Output": {
          "Region": "ap-guangzhou",
          "Bucket": "examplebucket-1250000000",
          "Object": "public/cover/黄梅戏经典选段/黄梅戏经典选段.png"
        }
      },
      "Subtitle": {
        "Code": "Success",
        "Message": "success",
        "JobId": "j29a82fea08ba11edb9335254008618d9",
        "TemplateId": "t182c0ca7d91ca40969a3fc97c5559091a",
        "TemplateName": "字幕提取",
        "Output": {
          "Region": "ap-guangzhou",
          "Bucket": "examplebucket-1250000000",
          "Object": "public/video/黄梅戏经典选段/黄梅戏经典选段.srt"
        }
      }
    }
  }
}
```

### 响应示例

**成功：**
```json
{
  "message": "Transcode callback processed successfully",
  "opera_id": 123,
  "updates": {
    "video_path": "public/video/黄梅戏经典选段/黄梅戏经典选段_hls.m3u8",
    "avatar": "public/cover/黄梅戏经典选段/黄梅戏经典选段.png",
    "srt_path": "public/video/黄梅戏经典选段/黄梅戏经典选段.srt"
  }
}
```

**失败（找不到记录）：**
```json
{
  "error": "Opera not found with original path: tmp/video_20260103_123456.mp4"
}
```

## 腾讯云数据万象配置步骤

### 1. 创建工作流

在腾讯云控制台配置工作流，包含以下节点：

#### 转码节点（Transcode）
- **输出格式**：HLS (m3u8)
- **输出路径**：`public/video/${InputName}/${InputName}_hls.m3u8`
- **编码参数**：根据需求配置分辨率、码率等

#### 截图节点（Snapshot）
- **输出格式**：PNG
- **输出路径**：`public/cover/${InputName}/${InputName}.png`
- **截图时间**：建议第 5 秒或视频 10% 位置

#### 字幕提取节点（Subtitle，可选）
- **输出格式**：SRT
- **输出路径**：`public/video/${InputName}/${InputName}.srt`
- **语言**：中文

### 2. 配置回调地址

在工作流设置中配置回调 URL：
```
https://your-domain.com/api/v1/callback/transcode
```

### 3. 设置触发条件

- **触发目录**：`tmp/`
- **触发事件**：文件上传完成

### 4. 测试工作流

1. 上传测试视频到 `tmp/` 目录
2. 检查工作流执行日志
3. 验证回调是否成功调用
4. 确认 Opera 记录是否正确更新

## 数据库字段说明

Opera 表中相关字段：

| 字段 | 类型 | 说明 | 示例 |
|------|------|------|------|
| `video_path` | varchar(255) | 视频路径，初始为 tmp 路径，回调后更新为 HLS 路径 | `public/video/黄梅戏经典选段/黄梅戏经典选段_hls.m3u8` |
| `avatar` | varchar(255) | 封面路径，回调后自动填充 | `public/cover/黄梅戏经典选段/黄梅戏经典选段.png` |
| `srt_path` | varchar(255) | 字幕路径，回调后自动填充 | `public/video/黄梅戏经典选段/黄梅戏经典选段.srt` |

## 注意事项

1. **路径匹配**：回调时通过原始 tmp 路径查找 Opera 记录，确保创建记录时使用的路径与上传路径一致
2. **容错处理**：如果某个任务失败，仅更新成功的任务结果
3. **安全性**：建议配置 IP 白名单或签名验证（可选扩展）
4. **重试机制**：腾讯云默认会重试失败的回调，建议接口实现幂等性
5. **日志记录**：建议添加详细的日志记录，便于排查问题

## 故障排查

### 回调未收到
- 检查工作流回调 URL 配置是否正确
- 确认服务器防火墙/安全组允许腾讯云 IP 访问
- 查看腾讯云工作流执行日志

### 找不到 Opera 记录
- 确认创建 Opera 时的 `video_path` 与上传文件路径一致
- 检查数据库中是否存在对应的记录
- 查看回调请求中的 `Object` 字段是否正确

### 路径更新失败
- 检查数据库连接是否正常
- 查看服务端日志中的错误信息
- 确认字段类型和长度是否足够

## 后续优化建议

1. **签名验证**：添加腾讯云回调签名验证，提高安全性
2. **异步处理**：对于大量回调，可使用消息队列异步处理
3. **状态管理**：添加转码状态字段（pending, processing, completed, failed）
4. **通知机制**：转码完成后通知上传用户（WebSocket/邮件）
5. **重试策略**：实现自定义重试逻辑，处理临时失败
