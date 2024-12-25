---
title: 黄梅戏文化平台
description: 数据库与智能处理大作业：黄梅戏文化平台
date: 2024-12-20 01:03:12
top: 1
tags:
  - 数据库
---
# 黄梅戏文化平台

基于 `Node.js` + `Express` + `MySQL` + `Redis` 实现的黄梅戏文化展示平台

Github: [lm379/hmx: 黄梅戏文化平台](https://github.com/lm379/hmx)
Gitee: [lm379/hmx: 黄梅戏文化平台](https://gitee.com/lm379/hmx)

## 简介

后端使用 Node.js 配合 Express 框架开发，前端为 HTML + CSS + JavaScript 开发（使用模板）

数据库使用MySQL + Redis，其中MySQL负责主体业务逻辑，Redis负责处理验证码

邮箱验证码使用Node.js的nodemailer库，短信使用阿里云的短信接口

视频/图片文件使用腾讯云对象存储，后端授予前端权限允许前端通过 POST 方法直接上传用户头像/视频文件等到对象存储

为项目添加了智能处理，通过 OpenAI 开源的whisper 语音识别模型识别文字，再调用大模型API对视频进行总结，智能处理代码修改自下面代码fork的项目，智能处理代码Github链接: [lm379/video_summarize: video_summarize（视频总结：上传视频通过prompt对视频内容经行总结输出）](https://github.com/lm379/video_summarize)

## 项目结构
```
├── bin                // 程序入口目录
├── dataBase           // MySQL配置
├── node_modules       // 项目依赖 
├── public             // 前端代码
├── routes             // Express路由
├── utils              // 自定义依赖
├── app.js             // Express主程序
└── views              // Express视图
```

## 数据库表结构

1. User表

| 名称                | 类型        | 长度  | Not Null | 键     | 默认值               | 注释    | Key Order |
| ----------------- | --------- | --- | -------- | ----- | ----------------- | ----- | --------- |
| User_id           | int       |     | TRUE     | TRUE  |                   | 自增UID | ASC       |
| Username          | varchar   | 50  | TRUE     | FALSE |                   | 昵称    |           |
| Phone             | varchar   | 15  | TRUE     | FALSE |                   | 手机号   |           |
| Email             | varchar   | 255 | FALSE    | FALSE | NULL              | 邮箱    |           |
| Registration_time | timestamp |     | TRUE     | FALSE | CURRENT_TIMESTAMP | 注册时间  |           |
| Password          | varchar   | 255 | TRUE     | FALSE |                   | 密码    |           |
| Sex               | enum      |     | TRUE     | FALSE | 'Other'           | 性别    |           |
| Icon              | varchar   | 255 | FALSE    | FALSE | NULL              | 头像    |           |
| Role              | enum      |     | TRUE     | FALSE | 'User'            | 权限    |           |
| Update_time       | timestamp |     | TRUE     | FALSE | CURRENT_TIMESTAMP | 更新时间  |           |
| Status            | enum      |     | FALSE    | FALSE | 'Normal'          | 状态    |           |

2. Opera表

| 名称           | 类型        | 长度  | Not Null | 键     | 默认值               | 注释   | Key Order |
| ------------ | --------- | --- | -------- | ----- | ----------------- | ---- | --------- |
| Opera_id     | int       |     | TRUE     | TRUE  |                   | 自增ID | ASC       |
| Opera_title  | varchar   | 100 | TRUE     | FALSE |                   | 曲名   |           |
| Artist       | varchar   | 100 | FALSE    | FALSE | NULL              | 艺术家  |           |
| Release_date | datetime  |     | FALSE    | FALSE | NULL              | 发行时间 |           |
| Duration     | time      |     | FALSE    | FALSE | NULL              | 时长   |           |
| Music_path   | varchar   | 255 | FALSE    | FALSE | NULL              |      |           |
| Create_time  | timestamp |     | FALSE    | FALSE | CURRENT_TIMESTAMP | 上传时间 |           |
| Update_time  | timestamp |     | FALSE    | FALSE | CURRENT_TIMESTAMP | 更新时间 |           |
| Video_path   | varchar   | 255 | TRUE     | FALSE |                   |      |           |
| Description  | text      |     | FALSE    | FALSE | NULL              | 简介   |           |
| Avatar       | varchar   | 255 | FALSE    | FALSE | NULL              | 缩略图  |           |
|              |           |     |          |       |                   |      |           |
3. Comments表

| 名称           | 类型        | Not Null | 键     | 默认值               | 键长度 | Key Order |
| ------------ | --------- | -------- | ----- | ----------------- | --- | --------- |
| Comment_id   | int       | TRUE     | TRUE  |                   | 0   | ASC       |
| User_id      | int       | FALSE    | FALSE | NULL              |     |           |
| Opera_id     | int       | FALSE    | FALSE | NULL              |     |           |
| Comment_text | text      | TRUE     | FALSE |                   |     |           |
| Create_time  | timestamp | FALSE    | FALSE | CURRENT_TIMESTAMP |     |           |
| Update_time  | timestamp | FALSE    | FALSE | CURRENT_TIMESTAMP |     |           |
4. Playhistroy表

| 名称        | 类型        | Not Null | 键     | 默认值               | Key Order |
| --------- | --------- | -------- | ----- | ----------------- | --------- |
| Play_id   | int       | TRUE     | TRUE  |                   | ASC       |
| User_id   | int       | FALSE    | FALSE | NULL              |           |
| Opera_id  | int       | FALSE    | FALSE | NULL              |           |
| Play_time | timestamp | FALSE    | FALSE | CURRENT_TIMESTAMP |           |
5. Favorite表

| 名称           | 类型        | Not Null | 键     | 默认值               | Key Order |
|--------------|-----------|----------|-------|-------------------|-----------|
| Favorite_id  | int       | TRUE     | TRUE  |                   | ASC       |
| User_id      | int       | FALSE    | FALSE | NULL              |           |
| Opera_id     | int       | FALSE    | FALSE | NULL              |           |
| Created_time | timestamp | FALSE    | FALSE | CURRENT_TIMESTAMP |

## 使用方法

提前安装好 Node.js npm MySQL Redis 环境
>Node.js的代码理论上是跨平台的，但我只在Ubuntu下测试过

然后进入命令行，输入

```bash
npm install nodemon -g
```

拉取源码，安装依赖

```bash
git clone https://github.com/lm379/hmx
cd hmx
npm install
```

设置环境变量，并修改.env中的相关字段

```bash
mv .env.example .env
```
启动
```bash
npm start
```

## 部分代码示例

1. MySQL数据查询 (以下为示例代码，通过动态构建SQL语句实现)

```javascript
const getData = (reqParam) => {
    return new Promise((resolve, reject) => {
        let query = "SELECT * FROM table WHERE TRUE";
        let params = [];
        if (reqParam) {
            query += ` AND reqParam = ?`;
            params.push(reqParam);
        }

        if (params.length === 0) {
            return reject(new Error("参数错误"));
        }

        connection.query(query, params, (err, data) => {
            if (err) return reject(err);
            if (data.length === 0) return reject(resolve({ code: 1, msg: errmsg }));
            return resolve({ code: 0, data: data });
        });
    });
}
```

2. Redis验证码处理
参考了项目[TMDOG666/email-resend-demo](https://gitee.com/mbjdot/email-resend-demo)

```javascript
// 存储验证码
async function storeVerificationCode(account, verificationCode, verifyType, expiresIn = 300) {
    try {
        const data = JSON.stringify({ verificationCode, verifyType });
        await setData(account, data, expiresIn);
        return true;
    } catch (error) {
        throw error;
    }
}
// 生成6位数字验证码
function generateRandomCode() {
    const randomCode = Math.floor(100000 + Math.random() * 900000).toString();
    return randomCode;
}
```

3. 发送短信验证码
>注意这里有个坑，阿里云的短信验证码接口无论发送成功与否都会返回状态码200，所以不能通过返回的HTTP状态码来判断是否发送成功

```javascript
async function sendSms(phoneNumbers, Code) {
    try {
        const templateParam = JSON.stringify({
            code: Code
        });
        await Client.main({
            phoneNumbers: phoneNumbers,
            templateParam: templateParam
        });
        return { code: 0, msg: "短信验证码发送成功，请注意查收" };
    } catch (error) {
        return { code: 1, msg: "短信验证码发送失败" + error.message };
    }
}
```

4. 发送邮件验证码

```javascript
let mailOptions = {
        from: `"黄梅戏文化平台" <${process.env.EMAIL_USER}>`,
        to: email,
        subject: subject,
        text: message
    };
    return new Promise((resolve, reject) => {
        transporter.sendMail(mailOptions, (error, info) => {
            if (error) {
                reject({ code: 1, msg: "邮箱验证码发送失败" });
            } else {
                resolve({ code: 0, msg: "邮箱验证码发送成功，请注意查收" });
            }
        });
    });
}
```
效果![](https://r2.lm379.cn/2024/12/0abdf24d77f2112823008f7f8d7499b6.png)
5. GPT总结视频，完整项目见[lm379/video_summarize: video_summarize（视频总结：上传视频通过prompt对视频内容经行总结输出）](https://github.com/lm379/video_summarize)     
下面是针对本项目的扩展代码，使用Python的Flask框架开启HTTP接口，让Node.js可以直接调用该接口，而不需要引入Pyrunner之类的库

```python
from flask import Flask, request, jsonify
import json
import argparse
import whisper
import os
from utils import OpenaiMgr

app = Flask(__name__)

# Whisper模型
WHISPER_MODEL = "large-v2.pt"
PROJECT_PATH = os.path.dirname(os.path.dirname(__file__))
UP_PROJECT_PATH = os.path.join(PROJECT_PATH, "video_summarize")
MODE_PATH = os.path.join(UP_PROJECT_PATH, "model")
print(MODE_PATH+"/" + WHISPER_MODEL)
whisper_model =  whisper.load_model(MODE_PATH+ "/" + WHISPER_MODEL)  

def whisper_to_str(content):
    return whisper_model.transcribe(content)  
    
@app.route('/summarize', methods=['POST'])
def summarize_video():
    data = request.json
    video_path = data.get('video_path')
    prompt = data.get('prompt')
    print("1.提取视频文案")
    video_content = whisper_to_str(video_path)
    question = prompt + video_content["text"]
    print("2.GPT总结文案")
    s1, content = OpenaiMgr.call(question)
    print(content)
    return jsonify({"code":0, "content": content})  

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=7888)
```