package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"grie-engine/internal/shm"
)

func main() {
	var (
		size   = flag.Int("size", 10*1024*1024, "Size of shared memory (bytes)")
		count  = flag.Int("count", 5, "Number of messages to write")
		delay  = flag.Duration("delay", 100*time.Millisecond, "Delay between writes")
		remove = flag.Bool("remove", false, "Remove file after closing")
	)
	flag.Parse()

	// Create shared memory
	mgr, err := shm.Create(*size)
	if err != nil {
		log.Fatalf("Failed to create shared memory: %v", err)
	}
	defer mgr.Close()

	fmt.Printf("📝 GRIE Writer Started\n")
	fmt.Printf("📂 Shared Memory: %s\n", mgr.Path())
	fmt.Printf("💾 Size: %d bytes\n", mgr.Size())
	fmt.Printf("📤 Messages: %d\n", *count)
	fmt.Printf("\n---\n\n")

	// Write messages
	for i := 1; i <= *count; i++ {
		message := fmt.Sprintf("Message #%d - Timestamp: %s", i, time.Now().Format(time.RFC3339Nano))

		fmt.Printf("[%d/%d] Writing: %q\n", i, *count, message)
		if err := mgr.WriteData([]byte(message)); err != nil {
			log.Fatalf("Write failed: %v", err)
		}

		stats := mgr.GetStats()
		fmt.Printf("  ✅ State: %s | SeqNum: %v | DataLen: %v | Writes: %v | Reads: %v\n",
			stats["state"], stats["seq_num"], stats["data_len"],
			stats["total_writes"], stats["total_reads"])

		if i < *count {
			fmt.Printf("  ⏳ Waiting %v...\n", *delay)
			time.Sleep(*delay)
		}
	}

	fmt.Printf("\n---\n")
	fmt.Printf("✅ Writer finished successfully!\n")
	fmt.Printf("⏱️  Waiting for readers...\n")
	fmt.Printf("📌 Path: %s\n", mgr.Path())
	fmt.Printf("💡 To run reader: go run cmd/reader/main.go %s\n", mgr.Path())
	fmt.Printf("\nPress Ctrl+C to exit and cleanup...\n")

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal
	sig := <-sigChan
	fmt.Printf("\n🛑 Received signal: %v. Shutting down...\n", sig)

	// Cleanup if requested
	if *remove {
		fmt.Printf("🗑️  Removing shared memory...\n")
		if err := os.Remove(mgr.Path()); err != nil {
			log.Printf("⚠️  Failed to remove file: %v", err)
		}
	}

	fmt.Printf("✅ Shutdown complete.\n")
}
