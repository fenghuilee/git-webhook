package config

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port   int    `yaml:"port"`
		Secret string `yaml:"secret"`
	} `yaml:"server"`
	Projects []Project `yaml:"projects"`
	mu       sync.RWMutex
}

type Project struct {
	Name     string   `yaml:"name"`
	Path     string   `yaml:"path"`
	Branch   string   `yaml:"branch"`
	Secret   string   `yaml:"secret"`
	Commands []string `yaml:"commands"`
}

// LoadConfig 加载配置文件并返回配置对象
func LoadConfig(filename string) (*Config, error) {
	config := &Config{}
	if err := config.load(filename); err != nil {
		return nil, err
	}

	// 启动配置文件监控
	go config.watch(filename)

	return config, nil
}

// load 从文件加载配置
func (c *Config) load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}

	// 规范化所有项目路径
	for i := range c.Projects {
		c.Projects[i].Path = filepath.Clean(c.Projects[i].Path)
	}

	return nil
}

// watch 监控配置文件变化
func (c *Config) watch(filename string) {
	// 规范化配置文件路径
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()

	// 添加文件到监控列表
	if err := watcher.Add(filepath.Dir(absPath)); err != nil {
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// 文件被修改时重新加载配置
			if filepath.Clean(event.Name) == absPath && event.Has(fsnotify.Write) {
				if err := c.load(filename); err != nil {
					// 这里可以添加错误日志记录
					continue
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			// 这里可以添加错误日志记录
			_ = err
		}
	}
}

// GetProjects 线程安全地获取项目列表
func (c *Config) GetProjects() []Project {
	c.mu.RLock()
	defer c.mu.RUnlock()
	projects := make([]Project, len(c.Projects))
	copy(projects, c.Projects)
	return projects
}

// GetServerConfig 线程安全地获取服务器配置
func (c *Config) GetServerConfig() struct {
	Port   int
	Secret string
} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return struct {
		Port   int
		Secret string
	}{
		Port:   c.Server.Port,
		Secret: c.Server.Secret,
	}
}
