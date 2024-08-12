package main

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var PoliciesDirectory = "policies"

type KubeArmorPolicy struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Severity int      `yaml:"severity"`
		Message  string   `yaml:"message"`
		Tags     []string `yaml:"tags"`
		Selector struct {
			MatchLabels map[string]string `yaml:"matchLabels"`
		} `yaml:"selector"`
		Process struct {
			MatchPaths []struct {
				Path string `yaml:"path"`
			} `yaml:"matchPaths"`
		} `yaml:"process"`
		Action string `yaml:"action"`
	} `yaml:"spec"`
}

func ReadDefaultPolicy() (KubeArmorPolicy, error) {
	data, err := os.ReadFile("policy.yaml")
	if err != nil {
		return KubeArmorPolicy{}, err
	}

	var policy KubeArmorPolicy

	err = yaml.Unmarshal(data, &policy)

	return policy, err
}

func GeneratePoliciesForContainers(containerIDs ...string) ([]KubeArmorPolicy, error) {
	defaultPolicy, err := ReadDefaultPolicy()
	if err != nil {
		return nil, err
	}

	policies := make([]KubeArmorPolicy, 0)

	for _, container := range containerIDs {
		container, err := GetContainerByID(container)
		if err != nil {
			return nil, err
		}
		// container.Name is in the format "/containerName"
		containerName := strings.TrimPrefix(container.Names[0], "/")
		policy := defaultPolicy
		policy.Spec.Selector.MatchLabels["kubearmor.io/container.name"] = containerName
		err = WritePolicyToFile(policy, PoliciesDirectory+"/policy_"+containerName+".yaml")
		if err != nil {
			return nil, err
		}
		policies = append(policies, policy)
	}

	return policies, nil
}

func WritePolicyToFile(policy KubeArmorPolicy, filename string) error {
	data, err := yaml.Marshal(policy)
	if err != nil {
		return err
	}
	directory := filepath.Dir(filename)

	info, err := os.Stat(directory)
	if os.IsNotExist(err) {
		err = os.MkdirAll(directory, 0755)
		if err != nil {
			return err
		}
	} else if !info.IsDir() {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
