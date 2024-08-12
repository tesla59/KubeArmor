package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

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
