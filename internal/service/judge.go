package service

import (
	"fmt"
	"os"
	"os/exec"
)

func Execute(code string) {

	tmpFile, err := os.CreateTemp("", "*.py")
	if err != nil {
		fmt.Printf("CreateTemp %s\n", err.Error())
		panic(err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(code)
	if err != nil {
		fmt.Printf("WriteString %s\n", err.Error())
		panic(err)
	}
	tmpFile.Close()

	cmd := exec.Command("python3", tmpFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("runErr %s\n", err.Error())
		return
	}
}
