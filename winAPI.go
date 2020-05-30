// +build windows

package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/process"
)

// BOOL WINAPI NtSuspendProcess
//   _In_ HANDLE hProcess,
func NtSuspendProcess(hProcess syscall.Handle) bool {
	ntdll, _ := syscall.LoadLibrary("ntdll.dll")
	ntSuspendProcess, _ := syscall.GetProcAddress(ntdll, "NtSuspendProcess")
	ret, _, _ := syscall.Syscall(ntSuspendProcess, 1,
		uintptr(hProcess),
		0,
		0)
	return ret == 0
}

// BOOL WINAPI NtResumeProcess
//   _In_ HANDLE hProcess,
func NtResumeProcess(hProcess syscall.Handle) bool {
	ntdll, _ := syscall.LoadLibrary("ntdll.dll")
	ntResumeProcess, _ := syscall.GetProcAddress(ntdll, "NtResumeProcess")
	ret, _, _ := syscall.Syscall(ntResumeProcess, 1,
		uintptr(hProcess),
		0,
		0)
	return ret == 0
}

func OpenProcess(processId int32) (h syscall.Handle, e error) {
	_PROCESS_ALL_ACCESS := syscall.STANDARD_RIGHTS_REQUIRED | syscall.SYNCHRONIZE | 0xfff
	ph, err := syscall.OpenProcess(uint32(_PROCESS_ALL_ACCESS), false, uint32(processId))
	if err != nil {
		log.Fatal(err)
	}
	return ph, err
}

func throttle(processes []*process.Process, cpu int) {
	gather := make(chan bool, len(processes))
	for _, proc := range processes {
		stopChan := make(chan bool, 1)
		go func(proc *process.Process, stopChannel chan bool) {
			pid := proc.Pid
			handle, err := OpenProcess(int32(pid))
			for {
				select {
				case <-stopChannel:
					break
				default:
					if handle != 0 {
						NtSuspendProcess(handle)
						time.Sleep(time.Duration(100*(100-cpu)) * time.Microsecond)
						NtResumeProcess(handle)
						time.Sleep(time.Duration(100*cpu) * time.Microsecond)
					} else {
						fmt.Println("Handle error")
						log.Fatal(err)
					}
				}
			}
		}(proc, stopChan)

		go func(proc *process.Process, stopChannel chan bool) {
			for true {
				time.Sleep(time.Duration(2) * time.Second)
				_, err := OpenProcess(int32(proc.Pid))
				if err != nil {
					stopChannel <- true
					gather <- true
				}
			}
		}(proc, stopChan)
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
		go func(proc *process.Process) {
			pid := proc.Pid
			handle, _ := OpenProcess(int32(pid))
			if handle != 0 {
				NtResumeProcess(handle)
			}
		}(proc)
	}
}

func suspendProcesses(processes []*process.Process) {
	for _, proc := range processes {
		go func(proc *process.Process) {
			pid := proc.Pid
			handle, _ := OpenProcess(int32(pid))
			if handle != 0 {
				NtSuspendProcess(handle)
			}
		}(proc)
	}
}
