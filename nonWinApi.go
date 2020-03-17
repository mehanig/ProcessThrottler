// +build !windows

package main

import (
	"os"
	"time"

	"github.com/shirou/gopsutil/process"
)

func throttle(processes []*process.Process, cpu int) {
	gather := make(chan bool, len(processes))
	for _, p := range processes {
		stopChan := make(chan bool, 1)
		go func(proc *process.Process, stopChannel chan bool) {
			for {
				select {
				case <-stopChannel:
					break
				default:
					proc.Suspend()
					time.Sleep(time.Duration(10*(100-cpu)) * time.Microsecond)
					proc.Resume()
					time.Sleep(time.Duration(10*cpu) * time.Microsecond)
				}
			}
		}(p, stopChan)

		go func(proc *process.Process, stopChannel chan bool) {
			for true {
				time.Sleep(time.Duration(2) * time.Second)
				_, err := proc.Status()
				if err != nil {
					stopChannel <- true
					gather <- true
				}
			}
		}(p, stopChan)
	}

	for range processes {
		select {
		case <-gather:
		}
	}
	os.Exit(0)
}

func resumeSuspended(processes []*process.Process) {
	for _, proc := range processes {
		proc.Resume()
	}
}
