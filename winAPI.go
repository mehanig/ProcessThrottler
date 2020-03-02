// +build windows

package main

import (
	"log"
	"syscall"
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

func OpenProcess(processId int) (h syscall.Handle, e error) {
	_PROCESS_ALL_ACCESS := syscall.STANDARD_RIGHTS_REQUIRED | syscall.SYNCHRONIZE | 0xfff
	ph, err := syscall.OpenProcess(uint32(_PROCESS_ALL_ACCESS), false, uint32(processId))
	if err != nil {
		log.Fatal(err)
	}
	return ph, err
}
