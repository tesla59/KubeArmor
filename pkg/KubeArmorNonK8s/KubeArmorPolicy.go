package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

// ReadDefaultPolicy reads the default policy from the policy.yaml file and returns KubeArmorPolicy{} and error
func ReadDefaultPolicy() (KubeArmorPolicy, error) {
	data, err := os.ReadFile("policy.yaml")
	if err != nil {
		return KubeArmorPolicy{}, err
	}

	var policy KubeArmorPolicy

	err = yaml.Unmarshal(data, &policy)

	return policy, err
}

// GeneratePoliciesForContainers generates policies for the specified container IDs. Uses the default policy as a template
func GeneratePoliciesForContainers(ctx context.Context, containerIDs ...string) ([]KubeArmorPolicy, error) {
	defaultPolicy, err := ReadDefaultPolicy()
	if err != nil {
		return nil, err
	}

	policies := make([]KubeArmorPolicy, 0)

	for _, container := range containerIDs {
		container, err := GetContainerByID(ctx, container)
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

// WritePolicyToFile writes the policy to the specified file. Called by GeneratePoliciesForContainers
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

// ApplyPolicy applies the policy to the KubeArmor VM
func ApplyPolicy(ctx context.Context, policy string) error {
	info, err := os.Stat(policy)
	if os.IsNotExist(err) {
		return err
	} else if info.IsDir() {
		return err
	}
	karmor := exec.CommandContext(ctx, "karmor", "vm", "policy", "add", policy)
	karmor.Stdout = os.Stdout
	karmor.Stderr = os.Stderr
	return karmor.Run()
}

// ApplyPolicies applies the policies in the specified directory to the KubeArmor VM
// If the path is a directory, it applies all the policies in the directory
// If the path is a file, it applies the policy in the file
func ApplyPolicies(ctx context.Context, policyPath string) error {
	info, err := os.Stat(policyPath)
	if os.IsNotExist(err) {
		return err
	}
	if info.IsDir() {
		var policies []string
		err := filepath.Walk(policyPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				policies = append(policies, path)
			}
			return nil
		})
		if err != nil {
			return err
		}
		for _, policy := range policies {
			err := ApplyPolicy(ctx, policy)
			if err != nil {
				return err
			}
		}
	} else {
		err := ApplyPolicy(ctx, policyPath)
		if err != nil {
			return err
		}
	}
	return nil
}
