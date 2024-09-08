package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	// set log.txt as input
	// file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()
	// log.SetOutput(file)
	log.SetOutput(os.Stdout)

	app := cli.NewApp()
	app.Usage = "used for a simple container"
	app.Name = "simple-docker"
	app.Commands = []*cli.Command{
		RunCommand,
		InitCommand,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
