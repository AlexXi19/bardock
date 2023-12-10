package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func executeConfig(services []string, config Config) error {
	logger.Debugf("Building services: %v", services)
	serviceConfigs := make(map[string]Service)
	for _, serviceName := range services {
		serviceConfigs[serviceName] = config.Services[serviceName]
	}

	for serviceName, service := range serviceConfigs {
		logger.Infof("Building service: %s", serviceName)
		err := executeServiceConfig(service, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func executeServiceConfig(service Service, config Config) error {
	bardockPath := cliConfig.FilePath
	relativeBuildContext := service.BuildContext
	if relativeBuildContext == "" {
		relativeBuildContext = "."
	}
	buildContext := filepath.Join(filepath.Dir(bardockPath), relativeBuildContext)
	cmdStr := fmt.Sprintf("cd %s", buildContext)

	relativeDockerfilePath, err := filepath.Rel(service.BuildContext, service.Dockerfile)
	if err != nil {
		return err
	}
	imageTag := fmt.Sprintf("%s/%s:%s", config.GlobalConfig.RegistryURL, service.Image, config.GlobalConfig.ImageTag)

	cmdStr += fmt.Sprintf(" && docker build -t %s -f %s .", imageTag, relativeDockerfilePath)
	logger.Debugf("Build command: %s", cmdStr)

	var push bool
	if (service.Push != nil) && !*service.Push {
		push = false
	} else {
		push = true
	}

	if push {
		cmdStr += fmt.Sprintf(" && docker push %s", imageTag)
	}

	// Execute
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
