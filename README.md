# Git Webhook 服务

一个轻量级的 Git Webhook 服务器，用于自动化部署和构建。支持 GitLab Webhook，可以通过配置文件动态管理多个项目。完全支持 Windows 和 Linux 平台。

## 特性

- 支持多个项目同时监控
- 配置文件热重载（无需重启服务）
- 支持分支过滤
- 安全验证（通过 Webhook Secret）
- 自定义命令执行
- 完整的跨平台支持（Windows/Linux）
  - Windows 下使用 cmd /C 执行命令
  - Linux 下使用 sh -c 执行命令
  - 智能路径处理

## 安装

```bash
# 克隆仓库
git clone [your-repository-url]

# 进入项目目录
cd git-webhook

# 安装依赖
go mod download

# 编译
## Windows
go build -o git-webhook.exe

## Linux/Mac
GOOS=linux go build -o git-webhook
```

## 配置

创建 `config.yaml` 文件：

```yaml
server:
  port: 8080
  secret: "your-gitlab-webhook-secret"

projects:
  # Windows 配置示例
  - name: "windows-project"
    path: "D:/projects/myapp"  # Windows 路径
    branch: "main"
    secret: "project-secret-1"
    commands:
      - "git pull origin main"
      - "npm install"
      - "npm run build"

  # Linux 配置示例
  - name: "linux-project"
    path: "/var/www/myapp"    # Linux 路径
    branch: "main"
    secret: "project-secret-2"
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
  - `path`: 项目本地路径（支持 Windows 和 Linux 格式）
  - `branch`: 触发构建的分支
  - `secret`: 项目特定的 Webhook 密钥
  - `commands`: 要执行的命令列表

## 使用方法

1. 启动服务：
```bash
# Windows
./git-webhook.exe

# Linux/Mac
./git-webhook
```

2. 在 GitLab 项目设置中配置 Webhook：
   - URL: `http://your-server:8080/webhook`
   - Secret Token: 设置为项目配置中的 secret
   - 触发事件: 选择 Push events

3. 测试 Webhook：
   - 在 GitLab 的 Webhook 设置页面点击"Test"按钮
   - 或进行一次代码提交

## 跨平台注意事项

1. Windows 环境：
   - 路径可以使用正斜杠（/）或双反斜杠（\\）
   - 支持 Windows 风格的命令（如 dir, copy 等）
   - 确保 Git 和其他命令行工具在系统环境变量中

2. Linux 环境：
   - 使用标准的 Unix 路径格式
   - 确保脚本有执行权限
   - 建议使用绝对路径

3. 通用建议：
   - 配置文件中的路径会自动被规范化
   - 命令输出会实时显示在控制台
   - 错误信息包含完整的命令输出

## 安全建议

1. 文件权限：
   - Windows: 确保运行用户有足够的文件访问权限
   - Linux: 建议使用适当的文件权限（如 755）

2. 网络安全：
   - 建议使用 HTTPS 进行 Webhook 通信
   - 为每个项目设置唯一的 Secret
   - 定期更新 Secret

3. 配置文件：
   - 支持热重载，修改后立即生效
   - 确保配置文件格式正确
   - 建议限制配置文件的访问权限

## 日志

服务会在控制台输出详细的运行日志，包括：
- 服务启动信息
- Webhook 请求处理状态
- 命令执行结果和输出
- 配置文件重载状态
- 错误信息和调试信息

## 许可证

本项目采用 MIT 许可证。查看 [LICENSE](LICENSE) 文件了解更多详情。

MIT 许可证是一个宽松的软件许可证，这意味着：

- ✔️ 可以商业使用
- ✔️ 可以修改源代码
- ✔️ 可以私有使用
- ✔️ 可以分发
- ✔️ 可以再授权

唯一的要求是在您的项目中包含原始许可证和版权声明。 