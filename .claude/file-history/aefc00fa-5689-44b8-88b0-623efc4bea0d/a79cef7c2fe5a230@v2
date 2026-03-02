// GRIE Julia Simulator
// Simulates Julia reader.jl behavior for testing without Julia runtime
package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const (
	SHM_SIZE     = 10 * 1024 * 1024 // 10MB
	HEADER_SIZE  = 128
	STATE_IDLE   = int32(0)
	STATE_WRITING = int32(1)
	STATE_READY  = int32(2)
	STATE_READING = int32(3)
)

// EngineHeader matches Go header structure
type EngineHeader struct {
	Magic      uint64
	Version    uint32
	State      int32
	SeqNum     uint64
	Timestamp  int64
	DataLen    uint64
	_pad1      [24]byte
	WriterPID  int64
	ReaderPID  int64
	TotalWrites uint64
	TotalReads uint64
	_pad2      [32]byte
}

func readHeader(data []byte) *EngineHeader {
	header := (*EngineHeader)(unsafe.Pointer(&data[0]))
	return header
}

func stateName(state int32) string {
	switch state {
	case STATE_IDLE:
		return "IDLE"
	case STATE_WRITING:
		return "WRITING"
	case STATE_READY:
		return "READY"
	case STATE_READING:
		return "READING"
	default:
		return "UNKNOWN"
	}
}

// MatMul4x4 - Matrix multiplication 4×4
func matmul4x4(a, b *[4][4]float64) *[4][4]float64 {
	c := &[4][4]float64{}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				c[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return c
}

// ParseMatrixData - Parse binary matrix data from SHM
func parseMatrixData(data []byte) (*[4][4]float64, error) {
	if len(data) < 128 { // 4*4*8 = 128 bytes
		return nil, fmt.Errorf("insufficient data for 4×4 matrix")
	}

	matrix := &[4][4]float64{}
	offset := 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			bits := binary.LittleEndian.Uint64(data[offset : offset+8])
			matrix[i][j] = math.Float64frombits(bits)
			offset += 8
		}
	}
	return matrix, nil
}

// PerformMatMul - Simulate Julia MatMul operation
func performMatMul(shmData []byte, dataLen uint64) time.Duration {
	dataBytes := shmData[HEADER_SIZE : HEADER_SIZE+dataLen]

	matrixA, err := parseMatrixData(dataBytes[:128])
	if err != nil {
		fmt.Printf("  ❌ MatMul Error: %v\n", err)
		return 0
	}

	matrixB, err := parseMatrixData(dataBytes[128:])
	if err != nil {
		fmt.Printf("  ❌ MatMul Error: %v\n", err)
		return 0
	}

	start := time.Now()
	result := matmul4x4(matrixA, matrixB)
	elapsed := time.Since(start)

	fmt.Printf("  🎯 MatMul 4×4 completed in %v\n", elapsed)
	fmt.Printf("     Result[0][0] = %.6f\n", result[0][0])

	return elapsed
}

// FFTSimple - Simple FFT simulation (just log the operation)
func performFFT(shmData []byte, dataLen uint64) time.Duration {
	dataBytes := shmData[HEADER_SIZE : HEADER_SIZE+dataLen]

	if len(dataBytes) < 16*2*8 { // 16 complex numbers = 16*2*8 bytes
		fmt.Printf("  ❌ Insufficient data for FFT-16\n")
		return 0
	}

	start := time.Now()
	// Simulated FFT: just read the data
	sum := 0.0
	for i := 0; i < 16; i++ {
		bits := binary.LittleEndian.Uint64(dataBytes[i*16 : i*16+8])
		sum += math.Float64frombits(bits)
	}
	elapsed := time.Since(start)

	fmt.Printf("  🎯 FFT-16 completed in %v\n", elapsed)
	fmt.Printf("     Data sum = %.6f\n", sum)

	return elapsed
}

func main() {
	path := flag.String("path", "", "Path to SHM file")
	timeout := flag.Duration("timeout", 5*time.Second, "Timeout for polling")
	maxReads := flag.Int("count", 3, "Max number of reads")
	flag.Parse()

	if *path == "" {
		fmt.Println("Usage: julia_simulator -path <shm_path> [-timeout 5s] [-count 3]")
		fmt.Println("Example: julia_simulator -path /tmp/grie_shm_123456789")
		os.Exit(1)
	}

	// Open SHM file
	fd, err := syscall.Open(*path, syscall.O_RDWR, 0)
	if err != nil {
		fmt.Printf("❌ Error opening SHM: %v\n", err)
		os.Exit(1)
	}
	defer syscall.Close(fd)

	// mmap the file
	shmData, err := syscall.Mmap(fd, 0, SHM_SIZE, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		fmt.Printf("❌ Error mmap: %v\n", err)
		os.Exit(1)
	}
	defer syscall.Munmap(shmData)

	fmt.Println("📖 GRIE Julia Simulator Started")
	fmt.Printf("📂 Shared Memory: %s\n", *path)
	fmt.Printf("💾 Size: %.1f MB\n", float64(SHM_SIZE)/1024/1024)
	fmt.Println()
	fmt.Println("---\n")

	readCount := 0
	startTime := time.Now()
	deadline := startTime.Add(*timeout)

	for readCount < *maxReads && time.Now().Before(deadline) {
		header := readHeader(shmData)

		if header.State == STATE_READY {
			readCount++
			fmt.Printf("[%d/%d] Reading...\n", readCount, *maxReads)
			fmt.Printf("  📝 SeqNum: %d | DataLen: %d | State: %s\n",
				header.SeqNum, header.DataLen, stateName(header.State))

			// Perform computation based on data type
			// Check if it's matrix data (256 bytes = 2×4×4 matrices)
			if header.DataLen >= 256 {
				performMatMul(shmData, 256)
			} else if header.DataLen >= 128 {
				performFFT(shmData, header.DataLen)
			}

			fmt.Printf("  📈 TotalWrites: %d | TotalReads: %d\n",
				header.TotalWrites, header.TotalReads)

			// Transition READY → IDLE (atomic store in real scenario)
			header.State = STATE_IDLE

			fmt.Println()
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}

	fmt.Println("---")
	if readCount >= *maxReads {
		fmt.Printf("✅ Simulator finished successfully! (%d messages read)\n", readCount)
	} else {
		fmt.Printf("⏱️  Timeout reached after %v (%d messages read)\n",
			time.Since(startTime), readCount)
	}
}
