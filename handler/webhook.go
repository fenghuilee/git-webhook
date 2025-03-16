package handler

import (
	"encoding/json"
	"fmt"
	"git-webhook/config"
	"git-webhook/models"
	"io/ioutil"
	"net/http"
	"os/exec"
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

	body, err := ioutil.ReadAll(r.Body)
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
	for _, cmd := range project.Commands {
		command := exec.Command("sh", "-c", cmd)
		command.Dir = project.Path

		if err := command.Run(); err != nil {
			return fmt.Errorf("failed to execute command '%s': %v", cmd, err)
		}
	}
	return nil
}
