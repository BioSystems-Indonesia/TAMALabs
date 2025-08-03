//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

// showErrorMessage displays a native Windows error message box.
func showErrorMessage(title, message string) {
	user32 := syscall.NewLazyDLL("user32.dll")
	messageBoxW := user32.NewProc("MessageBoxW")
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	msgPtr, _ := syscall.UTF16PtrFromString(message)
	// MB_ICONERROR = 0x10, MB_OK = 0x0
	messageBoxW.Call(0, uintptr(unsafe.Pointer(msgPtr)), uintptr(unsafe.Pointer(titlePtr)), 0x10)
}
