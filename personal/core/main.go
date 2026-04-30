package main

import (
	"flag"
	"fmt"
	"os"
	"unsafe"
	"wacast/core/utils"
	"runtime"
	"syscall"
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procCreateMutex  = kernel32.NewProc("CreateMutexW")
	procGetLastError = kernel32.NewProc("GetLastError")
)

const (
	ERROR_ALREADY_EXISTS = 183
)

const Version = "1.0.0"

func main() {
	// 1. Flags
	headless := flag.Bool("headless", false, "Run without GUI (CLI mode)")
	flag.Parse()

	if *headless {
		runHeadless()
		return
	}

	// 1.5 Single Instance Protection (Windows only)
	if runtime.GOOS == "windows" {
		mutexName, _ := syscall.UTF16PtrFromString("WACAST_PERSONAL_MUTEX")
		_, _, _ = procCreateMutex.Call(0, 1, uintptr(unsafe.Pointer(mutexName)))
		ret, _, _ := procGetLastError.Call()
		if ret == uintptr(ERROR_ALREADY_EXISTS) {
			// Instead of fmt, we might want to show a message box, but for now exit
			fmt.Println("WACAST is already running.")
			os.Exit(0)
		}
	}

	// 2. Initialize Logger early
	_ = utils.InitLogger("info")

	// 3. GUI Mode (Default for Personal)
	fmt.Println("Memulai Control Panel...")
	cp := NewControlPanel()
	cp.Run()
}

func runHeadless() {
	app := NewApp()
	if err := app.Start(); err != nil {
		fmt.Printf("FATAL ERROR: %v\n", err)
		fmt.Println("\nTekan Enter untuk keluar...")
		fmt.Scanln()
		os.Exit(1)
	}

	// Wait for signal in headless mode
	utils.Info("Headless mode active. Press Ctrl+C to stop.")
	
	// Keep running until interrupted
	select {}
}
