# 黄梅戏数字化传播平台

基于 Go + Gin 框架开发的黄梅戏数字化传播平台，支持用户注册、登录、作品上传、点赞、收藏、评论等功能。

## 技术栈

- **后端框架**: Go + Gin
- **前端框架**: Vue 3 + TypeScript + Vite
- **数据库**: PostgreSQL + Redis
- **对象存储**: AWS S3 或 S3 协议兼容对象存储
- **认证**: Refresh Token + Access Token
- **容器化**: Docker 多阶段构建

## 快速开始

### 方式一：Docker 部署（推荐）

#### 1. 构建 Docker 镜像

```bash
docker build -t hmx:latest .
```

#### 2. 运行容器

```bash
docker run -d \
  -p 15379:15379 \
  -e DATABASE_URL="your_database_url" \
  -e REDIS_URL="your_redis_url" \
  -e AWS_ACCESS_KEY="your_aws_key" \
  -e AWS_SECRET_KEY="your_aws_secret" \
  -e AWS_REGION="your_region" \
  -e S3_BUCKET="your_bucket" \
  -e JWT_SECRET="your_jwt_secret" \
  -e SMTP_HOST="smtp.example.com" \
  -e SMTP_PORT="587" \
  -e SMTP_USER="your_email" \
  -e SMTP_PASSWORD="your_password" \
  --name hmx-app \
  hmx:latest
```

#### 3. 查看日志

```bash
docker logs -f hmx-app
```

### 方式二：本地开发

#### 1. 安装前端依赖并构建

```bash
cd frontend
npm install
npm run build
cd ..
```

#### 2. 安装 Go 依赖

```bash
go mod tidy
```

#### 3. 配置环境变量

创建配置文件或设置环境变量，包括：
- 数据库连接信息
- Redis 连接信息
- AWS S3 配置
- SMTP 邮件配置
- JWT 密钥

#### 4. 构建项目

```bash
go build -o hmx ./cmd/server/main.go
```

#### 5. 运行项目

```bash
./hmx
```

或直接运行：

```bash
go run cmd/server/main.go
```

---

## 项目结构

```
hmx/
├── cmd/server/          # 服务入口
│   ├── main.go         # 主程序
│   └── static/         # 前端构建产物（由 frontend/dist 生成）
├── config/             # 配置管理
├── database/           # 数据库初始化
├── frontend/           # Vue 3 前端项目
│   ├── src/           # 源代码
│   ├── package.json   # 前端依赖
│   └── vite.config.ts # Vite 配置
├── internal/           # 内部包
│   ├── api/v1/        # API 路由和处理器
│   ├── models/        # 数据模型
│   ├── repository/    # 数据访问层
│   └── services/      # 业务逻辑层
├── pkg/               # 公共工具包
├── Dockerfile         # Docker 构建文件
├── .dockerignore      # Docker 忽略文件
├── go.mod             # Go 依赖管理
└── README.md          # 项目文档
```

---

### 认证 (Authentication)

#### 发送验证码
发送邮箱验证码用于注册。

**接口**: `POST /auth/send-code`

**请求体**:
```json
{
  "email": "user@example.com"
}
```

**成功响应** (200):
```json
{
  "message": "Verification code sent (check logs)."
}
```

---

#### 用户注册
使用邮箱验证码注册新用户。

**接口**: `POST /auth/register`

**请求体**:
```json
{
  "username": "张三",
  "phone": "13800138000",
  "email": "user@example.com",
  "password": "password123",
  "code": "123456"
}
```

**成功响应** (201):
```json
{
  "message": "Registration successful"
}
```

**错误响应**:
- `400`: 验证码错误或已过期
- `409`: 邮箱或手机号已被注册

---

#### 用户登录
使用邮箱和密码登录，返回 Access Token 及 Refresh Token。

**接口**: `POST /auth/login`

**请求体**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**成功响应** (200):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**错误响应**:
- `401`: 邮箱或密码错误

#### 重置密码
使用邮箱和验证码重置密码。

**接口**: `POST /auth/forget`

**请求体**:
```json
{
    "email": "email@example.com",
    "code": "123456",
    "new_password": "new_password"
}
```

**成功响应** (200):
```json
{
	"code": 200,
	"msg": "success",
	"data": {
		"message": "Password reset successful"
	}
}
```

---

#### 刷新 Token
使用 Refresh Token 获取新的 Access Token。

**接口**: `POST /auth/refresh`

**请求体**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**成功响应** (200):
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**错误响应**:
- `401`: Refresh Token 无效或已过期

---

#### 登出
使当前的 Refresh Token 失效。

**接口**: `POST /auth/logout`

**请求头**: 需要 Access Token

**请求体**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**成功响应** (200):
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "message": "Logged out successfully"
  }
}
```


---

### 用户 (Users)

以下接口需要在请求头中携带 Access Token:
```
Authorization: Bearer <access token>
```

#### 获取当前用户信息
获取当前登录用户的详细信息。

**接口**: `GET /users/me`

**成功响应** (200):
```json
{
  "user_id": 1,
  "username": "张三",
  "phone": "13800138000",
  "email": "user@example.com",
  "sex": "Other",
  "icon": "https://s3.amazonaws.com/bucket/avatars/user1.jpg",
  "role": "User",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

---

#### 获取用户点赞的作品列表
获取当前用户点赞的所有作品。

**接口**: `GET /users/me/likes`

**查询参数**:
- `page` (int, 可选): 页码，默认 1
- `page_size` (int, 可选): 每页数量，默认 10

**成功响应** (200):
```json
{
  "list": [
    {
      "opera_id": 1,
      "opera_title": "天仙配",
      "release_date": "2025-01-01",
      "duration": "02:30:00",
      "music_path": "https://s3.amazonaws.com/bucket/music/opera1.mp3",
      "video_path": "https://s3.amazonaws.com/bucket/videos/opera1.mp4",
      "description": "经典黄梅戏作品",
      "avatar": "https://s3.amazonaws.com/bucket/avatars/opera1.jpg",
      "ai_summary": "这是一部...",
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "total": 25,
    "page": 1,
    "page_size": 10
  }
}
```

---

#### 获取用户收藏的作品列表
获取当前用户收藏的所有作品。

**接口**: `GET /users/me/favorites`

**查询参数**:
- `page` (int, 可选): 页码，默认 1
- `page_size` (int, 可选): 每页数量，默认 10

**响应格式**: 同"获取用户点赞的作品列表"

---

#### 获取用户观看历史
获取当前用户的观看历史记录。

**接口**: `GET /users/me/history`

**请求头**: 需要 Access Token

**查询参数**:
- `page` (int, 可选): 页码，默认 1
- `page_size` (int, 可选): 每页数量，默认 10

**成功响应** (200):
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "list": [
      {
        "opera_id": 1,
        "opera_title": "天仙配",
        "avatar": "https://...",
        "last_watched_at": "2025-01-01T12:00:00Z",
        "watch_duration": 1800
      }
    ],
    "pagination": {
      "total": 50,
      "page": 1,
      "page_size": 10
    }
  }
}
```

---

#### 更新用户信息
更新当前用户的个人资料。

**接口**: `PUT /users/me`

**请求头**: 需要 Access Token

**请求体**:
```json
{
  "username": "新用户名",
  "phone": "13900139000",
  "sex": "Male",
  "icon": "avatars/uuid-avatar.jpg"
}
```

**字段说明**:
- 所有字段均为可选
- `sex` 可选值: `Male`, `Female`, `Other`
- `icon` 为头像的 S3 object_key

**成功响应** (200):
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "message": "Profile updated successfully"
  }
}
```

---

#### 修改密码
修改当前用户的密码。

**接口**: `POST /users/me/password`

**请求头**: 需要 Access Token

**请求体**:
```json
{
  "old_password": "old_password123",
  "new_password": "new_password456"
}
```

**成功响应** (200):
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "message": "Password updated successfully"
  }
}
```

**错误响应**:
- `401`: 旧密码错误

---

#### 发送验证码到当前用户
向当前登录用户的邮箱发送验证码（用于敏感操作验证）。

**接口**: `POST /users/me/email-code`

**请求头**: 需要 Access Token

**成功响应** (200):
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "message": "Verification code sent"
  }
}
```

---

### 文件上传 (Uploads)

#### 请求预签名上传 URL
获取 S3 预签名上传 URL，用于客户端直接上传文件到 S3。

**接口**: `POST /uploads/presign`

**请求头**: 需要 Access Token

**请求体**:
```json
{
  "filename": "my_video.mp4",
  "content_type": "video/mp4",
  "upload_type": "videos"
}
```

**upload_type 可选值**:
- `videos`: 视频文件
- `avatars`: 头像/封面图片

**成功响应** (200):
```json
{
  "upload_url": "https://s3.amazonaws.com/bucket/videos/uuid-my_video.mp4?X-Amz-...",
  "object_key": "videos/uuid-my_video.mp4"
}
```

**使用方式**:
1. 客户端调用此接口获取 `upload_url` 和 `object_key`
2. 客户端使用 `upload_url` 直接 PUT 上传文件到 S3
3. 上传成功后，使用 `object_key` 创建作品

---

### 作品 (Operas)

#### 获取作品列表
获取所有已发布的作品列表，支持分页。

**接口**: `GET /operas/`

**查询参数**:
- `page` (int, 可选): 页码，默认 1
- `page_size` (int, 可选): 每页数量，默认 10

**成功响应** (200):
```json
{
  "list": [
    {
      "opera_id": 1,
      "opera_title": "天仙配",
      "release_date": "2025-01-01",
      "duration": "02:30:00",
      "music_path": "https://s3.amazonaws.com/bucket/music/opera1.mp3",
      "video_path": "https://s3.amazonaws.com/bucket/videos/opera1.mp4",
      "description": "经典黄梅戏作品",
      "avatar": "https://s3.amazonaws.com/bucket/avatars/opera1.jpg",
      "ai_summary": "这是一部...",
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

---

#### 获取作品详情
根据作品 ID 获取作品详细信息。

**接口**: `GET /operas/:id`

**路径参数**:
- `id` (int): 作品 ID

**成功响应** (200):
```json
{
  "opera_id": 1,
  "opera_title": "天仙配",
  "release_date": "2025-01-01",
  "duration": "02:30:00",
  "music_path": "https://s3.amazonaws.com/bucket/music/opera1.mp3",
  "video_path": "https://s3.amazonaws.com/bucket/videos/opera1.mp4",
  "description": "经典黄梅戏作品",
  "avatar": "https://s3.amazonaws.com/bucket/avatars/opera1.jpg",
  "ai_summary": "这是一部...",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

**错误响应**:
- `404`: 作品不存在

---

#### 记录观看历史
记录用户观看作品的历史（游客也可以记录，但不会关联到用户）。

**接口**: `POST /operas/:id/history`

**请求头**: 可选 Access Token（如果提供则关联用户）

**路径参数**:
- `id` (int): 作品 ID

**请求体**:
```json
{
  "watch_duration": 1800
}
```

**字段说明**:
- `watch_duration`: 观看时长（秒）

**成功响应** (200):
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "message": "History recorded"
  }
}
```

---

#### 创建作品
用户上传新作品（需要先通过预签名上传文件）。

**接口**: `POST /operas/`

**请求头**: 需要 Access Token

**请求体**:
```json
{
  "title": "女驸马",
  "description": "黄梅戏经典剧目",
  "video_path": "videos/uuid-video.mp4",
  "avatar_path": "avatars/uuid-cover.jpg",
  "artist_ids": [1, 2, 3]
}
```

**字段说明**:
- `video_path`: 必填，从预签名接口获取的 `object_key`
- `avatar_path`: 可选，作品封面图的 `object_key`
- `artist_ids`: 可选，关联的艺术家 ID 数组

**成功响应** (201):
```json
{
  "message": "Opera created successfully",
  "opera_id": 123
}
```

---

#### 点赞/取消点赞作品
对指定作品进行点赞或取消点赞。

**接口**: `POST /operas/:id/like`

**请求头**: 需要 Access Token

**路径参数**:
- `id` (int): 作品 ID

**成功响应** (200):
```json
{
  "message": "Liked" 
}
```
或
```json
{
  "message": "Unliked"
}
```

---

#### 收藏/取消收藏作品
对指定作品进行收藏或取消收藏。

**接口**: `POST /operas/:id/favorite`

**请求头**: 需要 Access Token

**路径参数**:
- `id` (int): 作品 ID

**成功响应** (200):
```json
{
  "message": "Favorited"
}
```
或
```json
{
  "message": "Unfavorited"
}
```

---

#### 分享作品
记录作品分享（游客也可以分享，但不会关联到用户）。

**接口**: `POST /operas/:id/share`

**请求头**: 可选 Access Token

**路径参数**:
- `id` (int): 作品 ID

**成功响应** (200 或 201):
```json
{
  "code": 201,
  "msg": "success",
  "data": {
    "message": "Share recorded",
    "new_share": true
  }
}
```

**字段说明**:
- `new_share`: `true` 表示首次分享（返回 201），`false` 表示重复分享（返回 200）

---

#### 创建评论
对指定作品发表评论。

**接口**: `POST /operas/:id/comments`

**请求头**: 需要 Access Token

**路径参数**:
- `id` (int): 作品 ID

**请求体**:
```json
{
  "content": "这部作品太精彩了！",
  "parent_id": 0
}
```

**字段说明**:
- `content`: 必填，评论内容
- `parent_id`: 可选，如果是回复某条评论，填写父评论 ID；顶级评论填 0 或不填

**成功响应** (201):
```json
{
  "message": "Comment created successfully",
  "comment_id": 456
}
```

---

#### 获取作品评论列表
获取指定作品的评论列表。

**接口**: `GET /operas/:id/comments`

**请求头**: 可选 Access Token（登录用户可以看到自己的点赞状态）

**路径参数**:
- `id` (int): 作品 ID

**查询参数**:
- `page` (int, 可选): 页码，默认 1
- `page_size` (int, 可选): 每页数量，默认 10

**成功响应** (200):
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "list": [
      {
        "comment_id": 1,
        "user_id": 10,
        "username": "张三",
        "user_icon": "https://...",
        "content": "这部作品太精彩了！",
        "parent_id": 0,
        "created_at": "2025-01-01T10:00:00Z",
        "replies": []
      }
    ],
    "pagination": {
      "total": 100,
      "page": 1,
      "page_size": 10
    }
  }
}
```

---

### 艺术家 (Artists)

#### 获取艺术家列表
获取所有艺术家的列表，支持分页。

**接口**: `GET /artists/`

**查询参数**:
- `page` (int, 可选): 页码，默认 1
- `page_size` (int, 可选): 每页数量，默认 10

**成功响应** (200):
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "list": [
      {
        "artist_id": 1,
        "name": "韩再芬",
        "avatar": "https://s3.amazonaws.com/bucket/artists/artist1.jpg",
        "bio": "著名黄梅戏表演艺术家",
        "birth_date": "1968-03-20",
        "created_at": "2025-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "total": 50,
      "page": 1,
      "page_size": 10
    }
  }
}
```

---

#### 获取艺术家详情
根据艺术家 ID 获取详细信息及其作品列表。

**接口**: `GET /artists/:id`

**路径参数**:
- `id` (int): 艺术家 ID

**成功响应** (200):
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "artist": {
      "artist_id": 1,
      "name": "韩再芬",
      "avatar": "https://...",
      "bio": "著名黄梅戏表演艺术家",
      "birth_date": "1968-03-20",
      "created_at": "2025-01-01T00:00:00Z"
    },
    "operas": [
      {
        "opera_id": 1,
        "opera_title": "女驸马",
        "avatar": "https://...",
        "description": "经典剧目"
      }
    ]
  }
}
```

**错误响应**:
- `404`: 艺术家不存在

---

### 评论 (Comments)

#### 删除评论
删除指定的评论（仅评论作者或管理员可删除）。

**接口**: `DELETE /comments/:id`

**请求头**: 需要 Access Token

**路径参数**:
- `id` (int): 评论 ID

**成功响应** (200):
```json
{
  "message": "Comment deleted successfully"
}
```

**错误响应**:
- `403`: 无权删除此评论
- `404`: 评论不存在

---

### 管理员 (Admin)

以下接口需要管理员权限（`role = "Administrator"`）。

#### 删除作品
管理员删除指定作品。

**接口**: `DELETE /admin/operas/:id`

**请求头**: 需要 Access Token (管理员)

**路径参数**:
- `id` (int): 作品 ID

**成功响应** (200):
```json
{
  "message": "Opera deleted successfully"
}
```

**错误响应**:
- `403`: 无管理员权限
- `404`: 作品不存在

---

## 状态码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未授权（未登录或 Token 无效） |
| 403 | 禁止访问（权限不足） |
| 404 | 资源不存在 |
| 409 | 资源冲突（如邮箱已注册） |
| 500 | 服务器内部错误 |

---

## 认证说明

### Access Token 使用

1. 登录成功后，服务器返回 Access Token
2. 后续需要认证的请求，在请求头中携带：
   ```
   Authorization: Bearer <your_access_token>
   ```
3. Token 包含用户 ID、用户名、角色等信息
4. 中间件会自动验证 Token 并将用户信息注入到上下文中

### 权限等级

- **游客**: 可以浏览作品列表和详情
- **普通用户** (`User`): 可以上传作品、点赞、收藏、评论
- **管理员** (`Administrator`): 拥有所有权限，可以删除任何作品和评论

---

## 文件上传流程

1. **客户端请求预签名 URL**:
   ```
   POST /api/v1/uploads/presign
   {
     "filename": "video.mp4",
     "content_type": "video/mp4",
     "upload_type": "videos"
   }
   ```

2. **服务器返回上传信息**:
   ```json
   {
     "upload_url": "https://s3.amazonaws.com/...",
     "object_key": "videos/uuid-video.mp4"
   }
   ```

3. **客户端直接上传到 S3**:
   ```bash
   curl -X PUT "<upload_url>" \
     -H "Content-Type: video/mp4" \
     --data-binary @video.mp4
   ```

4. **上传成功后创建作品**:
   ```
   POST /api/v1/operas/
   {
     "title": "作品标题",
     "video_path": "videos/uuid-video.mp4",
     "description": "..."
   }
   ```

