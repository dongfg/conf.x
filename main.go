package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Config X",
		Usage: "watch config from config center\nhttps://github.com/dongfg/conf.x",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Usage:    "Load configuration from 'FILE'",
				Aliases:  []string{"c"},
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			fbs, err := os.ReadFile(c.String("config"))
			if err != nil {
				return fmt.Errorf("failed to read config file: %v", err)
			}
			var x X
			if err := yaml.Unmarshal(fbs, &x); err != nil {
				return fmt.Errorf("failed to parse YAML: %v", err)
			}
			watch(x)
			waitForExitSignal()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// 阻塞直到收到退出信号（如 Ctrl+C）
func waitForExitSignal() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
