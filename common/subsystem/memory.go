package subsystem

import (
	"os"
	"path"
	"strconv"
)

type MemorySubSystem struct {
	MemoryLimit string
	MountPath   string
}

func (s *MemorySubSystem) Set(cgroupPath string) error {
	var err error
	s.MountPath, err = FindCgroupMountPoint("memory", cgroupPath)
	if err != nil {
		return err
	}

	err = os.Mkdir(s.MountPath, 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(path.Join(s.MountPath, "memory.limit_in_bytes"), []byte(s.MemoryLimit), 0644)
	if err != nil {
		return err
	}
	// 限制内存和磁盘空间交换
	// 没有这条无法限制最大的内存分配，因为多余的内存会交换到磁盘
	err = os.WriteFile(path.Join(s.MountPath, "memory.memsw.limit_in_bytes"), []byte(s.MemoryLimit), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	err := os.WriteFile(path.Join(s.MountPath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *MemorySubSystem) Destroy(cgroupPath string) error {
	return os.RemoveAll(s.MountPath)
}
