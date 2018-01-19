package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
)

func ReapChildren(child *os.Process) error {
	// forward all signals to child
	c := make(chan os.Signal)
	signal.Notify(c)
	go func() {
		for sig := range c {
			child.Signal(sig)
		}
	}()

	glog.V(2).Infof("waiting for child %v to exit; forwarding all signals", child.Pid)

	var wstatus syscall.WaitStatus
	var err error
	for pid := -1; pid != child.Pid; {
		pid, err = syscall.Wait4(-1, &wstatus, 0, nil)
		if err != nil {
			return err
		}
		glog.Infof("reaped pid %v", pid)
	}

	if exitCode := wstatus.ExitStatus(); exitCode != 0 {
		glog.Errorf("child exited with code %v", exitCode)
	} else {
		glog.V(2).Infof("child exited with code %v", exitCode)
	}

	signal.Stop(c)
	close(c)
	return nil
}
