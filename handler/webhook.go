package handler

import (
	"encoding/json"
	"fmt"
	"git-webhook/config"
	"git-webhook/models"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type WebhookHandler struct {
	config *config.Config
}

func NewWebhookHandler(cfg *config.Config) *WebhookHandler {
	return &WebhookHandler{
		config: cfg,
	}
}

func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	signature := r.Header.Get("X-Gitlab-Token")
	if signature == "" {
		http.Error(w, "Missing GitLab signature", http.StatusBadRequest)
		return
	}

	var payload models.GitLabWebhook
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Failed to parse webhook payload", http.StatusBadRequest)
		return
	}

	// 获取分支名称
	branch := strings.TrimPrefix(payload.Ref, "refs/heads/")

	// 查找匹配的项目配置
	for _, project := range h.config.GetProjects() {
		if project.Branch == branch {
			// 验证项目密钥
			if !verifySignature(project.Secret, signature) {
				http.Error(w, "Invalid signature", http.StatusUnauthorized)
				return
			}

			// 执行配置的命令
			if err := executeCommands(project); err != nil {
				http.Error(w, fmt.Sprintf("Failed to execute commands: %v", err), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "No matching project configuration found", http.StatusNotFound)
}

func verifySignature(secret, signature string) bool {
	return secret == signature
}

func executeCommands(project config.Project) error {
	// 规范化路径
	project.Path = filepath.Clean(project.Path)

	// 检查项目目录是否存在
	if _, err := os.Stat(project.Path); os.IsNotExist(err) {
		// 创建项目目录
		if err := os.MkdirAll(project.Path, 0755); err != nil {
			return fmt.Errorf("failed to create project directory: %v", err)
		}

		// 初始化 git 仓库并拉取指定分支
		commands := []string{
			fmt.Sprintf("git clone -b %s %s .", project.Branch, project.Repository),
		}

		for _, cmd := range commands {
			var command *exec.Cmd
			if runtime.GOOS == "windows" {
				command = exec.Command("cmd", "/C", cmd)
			} else {
				command = exec.Command("sh", "-c", cmd)
			}

			command.Dir = project.Path
			output, err := command.CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to execute init command '%s': %v\nOutput: %s", cmd, err, string(output))
			}
			fmt.Printf("Init command '%s' output:\n%s\n", cmd, string(output))
		}
	}

	// 执行配置的命令
	for _, cmd := range project.Commands {
		var command *exec.Cmd
		if runtime.GOOS == "windows" {
			command = exec.Command("cmd", "/C", cmd)
		} else {
			command = exec.Command("sh", "-c", cmd)
		}

		// 设置工作目录
		command.Dir = project.Path

		// 获取命令的输出
		output, err := command.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to execute command '%s': %v\nOutput: %s", cmd, err, string(output))
		}

		// 打印命令输出
		fmt.Printf("Command '%s' output:\n%s\n", cmd, string(output))
	}
	return nil
}
