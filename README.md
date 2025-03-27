## Config X
> 自动从 Apollo 配置中心拉取namespace配置到本地并保持同步

[![lint](https://github.com/dongfg/conf.x/actions/workflows/lint.yaml/badge.svg)](https://github.com/dongfg/conf.x/actions/workflows/lint.yaml)
![GitHub Release](https://img.shields.io/github/v/release/dongfg/conf.x)


## 安装

go install

```shell
go install github.com/dongfg/conf.x@latest
```

or [download binary](https://github.com/dongfg/conf.x/releases)

## 使用

```shell
conf.x -c .\config.yaml 
```

## 配置文件说明
```yaml
appID: test # Apollo appID
cluster: default # Apollo cluster
host: # Apollo IP
secret: # 可选, Apollo secret
namespaces:
  - name: application.properties # namespace
    localPath: /tmp/application.properties # 同步到本地的文件路径
    postUpdate: # 可选, 配置更新以后执行 
      - cat
      - /tmp/application.properties
  - name: daemon.json
    localPath: /tmp/daemon.json
  - name: nginx.txt
    localPath: /tmp/nginx.conf
```

## Systemd Service 示例
```text
[Unit]
Description=Config X Service
After=network.target

[Service]
User=syncuser
Group=syncuser
WorkingDirectory=/opt/auto-sync
ExecStart=/usr/local/bin/conf.x -c config.local.yaml
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
```