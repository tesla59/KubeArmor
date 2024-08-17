package monitor

import (
	"fmt"
	kl "github.com/kubearmor/KubeArmor/KubeArmor/common"
	cfg "github.com/kubearmor/KubeArmor/KubeArmor/config"
	"github.com/kubearmor/KubeArmor/KubeArmor/feeder"
	tp "github.com/kubearmor/KubeArmor/KubeArmor/types"
	"strings"
	"sync"
	"testing"
)

var systemMonitor *SystemMonitor

func TestMain(m *testing.M) {
	// containers
	Containers := map[string]tp.Container{}
	ContainersLock := new(sync.RWMutex)

	// pid map
	ActiveHostPidMap := map[string]tp.PidMap{}
	ActivePidMapLock := new(sync.RWMutex)

	// node
	node := tp.Node{}
	nodeLock := new(sync.RWMutex)

	node.KernelVersion = kl.GetCommandOutputWithoutErr("uname", []string{"-r"})
	node.KernelVersion = strings.TrimSuffix(node.KernelVersion, "\n")

	// load configuration
	if err := cfg.LoadConfig(); err != nil {
		fmt.Println("[FAIL] Failed to load configuration")
		return
	}

	// configuration
	cfg.GlobalCfg.Policy = true
	cfg.GlobalCfg.HostPolicy = true

	// create logger
	logger := feeder.NewFeeder(&node, &nodeLock)
	if logger == nil {
		fmt.Println("[FAIL] Failed to create logger")
		return
	}
	fmt.Println("[PASS] Created logger")

	// montor lock
	monitorLock := new(sync.RWMutex)

	// Create System Monitor
	systemMonitor = NewSystemMonitor(&node, &nodeLock, logger, &Containers, &ContainersLock, &ActiveHostPidMap, &ActivePidMapLock, &monitorLock)
	if systemMonitor == nil {
		fmt.Println("[FAIL] Failed to create SystemMonitor")

		if err := logger.DestroyFeeder(); err != nil {
			fmt.Println("[FAIL] Failed to destroy logger")
			return
		}
		return
	}
	fmt.Println("[PASS] Created SystemMonitor")
	go systemMonitor.UpdateLogs()
	m.Run()
}

func FuzzUpdateLogs(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		con := ContextCombined{
			ContainerID: string(data),
			ContextSys:  SyscallContext{},
			ContextArgs: nil,
		}
		systemMonitor.ContextChan <- con
	})
}
