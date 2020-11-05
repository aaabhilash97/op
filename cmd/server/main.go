package main

import (
	"fmt"
	"os"

	"github.com/aaabhilash97/op/internal/appconfig"
	"github.com/aaabhilash97/op/pkg/cmd"
	"github.com/aaabhilash97/op/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	if err := cmd.RunServer(appconfig.Appconfig); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		logger.Error("Runserver", zap.Error(err))
		os.Exit(1)
	}
}
