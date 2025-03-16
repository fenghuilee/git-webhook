# Git Webhook 服务

一个轻量级的 Git Webhook 服务器，用于自动化部署和构建。支持 GitLab Webhook，可以通过配置文件动态管理多个项目。

## 特性

- 支持多个项目同时监控
- 配置文件热重载（无需重启服务）
- 支持分支过滤
- 安全验证（通过 Webhook Secret）
- 自定义命令执行
- 跨平台支持（Windows/Linux）

## 安装

```bash
# 克隆仓库
git clone [your-repository-url]

# 进入项目目录
cd git-webhook

# 安装依赖
go mod download

# 编译
go build -o git-webhook.exe  # Windows
go build -o git-webhook      # Linux/Mac
```

## 配置

创建 `config.yaml` 文件：

```yaml
server:
  port: 8080
  secret: "your-gitlab-webhook-secret"

projects:
  - name: "example-project"
    path: "D:/path/to/your/project"  # Windows 路径示例
    branch: "main"
    secret: "project-specific-secret"
    commands:
      - "git pull origin main"
      - "npm install"
      - "npm run build"
```

### 配置说明

- `server.port`: 服务器监听端口
- `server.secret`: 全局 Webhook 密钥（可选）
- `projects`: 项目列表
  - `name`: 项目名称
  - `path`: 项目本地路径
  - `branch`: 触发构建的分支
  - `secret`: 项目特定的 Webhook 密钥
  - `commands`: 要执行的命令列表

## 使用方法

1. 启动服务：
```bash
./git-webhook.exe  # Windows
./git-webhook      # Linux/Mac
```

2. 在 GitLab 项目设置中配置 Webhook：
   - URL: `http://your-server:8080/webhook`
   - Secret Token: 设置为项目配置中的 secret
   - 触发事件: 选择 Push events

3. 测试 Webhook：
   - 在 GitLab 的 Webhook 设置页面点击"Test"按钮
   - 或进行一次代码提交

## 注意事项

1. Windows 环境注意事项：
   - 路径使用正斜杠（/）或双反斜杠（\\）
   - 确保执行命令的用户有足够权限
   - Git 命令需要配置在系统环境变量中

2. 安全建议：
   - 建议使用 HTTPS 进行 Webhook 通信
   - 为每个项目设置唯一的 Secret
   - 定期更新 Secret

3. 配置文件更新：
   - 配置文件支持热重载，修改后立即生效
   - 确保配置文件格式正确，错误的格式可能导致加载失败

## 日志

服务会在控制台输出基本的运行日志，包括：
- 服务启动信息
- Webhook 请求处理状态
- 命令执行结果
- 配置文件重载状态

## 许可证

[您的许可证类型] 