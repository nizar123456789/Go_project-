package logger

import (
    "fmt"
    "log"
    "os"
    "time"
)

var (
    // Logger is the global logger instance
    Logger *log.Logger
)

// Initialize sets up the logging environment
func Initialize() error {
    // Create log directory if it doesn't exist
    logDir := "logs"
    if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
        return fmt.Errorf("error creating log directory: %w", err)
    }

    // Generate log file name with timestamp
    logFileName := fmt.Sprintf("%s.log", time.Now().Format("01_02_2006_15_04_05"))
    logFilePath := fmt.Sprintf("%s/%s", logDir, logFileName)

    // Open log file for writing
    file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        return fmt.Errorf("error opening log file: %w", err)
    }

    // Create a new logger
    Logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)

    return nil
}

// LogInfo logs an info message
func LogInfo(message string) {
    if Logger != nil {
        Logger.Println(message)
    }
}

// LogError logs an error message
func LogError(message string) {
    if Logger != nil {
        Logger.Println("ERROR: " + message)
    }
}