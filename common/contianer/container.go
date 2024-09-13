package container

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File, error) {
	log.Println("start new parent process")

	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	// /proc/self/exe即程序本身
	// 调用程序本身进行初始化
	cmd := exec.Command("/proc/self/exe", "init")
	// 先clone出新的隔离环境，然后再使用cmd start执行init操作
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	// cmd执行的时候就会带着readPipe去执行
	cmd.ExtraFiles = []*os.File{readPipe}

	return cmd, writePipe, nil
}

func RunContainerInitProcess(command string, args []string) error {
	log.Println("start run container init process")

	var err error
	// 挂载进程proc信息
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		return err
	}

	// 获取父进程的管道信息
	cmds, err := ReadUserCommand()
	if err != nil {
		return err
	}
	path, err := exec.LookPath(cmds[0])
	if err != nil {
		return err
	}

	// 覆盖旧的进程的镜像、数据和堆栈，PID等信息
	if err = syscall.Exec(path, cmds[0:], os.Environ()); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// index0 标准输入
// index1 标准输出
// index2 标准错误
// index3 带过来的第一个FD，也就是readPipe
const fdIndex = 3

func ReadUserCommand() ([]string, error) {
	pipe := os.NewFile(uintptr(fdIndex), "pipe")
	msg, err := io.ReadAll(pipe)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(msg), " "), nil
}
