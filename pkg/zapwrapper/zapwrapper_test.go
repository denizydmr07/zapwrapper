package zapwrapper_test

import (
	"os"
	"testing"
	"time"

	"github.com/denizydmr07/zapwrapper/pkg/zapwrapper"
	"go.uber.org/zap/zapcore"
)

const testLogDir = "./test_logs"

func setup() {
	// Create test log directory
	os.Mkdir(testLogDir, 0755)
}

func teardown() {
	// Remove test log directory and its contents
	os.RemoveAll(testLogDir)
}

// TestLogger tests the NewLogger function by creating multiple log files.
// Primary focus is on ensuring that the log rotation mechanism works
// correctly when new log files are created with each run
func TestLogger(t *testing.T) {
	setup()
	defer teardown()

	// Create a logger with a small maxBackup to test rotation
	maxBackup := 3
	logger := zapwrapper.NewLogger(testLogDir, maxBackup, zapcore.DebugLevel)

	// Close the logger to ensure the file is written
	defer logger.Sync()

	// Generate multiple runs to create multiple log files
	for i := 0; i < 5; i++ {
		// Create a new logger instance to simulate a new run
		logger = zapwrapper.NewLogger(testLogDir, maxBackup, zapcore.DebugLevel)
		logger.Info("This is a test log message", zapcore.Field{
			Key:     "iteration",
			Type:    zapcore.Int64Type,
			Integer: int64(i),
		})
		time.Sleep(1 * time.Second) // Ensure logs have different timestamps

		// Close the logger to ensure the file is written
		logger.Sync()
	}

	// Check the log directory
	files, err := os.ReadDir(testLogDir)
	if err != nil {
		t.Fatalf("Failed to read log directory: %v", err)
	}

	// Verify the number of log files
	if len(files) > maxBackup+1 {
		t.Fatalf("Expected a maximum of %d log files, found %d", maxBackup+1, len(files))
	}

	// Check if the files are named correctly
	for _, file := range files {
		if file.IsDir() {
			t.Fatalf("Expected file, found directory: %s", file.Name())
		}
		if !file.Type().IsRegular() {
			t.Fatalf("Expected regular file, found: %s", file.Name())
		}
		if len(file.Name()) == 0 {
			t.Fatalf("Found file with empty name")
		}
	}
}
