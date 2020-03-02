// +build !windows

package main

import (
	"time"

	"github.com/shirou/gopsutil/process"
)

func throttle(processes []*process.Process, cpu int) {
	for _, p := range processes {
		go func(proc *process.Process) {
			for true {
				proc.Suspend()
				time.Sleep(time.Duration(10*(100-cpu)) * time.Microsecond)
				proc.Resume()
				time.Sleep(time.Duration(10*cpu) * time.Microsecond)
			}
		}(p)
	}
}
