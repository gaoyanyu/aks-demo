package util

import (
	"errors"
	"fmt"

	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"k8s.io/klog/v2"
)

func ExecCMD(infoFile, errFile *os.File, command string, args ...string) (int, error, string) {
	cmd := exec.Command(command, args...)
	//设置该cmd在原来的进程组
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	defer func() {
		if infoFile != nil {
			err := infoFile.Close()
			if err != nil {
				klog.Errorf("close info file failed! ", infoFile)
			}
		}

		if errFile != nil {
			err := errFile.Close()
			if err != nil {
				klog.Errorf("close error file failed! ", infoFile)
			}
		}

	}()

	if infoFile != nil {
		cmd.Stdout = infoFile
	}

	if errFile != nil {
		cmd.Stderr = errFile
	}

	if err := cmd.Start(); err != nil {
		return -999, err, ""
	}

	output, _ := ioutil.ReadAll(stdout)
	errput, _ := ioutil.ReadAll(stderr)

	err := make(chan error, 1)
	go func() {
		err <- cmd.Wait()
	}()

	select {
	case er := <-err:
		if er == nil {
			return 0, nil, string(output)
		} else {
			if ex, ok := er.(*exec.ExitError); ok {
				return ex.Sys().(syscall.WaitStatus).ExitStatus(), er, string(errput) //获取命令执行返回状态，相当于shell: echo $?
			}
			return -999, er, string(errput)
		}
	}
}

func EnsureProcessExist(cmd string) (bool, error) {
	res, err, output := ExecCMD(nil, nil, "bash", "-c", cmd)
	if err != nil {
		return false, err
	}
	if res != 0 {
		return false, errors.New(fmt.Sprintf("Exec cmd error, code:%d,err:%+v, output: %s", res, err, output))
	}
	count, err := strconv.Atoi(strings.Trim(strings.Trim(output, " "), "\n"))
	if err != nil {
		return false, err
	}
	return count >= 1, nil
}
