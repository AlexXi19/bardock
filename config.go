package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Services     map[string]Service `yaml:"services"`
	GlobalConfig GlobalConfig       `yaml:"global_config"`
}

type Service struct {
	Dockerfile   string `yaml:"dockerfile"`
	Image        string `yaml:"image_name"`
	BuildContext string `yaml:"build_context"`
	Push         *bool  `yaml:"push"`
}

type GlobalConfig struct {
	RegistryURL string `yaml:"registry_url"`
	ImageTag    string `yaml:"image_tag"`
}

// constant
var (
	defaultBardockPath = "bardock.yaml"
	validImageTags     = []string{"latest", "git_sha"}
)

func validateConfig(config *Config) error {
	logger.Debugf("validateConfig: %v", config)
	allErrors := []error{}

	imageTag := cliConfig.ImageTag
	if imageTag != "" {
		config.GlobalConfig.ImageTag = imageTag
	} else {
		if config.GlobalConfig.ImageTag == "" {
			config.GlobalConfig.ImageTag = "latest"
		} else {
			match := false
			for _, validImageTag := range validImageTags {
				if config.GlobalConfig.ImageTag == validImageTag {
					match = true
				}
			}

			if !match {
				allErrors = append(allErrors, fmt.Errorf("image_tag must be one of %v", validImageTags))
			}
		}
	}

	if config.GlobalConfig.ImageTag == "git_sha" {
		cmd := exec.Command("git", "rev-parse", "HEAD")
		output, err := cmd.Output()
		// TODO: Get stdout
		if err != nil {
			return fmt.Errorf("error getting git sha")
		}
		config.GlobalConfig.ImageTag = string(output)[:7]
	}

	if config.GlobalConfig.RegistryURL == "" {
		allErrors = append(allErrors, fmt.Errorf("registry_url is required"))
	}

	for serviceName, service := range config.Services {
		if service.Image == "" {
			allErrors = append(allErrors, fmt.Errorf("services.%s.image is required", serviceName))
		}
		if service.Dockerfile == "" {
			allErrors = append(allErrors, fmt.Errorf("services.%s.dockerfile is required", serviceName))
		}
	}

	combinedErrors := ""
	for _, err := range allErrors {
		combinedErrors += err.Error() + "\n"
	}

	if combinedErrors != "" {
		return fmt.Errorf(combinedErrors)
	}

	return nil
}

func parseConfig(filePath string) (Config, error) {
	logger.Debugf("parseConfig: %s", filePath)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
