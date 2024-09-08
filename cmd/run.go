package main

import (
	"log"
	"os"
	container "simple-docker/common/contianer"
	"strings"

	"github.com/urfave/cli/v2"
)

var RunCommand = &cli.Command{
	Name:  "run",
	Usage: "run a container",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "it",
			Usage: "interactive and pseudo tty", // 额外开启的交互式终端
		},
		&cli.StringFlag{
			Name: "mem",
			Usage: "-mem 100m", // 限制内存使用
		},
		&cli.StringFlag{
			Name: "cpu",
			Usage: "-cpu 100", // 限制cpu使用
		},
		&cli.StringFlag{
			Name: "cpuset",
			Usage: "-cpuset 2,4",
		},
	},
	Action: func(ctx *cli.Context) error {
		var (
			cmd = ctx.Args().Get(0)
			tty = ctx.Bool("it")
		)
		Run(tty, []string{cmd})
		return nil
	},
}

func Run(tty bool, cmds []string) {
	parent, writePipe, err := container.NewParentProcess(tty)
	if err != nil {
		log.Println(err)
	}
	if err := parent.Start(); err != nil {
		log.Println(err)
	}
	// 将参数发送通过管道发送给子进程
	SendInitCommand(cmds, writePipe)
	_ = parent.Wait()
	os.Exit(-1)
}

func SendInitCommand(cmds []string, writePipe *os.File) {
	command := strings.Join(cmds, " ")
	_, _ = writePipe.WriteString(command)
	_ = writePipe.Close()
}

var InitCommand = &cli.Command{
	Name:  "init",
	Usage: "init container process",
	Action: func(ctx *cli.Context) error {
		cmd := ctx.Args().Get(0)
		return container.RunContainerInitProcess(cmd, nil)
	},
}
