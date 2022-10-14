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
	fileRepoPath  string
	command       string
	gatewayName   string
	gatewayDomain string
	natsUrl       string
)

const (
	commandValidateConf       = "validate"
	commandGenerateSampleConf = "generate-sample"
	commandGenerateEnvoyConf  = "generate-envoy"
)

type commandHandler func() error

var (
	commandsTable = map[string]commandHandler{
		commandValidateConf: func() error {
			return validateConfigCommand()
		},
		commandGenerateSampleConf: func() error {
			return generateSampleConfCommand()
		},
		commandGenerateEnvoyConf: func() error {
			return generateEnvoyConfigCommand()
		},
	}
)

func init() {
	flag.StringVar(&fileRepoPath, "file", "", "YAML file path for configuration")
	flag.StringVar(&command, "command", commandValidateConf, "Command to invoke")
	flag.StringVar(&gatewayName, "gateway-name", "localhost", "Command to invoke")
	flag.StringVar(&gatewayDomain, "gateway-domain", "localhost", "Command to invoke")
	flag.StringVar(&natsUrl, "nats-url", "tls://nats-server:4222", "NATS URL for messaging")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s Usage:\n", os.Args[0])
		flag.PrintDefaults()

		fmt.Fprintf(os.Stderr, "\nAvailable commands:\n")
		for commandName, _ := range commandsTable {
			fmt.Fprintf(os.Stderr, "\t%s\n", commandName)
		}
	}
}

func main() {
	flag.Parse()
	logger.Init("confli")

	ch := commandsTable[command]
	if ch == nil {
		logger.Fatalf("Unknown command: %s", command)
	}

	err := ch()
	if err != nil {
		logger.Errorf("Command exec returned error: %v", err)
	}
}

func validateConfigCommand() error {
	if utils.IsEmptyString(fileRepoPath) {
		flag.Usage()
		os.Exit(-1)
	}

	_, err := config.NewConfigFileRepository(fileRepoPath, false, false)
	if err != nil {
		logger.Fatalf("Failed to create config repo: %v", err)
	}

	logger.Infof("Config file loaded and validated from: %s", fileRepoPath)
	return nil
}

func generateSampleConfCommand() error {
	return newSampleConfigGenerator(fileRepoPath).generate()
}

func generateEnvoyConfigCommand() error {
	return newEnvoyConfigGenerator(fileRepoPath).generate()
}
