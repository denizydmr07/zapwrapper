// Description: A wrapper around the zap.Logger that writes to both the console and a file.
package zapwrapper

// Import the required packages
import (
	"os"
	"sort"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Define ANSI color codes and the default log file path
const (
	// Default values for initialization
	DefaultFilepath   = "./logs" // Default log file path
	DefaultMaxSize    = 10       // Default max size of each log file in megabytes
	DefaultMaxBackups = 5        // Default max number of log files to retain
	DefaultMaxAge     = 30       // Default max number of days to retain a log file
	DefaultLogLevel   = zapcore.DebugLevel

	// color codes for the console
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorReset   = "\033[0m"
)

// NewLogger creates a logger with the specified log file path, max number of
// log files to retain, max size of each log file, and max number of days to
// retain a log file.
//
// Parameters:
//   - filepath: the path to the directory where the log files will be stored
//   - maxBackup: the maximum number of log files to retain
//   - logLevel: the log level (e.g., zapcore.InfoLevel, zapcore.DebugLevel, etc.)
//
// Returns:
//   - a new logger that writes to both the console and a file
//
// Example:
//
//	Logger := zapwrapper.NewLogger(
//	  zapwrapper.DefaultFilepath,
//	  zapwrapper.DefaultMaxBackup,
//	  zapwrapper.DefaultLogLevel,
//
// )
func NewLogger(filepath string, maxBackup int, logLevel zapcore.Level) *zap.Logger {
	// append timestamp to the log file (only the hour, minute, second includedin the timestamp)
	// formatting the timestamp as (day-month-year hour-minute-second)
	timestamp := time.Now().Format("02-01-06_15-04-05")
	filename := filepath + "/logs_" + timestamp + ".log"

	// Custom encoder configuration for the console
	consoleEncoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		// Add color to the encoded log level
		EncodeLevel: func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			var color string
			switch level {
			case zapcore.DebugLevel:
				color = colorCyan
			case zapcore.InfoLevel:
				color = colorGreen
			case zapcore.WarnLevel:
				color = colorYellow
			case zapcore.ErrorLevel:
				color = colorRed
			case zapcore.DPanicLevel:
				color = colorMagenta
			case zapcore.PanicLevel:
				color = colorMagenta
			case zapcore.FatalLevel:
				color = colorRed
			}
			enc.AppendString(color + level.CapitalString() + colorReset)
		},
		// Encode the time in the specified format
		EncodeTime:     zapcore.TimeEncoderOfLayout("02-01-06 15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Custom encoder configuration for the file (without color)
	fileEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("02-01-2006 15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create a core that writes to both the console and the file
	consoleCore := zapcore.NewCore(
		// Use the custom console encoder configuration
		zapcore.NewConsoleEncoder(consoleEncoderConfig),
		zapcore.AddSync(os.Stdout), // Write to the console
		logLevel,                   // log level (e.g., zapcore.InfoLevel, zapcore.DebugLevel, etc.
	)

	// Create a core that writes to a file
	fileCore := zapcore.NewCore(
		// Use the custom file encoder configuration
		zapcore.NewConsoleEncoder(fileEncoderConfig),
		zapcore.AddSync(&lumberjack.Logger{ //lumberjack.Logger is used to handle log rotation
			Filename: filename, // Log file name
		}),
		logLevel, // log level (e.g., zapcore.InfoLevel, zapcore.DebugLevel, etc.)

	)

	// Combine the cores
	core := zapcore.NewTee(consoleCore, fileCore)

	// check the filepath, if it exists and has more than maxBackup files,
	// delete the oldest file
	if _, err := os.Stat(filepath); err == nil {
		files, _ := os.ReadDir(filepath) // read the directory
		if len(files) > maxBackup {      // if the number of files is greater than maxBackup
			// sort the files by their names
			// (the files are named logs_15-04-05.log, logs_15-04-06.log, etc.)
			// so the oldest file is the first one
			sort.Slice(files, func(i, j int) bool {
				return files[i].Name() < files[j].Name()
			})
			// delete the oldest file
			os.Remove(filepath + "/" + files[0].Name())
		}
	}

	// Build the logger with the combined core and return it
	return zap.New(core)
}
