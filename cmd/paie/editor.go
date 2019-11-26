package main

import "os/exec"

func openEditor(path string) error {
	cmd := exec.Command("code", "-w", path)
	err := cmd.Run()
	return err
}
