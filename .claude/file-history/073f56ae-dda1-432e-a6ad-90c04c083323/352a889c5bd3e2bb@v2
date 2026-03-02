package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	var (
		addr      = flag.String("addr", "localhost:9090", "Monitoring server address")
		checkType = flag.String("check", "health", "Check type: health, ready, metrics, all")
		jsonFlag  = flag.Bool("json", false, "Output in JSON format")
		timeout   = flag.Duration("timeout", 5*time.Second, "Request timeout")
	)
	flag.Parse()

	client := &http.Client{Timeout: *timeout}

	switch *checkType {
	case "health":
		exitCode := checkHealth(client, *addr, *jsonFlag)
		os.Exit(exitCode)
	case "ready":
		exitCode := checkReady(client, *addr, *jsonFlag)
		os.Exit(exitCode)
	case "metrics":
		exitCode := checkMetrics(client, *addr, *jsonFlag)
		os.Exit(exitCode)
	case "all":
		health := checkHealth(client, *addr, *jsonFlag)
		ready := checkReady(client, *addr, *jsonFlag)
		metrics := checkMetrics(client, *addr, *jsonFlag)

		if *jsonFlag {
			result := map[string]interface{}{
				"health":  health == 0,
				"ready":   ready == 0,
				"metrics": metrics == 0,
				"time":    time.Now().UTC().Format(time.RFC3339),
			}
			b, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(b))
		}

		// Return error if any check failed
		if health != 0 || ready != 0 || metrics != 0 {
			os.Exit(1)
		}
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unknown check type: %s\n", *checkType)
		fmt.Fprintf(os.Stderr, "valid types: health, ready, metrics, all\n")
		os.Exit(2)
	}
}

func checkHealth(client *http.Client, addr string, jsonOut bool) int {
	url := fmt.Sprintf("http://%s/health", addr)
	resp, err := client.Get(url)
	if err != nil {
		if !jsonOut {
			fmt.Fprintf(os.Stderr, "❌ health check failed: %v\n", err)
		}
		return 2 // Connection error
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if !jsonOut {
			fmt.Fprintf(os.Stderr, "❌ health check failed: status %d\n", resp.StatusCode)
		}
		return 1 // Health check failed
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		if !jsonOut {
			fmt.Fprintf(os.Stderr, "❌ failed to parse response: %v\n", err)
		}
		return 1
	}

	if !jsonOut {
		fmt.Printf("✅ health check passed: %v\n", result["status"])
	}
	return 0
}

func checkReady(client *http.Client, addr string, jsonOut bool) int {
	url := fmt.Sprintf("http://%s/ready", addr)
	resp, err := client.Get(url)
	if err != nil {
		if !jsonOut {
			fmt.Fprintf(os.Stderr, "❌ readiness check failed: %v\n", err)
		}
		return 2 // Connection error
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		if !jsonOut {
			fmt.Fprintf(os.Stderr, "❌ failed to parse response: %v\n", err)
		}
		return 1
	}

	if resp.StatusCode == http.StatusServiceUnavailable {
		if !jsonOut {
			fmt.Fprintf(os.Stderr, "❌ readiness check failed: not ready\n")
		}
		return 1
	}

	if resp.StatusCode != http.StatusOK {
		if !jsonOut {
			fmt.Fprintf(os.Stderr, "❌ readiness check failed: status %d\n", resp.StatusCode)
		}
		return 1
	}

	if !jsonOut {
		fmt.Printf("✅ readiness check passed: %v\n", result["status"])
	}
	return 0
}

func checkMetrics(client *http.Client, addr string, jsonOut bool) int {
	url := fmt.Sprintf("http://%s/metrics", addr)
	resp, err := client.Get(url)
	if err != nil {
		if !jsonOut {
			fmt.Fprintf(os.Stderr, "❌ metrics check failed: %v\n", err)
		}
		return 2 // Connection error
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if !jsonOut {
			fmt.Fprintf(os.Stderr, "❌ metrics check failed: status %d\n", resp.StatusCode)
		}
		return 1
	}

	// Check for prometheus format
	if resp.Header.Get("Content-Type") != "text/plain; version=0.0.4; charset=utf-8" {
		if !jsonOut {
			fmt.Fprintf(os.Stderr, "⚠️  unexpected content-type: %s\n", resp.Header.Get("Content-Type"))
		}
		// Don't fail, just warn
	}

	if !jsonOut {
		fmt.Printf("✅ metrics check passed\n")
	}
	return 0
}
