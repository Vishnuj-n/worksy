package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// StatsData holds daily session counts and streak, persisted to stats.json.
type StatsData struct {
	// Today's date as "YYYY-MM-DD" — used to detect day rollover
	Date          string `json:"date"`
	SessionsToday int    `json:"sessionsToday"`
	// LastActiveDate is the most recent day a session was completed ("YYYY-MM-DD")
	LastActiveDate string `json:"lastActiveDate"`
	Streak         int    `json:"streak"`
}

// StatsService tracks completed sessions per day and running streak.
type StatsService struct {
	mu       sync.Mutex
	data     StatsData
	filePath string
}

func NewStatsService() *StatsService {
	appData, _ := os.UserCacheDir()
	dir := filepath.Join(appData, "FocusPlay")
	_ = os.MkdirAll(dir, 0755)
	ss := &StatsService{
		filePath: filepath.Join(dir, "stats.json"),
	}
	ss.load()
	ss.rolloverIfNeeded()
	return ss
}

// GetStats returns the current stats snapshot.
func (ss *StatsService) GetStats() StatsData {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.rolloverIfNeededUnlocked()
	return ss.data
}

// RecordSessionComplete increments today's count and updates the streak.
func (ss *StatsService) RecordSessionComplete() StatsData {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	today := todayStr()
	ss.rolloverIfNeededUnlocked()

	ss.data.SessionsToday++

	// Streak logic
	switch ss.data.LastActiveDate {
	case "":
		// First ever session
		ss.data.Streak = 1
	case today:
		// Already recorded a session today — streak stays the same
	default:
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		if ss.data.LastActiveDate == yesterday {
			ss.data.Streak++ // consecutive day
		} else {
			ss.data.Streak = 1 // gap — reset
		}
	}
	ss.data.LastActiveDate = today
	ss.save()
	return ss.data
}

// ── internal ──────────────────────────────────────────────────────────────────

func todayStr() string {
	return time.Now().Format("2006-01-02")
}

func (ss *StatsService) rolloverIfNeeded() {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.rolloverIfNeededUnlocked()
}

// rolloverIfNeededUnlocked resets sessionsToday when the calendar day changes.
// Must be called with ss.mu held.
func (ss *StatsService) rolloverIfNeededUnlocked() {
	today := todayStr()
	if ss.data.Date != today {
		ss.data.Date = today
		ss.data.SessionsToday = 0
		ss.save()
	}
}

func (ss *StatsService) load() {
	data, err := os.ReadFile(ss.filePath)
	if err != nil {
		ss.data = StatsData{Date: todayStr()}
		return
	}
	_ = json.Unmarshal(data, &ss.data)
}

func (ss *StatsService) save() {
	data, err := json.MarshalIndent(ss.data, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(ss.filePath, data, 0644)
}
