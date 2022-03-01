package jobobject

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")

	// modKernelbase = windows.NewLazySystemDLL("kernelbase.dll")

	procIsProcessInJob = modkernel32.NewProc("IsProcessInJob")
	procOpenJobObjectW = modkernel32.NewProc("OpenJobObjectW")
	// procIsProcessInJob = modkernel32.NewProc("IsProcessInJob")
	// procOpenJobObjectW = modkernel32.NewProc("OpenJobObjectW")
)

func IsProcessInJob(pHandle windows.Handle, jHandle windows.Handle, result *bool) (err error) {
	r1, _, e1 := syscall.Syscall(procIsProcessInJob.Addr(), 3, uintptr(pHandle), uintptr(jHandle), uintptr(unsafe.Pointer(result)))
	if r1 == 0 {
		if e1 != 0 {
			err = e1
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func OpenJobObjectW(desiredAccess uint32, inheritHandle bool, lpName *uint16) (handle windows.Handle, err error) {
	var p0 uint32
	if inheritHandle {
		p0 = 1
	} else {
		p0 = 0
	}
	r0, _, e1 := syscall.Syscall(procOpenJobObjectW.Addr(), 3, uintptr(desiredAccess), uintptr(p0), uintptr(unsafe.Pointer(lpName)))
	fmt.Printf("--- r0: %v\n", r0)
	handle = windows.Handle(r0)
	if handle == 0 {
		if e1 != 0 {
			err = e1
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
