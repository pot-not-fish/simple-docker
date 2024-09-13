package subsystem

import (
	"os"
	"path"
	"strconv"
)

type CPUSubSystem struct {
	CpuLimmit string
	MountPath string
}

func (s *CPUSubSystem) Set(cgroupPath string) error {
	var err error
	s.MountPath, err = FindCgroupMountPoint("cpu", cgroupPath)
	if err != nil {
		return err
	}

	err = os.Mkdir(s.MountPath, 0644)
	if err != nil {
		return err
	}

	// 默认值为100000 相当于100ms
	err = os.WriteFile(path.Join(s.MountPath, "cpu.cfs_period_us"), []byte("100000"), 0644)
	if err != nil {
		return err
	}

	percent, err := strconv.Atoi(s.CpuLimmit)
	if err != nil {
		return err
	}
	cpuLimit := strconv.Itoa(percent * 1000)

	err = os.WriteFile(path.Join(s.MountPath, "cpu.cfs_quota_us"), []byte(cpuLimit), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *CPUSubSystem) Apply(cgroupPath string, pid int) error {
	err := os.WriteFile(path.Join(s.MountPath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *CPUSubSystem) Destroy(cgroupPath string) error {
	return os.RemoveAll(s.MountPath)
}
