package main

import (
	"fmt"
	"log/slog"

	"github.com/mrhinton101/fluyt/cmd"
	"github.com/mrhinton101/fluyt/logger"
)

func main() {
	logfile := logger.InitLogger("fluyt.json")
	defer logfile.Close()
	defer func() {
		if r := recover(); r != nil {

			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       fmt.Errorf("%v", r),
				Component: "runtime",
				Action:    "panic recovery",
				Msg:       "unexpected panic has occured",
				Target:    "localhost",
			})
		}
	}()

	cmd.Execute()

	panic("test")

}
