package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/mrhinton101/fluyt/cmd/fluyt/cli"
	"github.com/mrhinton101/fluyt/cmd/fluyt/web"
	"github.com/mrhinton101/fluyt/internal/app/core/logger"
)

func main() {
	logger.ProgramLevel.Set(slog.LevelError)
	logfile := logger.InitLogger("fluytLogs.json")
	defer logfile.Close()

	if len(os.Args) != 2 {
		fmt.Println("Usage: fluyt <mode>")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	// case "tui":
	// 	fmt.Println("Starting TUI...")
	// 	tui.Run()
	case "cli":
		fmt.Println("Starting CLI...")
		cli.Execute()
	case "web":
		fmt.Println("Starting web...")
		web.StartServer()
		return

	default:
		fmt.Printf("Unknown input: %s\n", command)
		os.Exit(1)
	}
}
