package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// LogLevel represents log severity
type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

// LogEntry is a parsed log entry for the viewer
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	IP        string `json:"ip"`
	User      string `json:"user"`
	Action    string `json:"action"`
	Details   string `json:"details"`
}

var (
	logFile *os.File
	logMu   sync.Mutex
)

func initLog() {
	var err error
	logFile, err = os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal("Cannot open log file:", err)
	}
}

// logEvent writes a structured log entry
func logEvent(level LogLevel, ip, user, action, details string) {
	logMu.Lock()
	defer logMu.Unlock()
	ts := time.Now().Format("2006-01-02 15:04:05")
	if user == "" {
		user = "-"
	}
	if ip == "" {
		ip = "-"
	}
	line := fmt.Sprintf("[%s] [%s] IP=%s USER=%s ACTION=%s %s\n", ts, level, ip, user, action, details)
	if logFile != nil {
		logFile.WriteString(line)
	}
	log.Print(strings.TrimSpace(line))
}

// Convenience helpers
func logInfo(ip, user, action, details string)  { logEvent(LevelInfo, ip, user, action, details) }
func logWarn(ip, user, action, details string)  { logEvent(LevelWarn, ip, user, action, details) }
func logError(ip, user, action, details string) { logEvent(LevelError, ip, user, action, details) }
func logDebug(ip, user, action, details string) { logEvent(LevelDebug, ip, user, action, details) }

// writeLog is kept for backwards-compat — defaults to INFO level
func writeLog(ip, user, action, details string) { logInfo(ip, user, action, details) }

// getIP extracts the client IP from a request, honoring proxies only if from trusted source
func getIP(r *http.Request) string {
	remoteIP := strings.Split(r.RemoteAddr, ":")[0]
	isTrustedProxy := remoteIP == "127.0.0.1" || remoteIP == "::1" || remoteIP == "localhost"
	if isTrustedProxy {
		if f := r.Header.Get("X-Forwarded-For"); f != "" {
			return strings.Split(f, ",")[0]
		}
		if f := r.Header.Get("X-Real-IP"); f != "" {
			return f
		}
	}
	return remoteIP
}

// Regex matches both new format with [LEVEL] and old format without
var (
	logRegexNew = regexp.MustCompile(`^\[([\d\- :]+)\] \[(\w+)\] IP=(\S+) USER=(\S+) ACTION=(\S+)\s*(.*)$`)
	logRegexOld = regexp.MustCompile(`^\[([\d\- :]+)\] IP=(\S+) USER=(\S+) ACTION=(\S+)\s*(.*)$`)
)

// parseLogFile reads heinen.log and returns parsed entries (newest first)
func parseLogFile() ([]LogEntry, error) {
	f, err := os.Open(LogFile)
	if err != nil {
		return []LogEntry{}, nil
	}
	defer f.Close()
	entries := make([]LogEntry, 0)
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		if m := logRegexNew.FindStringSubmatch(line); m != nil {
			entries = append(entries, LogEntry{Timestamp: m[1], Level: m[2], IP: m[3], User: m[4], Action: m[5], Details: m[6]})
		} else if m := logRegexOld.FindStringSubmatch(line); m != nil {
			entries = append(entries, LogEntry{Timestamp: m[1], Level: "INFO", IP: m[2], User: m[3], Action: m[4], Details: m[5]})
		}
	}
	// Reverse so newest is first
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}
	return entries, nil
}

// filterLogs applies search and level filters
func filterLogs(entries []LogEntry, search, level string) []LogEntry {
	if search == "" && level == "" {
		return entries
	}
	search = strings.ToLower(search)
	out := make([]LogEntry, 0, len(entries))
	for _, e := range entries {
		if level != "" && level != "ALL" && e.Level != level {
			continue
		}
		if search != "" {
			hay := strings.ToLower(e.IP + " " + e.User + " " + e.Action + " " + e.Details + " " + e.Timestamp)
			if !strings.Contains(hay, search) {
				continue
			}
		}
		out = append(out, e)
	}
	return out
}
