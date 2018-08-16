// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/open-falcon/falcon-plus/g"
	"github.com/spf13/cobra"
)

// open-falcon start 子命令定义
var Start = &cobra.Command{
	Use:   "start [Module ...]",
	Short: "Start Open-Falcon modules",
	Long: `
Start the specified Open-Falcon modules and run until a stop command is received.
A module represents a single node in a cluster.
Modules:
	` + "all " + strings.Join(g.AllModulesInOrder, " "),
	RunE:          start,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// 对应 --preq-order 参数
var PreqOrderFlag bool
// 对应 --console-output 参数
var ConsoleOutputFlag bool

func cmdArgs(name string) []string {
	return []string{"-c", g.Cfg(name)}
}

func openLogFile(name string) (*os.File, error) {
	logDir := g.LogDir(name)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logPath := g.LogPath(name)
	logOutput, err := os.OpenFile(logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return logOutput, nil
}

func execModule(co bool, name string) error {
	// 下面是agent模块的启动示例:
	// 		/<prefix>/agent/bin/falcon-agent -c /<prefix>/agent/config/cfg.json
	cmd := exec.Command(g.Bin(name), cmdArgs(name)...)

	// 对应 --console-output 参数
	if co {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// 重定向 stdout & stderr 到日志文件
	// PS: 日志设计成接口, 支持多种backend比较好, 例如: rsyslog
	logOutput, err := openLogFile(name)
	if err != nil {
		return err
	}
	defer logOutput.Close()
	cmd.Stdout = logOutput
	cmd.Stderr = logOutput
	return cmd.Start()
}

func checkStartReq(name string) error {
	if !g.HasModule(name) {
		return fmt.Errorf("%s doesn't exist", name)
	}

	if !g.HasCfg(name) {
		r := g.Rel(g.Cfg(name))
		return fmt.Errorf("expect config file: %s", r)
	}

	return nil
}

func isStarted(name string) bool {
	// 1s内, 每100ms检测模块是否成功启动
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if g.IsRunning(name) {
				return true
			}
		case <-time.After(time.Second):
			return false
		}
	}
}

func start(c *cobra.Command, args []string) error {
	// 参数去重, 防止重复启动模块
	args = g.RmDup(args)

	// 是否需要按特定顺序启动模块, 对应 --preq-order 参数
	if PreqOrderFlag {
		args = g.PreqOrder(args)
	}

	// 默认开启全部模块
	if len(args) == 0 {
		args = g.AllModulesInOrder
	}

	for _, moduleName := range args {
		// 判断模块和配置文件是否存在
		if err := checkStartReq(moduleName); err != nil {
			return err
		}

		// 跳过已运行的模块
		if g.IsRunning(moduleName) {
			fmt.Print("[", g.ModuleApps[moduleName], "] ", g.Pid(moduleName), "\n")
			continue
		}

		// 启动模块
		if err := execModule(ConsoleOutputFlag, moduleName); err != nil {
			return err
		}

		// 打印启动状态

		if isStarted(moduleName) {
			fmt.Print("[", g.ModuleApps[moduleName], "] ", g.Pid(moduleName), "\n")
			continue
		}

		return fmt.Errorf("[%s] failed to start", g.ModuleApps[moduleName])
	}
	return nil
}
