package main

import (
	"fmt"
	"git-webhook/config"
	"git-webhook/handler"
	"log"
	"net/http"
)

func main() {
	// 加载配置文件
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建 webhook 处理器
	webhookHandler := handler.NewWebhookHandler(cfg)

	// 设置路由
	http.HandleFunc("/webhook", webhookHandler.Handle)

	// 启动服务器
	serverConfig := cfg.GetServerConfig()
	addr := fmt.Sprintf(":%d", serverConfig.Port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
