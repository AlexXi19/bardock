package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

var (
	version   = "0.1.0"
	logger    = logrus.New()
	cliConfig = CommandLineConfig{}
)

type CommandLineConfig struct {
	FilePath string
	ImageTag string
	Push     bool
}

func main() {
	logger = &logrus.Logger{
		Out:   os.Stdout,
		Level: logrus.InfoLevel,
		Formatter: &easy.Formatter{
			LogFormat: "[%lvl%] %msg% \n",
		},
	}

	var showVersion bool
	var verbose bool
	var rootCmd = &cobra.Command{
		Use:   "bardock",
		Short: "Bardock: A monorepo Dockerfile management tool",
		PreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				logger.SetLevel(logrus.DebugLevel)
			}
			if showVersion {
				fmt.Printf("Bardock version %s\n", version)
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
			services := args
			err := run(services, cliConfig)
			if err != nil {
				logger.Error("%v", err)
				os.Exit(1)
			}

			fmt.Println("\nBardock done!")
			os.Exit(0)
		},
	}

	rootCmd.Flags().BoolVar(&showVersion, "version", false, "print the version number of bardock")
	rootCmd.Flags().StringVarP(&cliConfig.FilePath, "file", "f", "bardock.yaml", "path to the bardock YAML file")
	rootCmd.Flags().StringVarP(&cliConfig.ImageTag, "tag", "t", "", "override image tag")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.Flags().BoolVarP(&cliConfig.Push, "push", "p", false, "override push image to registry")

	rootCmd.Execute()
}

func run(services []string, cliConfig CommandLineConfig) error {
	config, err := parseConfig(cliConfig.FilePath)
	if err != nil {
		return err
	}

	err = validateConfig(&config)
	if err != nil {
		return err
	}

	err = executeConfig(services, config)
	if err != nil {
		return err
	}

	return nil
}
