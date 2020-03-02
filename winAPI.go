// +build windows

package main

import (
	"fmt"
	"log"
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
	for _, proc := range processes {
		go func(proc *process.Process) {
			pid := proc.Pid
			handle, err := OpenProcess(int32(pid))
			if handle != 0 {
				for true {
					NtSuspendProcess(handle)
					time.Sleep(time.Duration(100*(100-cpu)) * time.Microsecond)
					NtResumeProcess(handle)
					time.Sleep(time.Duration(100*cpu) * time.Microsecond)
				}
			} else {
				fmt.Println("Error")
				log.Fatal(err)
			}
		}(proc)

	}
}
