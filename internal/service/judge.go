package service

import (
	"bytes"
	"code-judge/internal/transport/http/model"
	"code-judge/pkg/logger"
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func InitSubmission(ctx context.Context, code string) (*model.ResponseModel, error) {
	log := logger.Ctx(ctx)

	tmpFile, err := os.CreateTemp("", "*.py")
	if err != nil {
		log.Error("error creating temp file: ", err)
		panic(err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(code)
	if err != nil {
		log.Error("error writing to temp file: ", err)
		return nil, err
	}
	tmpFile.Close()

	syncR, syncW, _ := os.Pipe()

	cmd := exec.CommandContext(ctx, "/proc/self/exe", "python3", tmpFile.Name())
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWCGROUP | syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC | syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWPID | syscall.CLONE_NEWNET | syscall.CLONE_NEWNS}
	cmd.Env = append(os.Environ(), "SANDBOX_ROLE=child")
	cmd.ExtraFiles = []*os.File{syncR}

	var outputBuf, errBuf bytes.Buffer
	cmd.Stdout = &outputBuf
	cmd.Stderr = &errBuf

	ticker := time.NewTicker(100 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				log.Debug("cmd canceled by timeout")
				cmd.Cancel()
			}
		}
	}()

	if err := cmd.Start(); err != nil {
		ticker.Stop()
		log.Error(fmt.Sprintf("error running command: %s", err))
		log.Debug(fmt.Sprintf("error output: %s", errBuf.String()))
		return nil, err
	}
	syncR.Close()
	childPid := cmd.Process.Pid

	uid := os.Getuid()
	gid := os.Getgid()

	os.WriteFile(fmt.Sprintf("/proc/%d/uid_map", childPid),
		[]byte(fmt.Sprintf("0 %d 1", uid)), 0)

	os.WriteFile(fmt.Sprintf("/proc/%d/setgroups", childPid),
		[]byte("deny"), 0)

	os.WriteFile(fmt.Sprintf("/proc/%d/gid_map", childPid),
		[]byte(fmt.Sprintf("0 %d 1", gid)), 0)

	syncW.Close()
	if err := cmd.Wait(); err != nil {
		ticker.Stop()
		log.Error(fmt.Sprintf("error running command: %s", err))
		log.Debug(fmt.Sprintf("error output: %s", errBuf.String()))
		return nil, err
	}

	ticker.Stop()

	log.Debug(fmt.Sprintf("output: %s", outputBuf.String()))
	log.Debug(fmt.Sprintf("errBuf: %s", errBuf.String()))

	return &model.ResponseModel{Output: outputBuf.String(), Error: errBuf.String()}, nil
}

func ExecuteSubmission(lang, filename string) (string, error) {

	cmd := exec.Command(lang, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return "", nil
}
