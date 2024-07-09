## zapwrapper

A personal wrapper around `zap.Logger` that writes logs to both the console and a file.

### usage

```go
package main

import (
	"github.com/yourusername/yourproject/zapwrapper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger := zapwrapper.NewLogger(
		zapwrapper.DefaultFilepath      // Log file path
        zapwrapper.DefaultMaxBackups    // Max number of log files to retain
		zapwrapper.DefaultLogLevel,     // Log level
	)

	defer logger.Sync() // Flush any buffered log entries

	logger.Info("This is an info message")
	logger.Debug("This is a debug message")
	logger.Error("This is an error message")
}
```
### note

This package is for personal use.


