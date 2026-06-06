package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type LogEntry struct {
	Environment string `json:"environment"`
	CommitHash  string `json:"commit_hash"`
	Message     string `json:"message"`
}

var (
	url         string
	apiKey      string
	environment string
	commitHash  string
	batchSize   int
	flushInterval time.Duration
)

func init() {
	flag.StringVar(&url, "url", "http://localhost:8080/api/logs/ingest", "TelemetryMend ingestion URL")
	flag.StringVar(&apiKey, "key", os.Getenv("TM_API_KEY"), "API Key (defaults to TM_API_KEY env)")
	flag.StringVar(&environment, "env", "production", "Environment name")
	flag.StringVar(&commitHash, "commit", "", "Commit hash")
	flag.IntVar(&batchSize, "batch", 10, "Batch size for sending logs")
	flag.DurationVar(&flushInterval, "interval", 5*time.Second, "Maximum time to wait before flushing logs")
}

func main() {
	flag.Parse()

	if apiKey == "" {
		log.Fatal("API Key is required. Use -key or set TM_API_KEY environment variable.")
	}

	log.Printf("Starting tm-ingest. Target: %s, Env: %s", url, environment)

	// Create a channel for logs
	logChan := make(chan string, 100)

	// Start a goroutine to read from stdin
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			logChan <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading stdin: %v", err)
		}
		close(logChan)
	}()

	// Start a goroutine to process and batch logs
	processLogs(logChan)
}

func processLogs(logChan <-chan string) {
	var batch []LogEntry
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		err := sendBatch(batch)
		if err != nil {
			log.Printf("Error sending batch: %v", err)
		} else {
			log.Printf("Successfully sent %d logs", len(batch))
		}
		batch = nil
	}

	for {
		select {
		case line, ok := <-logChan:
			if !ok {
				flush()
				return
			}
			batch = append(batch, LogEntry{
				Environment: environment,
				CommitHash:  commitHash,
				Message:     line,
			})
			if len(batch) >= batchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}

func sendBatch(batch []LogEntry) error {
	data, err := json.Marshal(batch)
	if err != nil {
		return fmt.Errorf("failed to marshal batch: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
