package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/logger"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
)

var (
	fileRepoPath string
	command      string
)

const (
	commandValidateConf = "validate"
)

func init() {
	flag.StringVar(&fileRepoPath, "file", "", "YAML file path for configuration")
	flag.StringVar(&command, "command", commandValidateConf, "Command to invoke")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s Usage:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	logger.Init("confli")

	if utils.IsEmptyString(fileRepoPath) {
		flag.Usage()
		os.Exit(-1)
	}

	switch command {
	case commandValidateConf:
		validateConfig()
		break
	default:
		logger.Errorf("Unknown command: %s", command)
	}
}

func validateConfig() {
	_, err := config.NewConfigFileRepository(fileRepoPath)
	if err != nil {
		logger.Fatalf("Failed to create config repo: %v", err)
	}

	logger.Infof("Config file loaded and validated from: %s", fileRepoPath)
}
