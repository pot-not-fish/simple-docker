package subsystem

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

type Subsystem interface {
	Set(cgroupPath string) error
	Apply(cgroupPath string, pid int) error
	Destroy(cgroupPath string) error
}

type ControlSubsystem struct {
	Subsystems []Subsystem
	CgroupPath string
	PID        int
}

func NewControlSubsystem(pid int, cgroupPath string) *ControlSubsystem {
	return &ControlSubsystem{
		CgroupPath: cgroupPath,
		PID:        pid,
	}
}

func (c *ControlSubsystem) Register(subsystem Subsystem) {
	c.Subsystems = append(c.Subsystems, subsystem)
}

func (c *ControlSubsystem) SetAll() {
	for _, v := range c.Subsystems {
		err := v.Set(c.CgroupPath)
		if err != nil {
			log.Println(err)
		}
	}
}

func (c *ControlSubsystem) ApplyAll() {
	for _, v := range c.Subsystems {
		err := v.Apply(c.CgroupPath, c.PID)
		if err != nil {
			log.Println(err)
		}
	}
}

func (c *ControlSubsystem) DestroyAll() {
	for _, v := range c.Subsystems {
		err := v.Destroy(c.CgroupPath)
		if err != nil {
			log.Println(err)
		}
	}
}

// MemoryLimit 100  byte
// CpuLimit    20   %
// CpuSet      1,2

// 查找cgroup的位置，如果没有找到，需要自行创建
func FindCgroupMountPoint(subsystem string, cgroupPath string) (string, error) {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return "", err
	}
	defer f.Close()

	// 33 25 0:29 / /sys/fs/cgroup/memory rw,nosuid,nodev,noexec,relatime shared:15 - cgroup cgroup rw,seclabel,memory
	// 28 25 0:24 / /sys/fs/cgroup/pids rw,nosuid,nodev,noexec,relatime shared:10 - cgroup cgroup rw,seclabel,pids
	// 29 25 0:25 / /sys/fs/cgroup/cpuset rw,nosuid,nodev,noexec,relatime shared:11 - cgroup cgroup rw,seclabel,cpuset
	// 这个cpu,cpuacct是一个文件名
	// 30 25 0:26 / /sys/fs/cgroup/cpu,cpuacct rw,nosuid,nodev,noexec,relatime shared:12 - cgroup cgroup rw,seclabel,cpuacct,cpu
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, " ")
		// rw,seclabel,memory
		subsystems := strings.Split(fields[len(fields)-1], ",")
		for _, opt := range subsystems {
			if opt == subsystem {
				// /sys/fs/cgroup/memory
				return path.Join(fields[4], cgroupPath), nil
			}
		}
	}
	if err = scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("could not found cgroup path")
}
