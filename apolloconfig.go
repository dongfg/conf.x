package main

import (
	"context"
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type XN struct {
	Name       string   `yaml:"name"`
	LocalPath  string   `yaml:"localPath"`
	PostUpdate []string `yaml:"postUpdate"`
}

type X struct {
	AppID      string `yaml:"appID"`
	Cluster    string `yaml:"cluster"`
	Host       string `yaml:"host"`
	Secret     string `yaml:"secret"`
	Namespaces []XN   `yaml:"namespaces"`
}

type ConfigChangeListener struct {
	namespaceMap map[string]XN
}

var PropSuffix = ".properties"

func (c *ConfigChangeListener) OnNewestChange(event *storage.FullChangeEvent) {
	n, ok := c.namespaceMap[event.Namespace]
	newest := event.Changes
	if !ok {
		n, ok = c.namespaceMap[event.Namespace+PropSuffix]
		if !ok {
			log.Printf("namespace %s not found", event.Namespace)
			return
		}
	}
	var raw = ""
	if !strings.HasSuffix(n.Name, PropSuffix) {
		raw = newest["content"].(string)
	} else {
		for k, v := range newest {
			raw += fmt.Sprintf("%s=%v\n", k, v)
		}
	}
	if sync(n, raw) {
		postUpdate(n.Name, n.PostUpdate)
	}
}

func (c *ConfigChangeListener) OnChange(*storage.ChangeEvent) {
}

func watch(x X) {
	var namespaceName = make([]string, 0)
	var namespaceMap = make(map[string]XN)
	for _, n := range x.Namespaces {
		namespaceName = append(namespaceName, n.Name)
		namespaceMap[n.Name] = n
	}
	c := &config.AppConfig{
		AppID:          x.AppID,
		Cluster:        x.Cluster,
		IP:             x.Host,
		NamespaceName:  strings.Join(namespaceName, ","),
		IsBackupConfig: false,
		Secret:         x.Secret,
	}
	client, _ := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	client.AddChangeListener(&ConfigChangeListener{
		namespaceMap,
	})
	log.Printf("[Init] 获取配置 %s 成功\n", c.NamespaceName)
}

// sync config to local, return false if success
func sync(n XN, raw string) bool {
	compared := compare(n, raw)
	if compared == 0 { // same
		return false
	}
	if compared == -1 {
		log.Printf("[Init] %s --> %s\n", n.Name, n.LocalPath)
	} else {
		log.Printf("[Update] %s --> %s\n", n.Name, n.LocalPath)
	}
	_ = os.WriteFile(n.LocalPath, []byte(raw), 0644)
	return true
}

// compare between config and local, -1: local not exist, 0: same, 1: not same
func compare(n XN, raw string) int {
	fb, err := os.ReadFile(n.LocalPath)
	if err != nil {
		return -1
	}
	if strings.HasSuffix(n.Name, PropSuffix) {
		oldProps, err := parseProperties(fb)
		if err != nil {
			return 1
		}
		newProps, err := parseProperties([]byte(raw))
		if err != nil {
			return 1
		}
		if compareProperties(oldProps, newProps) {
			return 0
		}
	} else if string(fb) == raw {
		return 0
	}
	return 1
}

func postUpdate(namespace string, args []string) {
	if len(args) == 0 {
		return
	}
	if _, err := exec.LookPath(args[0]); err != nil {
		log.Println(fmt.Errorf("[PostUpdate] %s@error: 命令不存在: %s", namespace, args[0]))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // 设置超时
	defer cancel()
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(fmt.Errorf("[PostUpdate] %s@error: %v, %s", namespace, err, string(output)))
		return
	}
	log.Println(fmt.Errorf("[PostUpdate] %s@%s", namespace, args))
}
