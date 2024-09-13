package subsystem

import (
	"fmt"
	"os"
	"testing"
)

func TestRemove(t *testing.T) {
	err := os.Remove("/sys/fs/cgroup/memory/simple_docker")
	if err != nil {
		fmt.Println(err)
	}
}
