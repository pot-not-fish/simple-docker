package subsystem

import (
	"os"
	"path"
	"strconv"
)

type CPUSetSubSystem struct {
	CpuSet    string
	MountPath string
}

func (s *CPUSetSubSystem) Set(cgroupPath string) error {
	var err error
	s.MountPath, err = FindCgroupMountPoint("cpu", cgroupPath)
	if err != nil {
		return err
	}
	err = os.Mkdir(s.MountPath, 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(path.Join(s.MountPath, "cpu.cfs_quota_us"), []byte(s.CpuSet), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *CPUSetSubSystem) Apply(cgroupPath string, pid int) error {
	err := os.WriteFile(path.Join(s.MountPath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *CPUSetSubSystem) Destroy(cgroupPath string) error {
	return os.RemoveAll(s.MountPath)
}
