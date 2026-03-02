package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"grie-engine/internal/shm"
)

func main() {
	var (
		timeout = flag.Int("timeout", 5000, "Read timeout in milliseconds")
		count   = flag.Int("count", 0, "Number of messages to read (0 = unlimited)")
	)
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <shm_path> [-timeout ms] [-count n]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s /tmp/grie_shm_123456789\n", os.Args[0])
		os.Exit(1)
	}

	path := args[0]

	// Open shared memory
	mgr, err := shm.Open(path)
	if err != nil {
		log.Fatalf("Failed to open shared memory: %v", err)
	}
	defer mgr.Close()

	fmt.Printf("📖 GRIE Reader Started\n")
	fmt.Printf("📂 Shared Memory: %s\n", path)
	fmt.Printf("💾 Size: %d bytes\n", mgr.Size())
	fmt.Printf("⏱️  Timeout: %dms\n", *timeout)
	if *count > 0 {
		fmt.Printf("📨 Max reads: %d\n", *count)
	}
	fmt.Printf("\n---\n\n")

	readCount := 0
	for {
		if *count > 0 && readCount >= *count {
			break
		}

		readCount++
		fmt.Printf("[%d] Reading... ", readCount)

		data, err := mgr.ReadData(*timeout)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			break
		}

		fmt.Printf("✅\n")
		fmt.Printf("  📝 Data: %q\n", string(data))

		stats := mgr.GetStats()
		fmt.Printf("  📊 State: %s | SeqNum: %v | DataLen: %v | Reads: %v | Writes: %v\n",
			stats["state"], stats["seq_num"], stats["data_len"],
			stats["total_reads"], stats["total_writes"])
	}

	fmt.Printf("\n---\n")
	fmt.Printf("✅ Reader finished successfully! (%d messages read)\n", readCount)
}
