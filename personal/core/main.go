package main

import (
	"flag"
	"fmt"
	"os"
	"wacast/core/utils"
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
