// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 Authors of KubeArmor

package feeder

import (
	"fmt"
	cfg "github.com/kubearmor/KubeArmor/KubeArmor/config"
	tp "github.com/kubearmor/KubeArmor/KubeArmor/types"
	"sync"
	"testing"
)

var logger *Feeder

func TestMain(m *testing.M) {
	node := tp.Node{}
	nodeLock := new(sync.RWMutex)

	// load configuration
	if err := cfg.LoadConfig(); err != nil {
		fmt.Println("[FAIL] Failed to load configuration")
		return
	}

	// create logger
	logger = NewFeeder(&node, &nodeLock)
	if logger == nil {
		fmt.Println("[FAIL] Failed to create logger")
		return
	}
	m.Run()
}

func TestFeeder(t *testing.T) {
	// node
	node := tp.Node{}
	nodeLock := new(sync.RWMutex)

	// load configuration
	if err := cfg.LoadConfig(); err != nil {
		t.Log("[FAIL] Failed to load configuration")
		return
	}

	// create logger
	logger := NewFeeder(&node, &nodeLock)
	if logger == nil {
		t.Log("[FAIL] Failed to create logger")
		return
	}
	t.Log("[PASS] Created logger")

	// destroy logger
	if err := logger.DestroyFeeder(); err != nil {
		t.Log("[FAIL] Failed to destroy logger")
		return
	}
	t.Log("[PASS] Destroyed logger")
}

func FuzzFeeder_PushLog(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		log := tp.Log{Message: string(data)}
		logger.PushLog(log)
	})
}

func FuzzU(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		s := string(data)
		_ = s
	})
}
