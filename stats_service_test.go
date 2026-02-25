package main

import (
	"path/filepath"
	"testing"
	"time"
)

func TestStatsServiceInit(t *testing.T) {
	tmpDir := t.TempDir()
	ss := &StatsService{filePath: filepath.Join(tmpDir, "stats.json")}
	ss.load()

	stats := ss.GetStats()
	if stats.SessionsToday != 0 {
		t.Errorf("Expected 0 sessions today on init, got %d", stats.SessionsToday)
	}
	if stats.Streak != 0 {
		t.Errorf("Expected 0 streak on init, got %d", stats.Streak)
	}
	if stats.Date == "" {
		t.Error("Date should be set on init")
	}
}

func TestStatsServiceRecordSessionComplete(t *testing.T) {
	tmpDir := t.TempDir()
	ss := &StatsService{filePath: filepath.Join(tmpDir, "stats.json")}
	ss.load()

	// First session completes
	stats1 := ss.RecordSessionComplete()
	if stats1.SessionsToday != 1 {
		t.Errorf("Expected 1 session today, got %d", stats1.SessionsToday)
	}
	if stats1.Streak != 1 {
		t.Errorf("Expected streak 1 for first session, got %d", stats1.Streak)
	}

	// Second session same day
	stats2 := ss.RecordSessionComplete()
	if stats2.SessionsToday != 2 {
		t.Errorf("Expected 2 sessions today, got %d", stats2.SessionsToday)
	}
	if stats2.Streak != 1 {
		t.Errorf("Expected streak still 1, got %d", stats2.Streak)
	}
}

func TestStatsServiceStreakIncrement(t *testing.T) {
	tmpDir := t.TempDir()
	ss := &StatsService{filePath: filepath.Join(tmpDir, "stats.json")}
	ss.load()

	// First day: record session
	ss.RecordSessionComplete()
	stats1 := ss.GetStats()
	if stats1.Streak != 1 {
		t.Errorf("Day 1: expected streak 1, got %d", stats1.Streak)
	}

	// Manually set LastActiveDate to yesterday and reset SessionsToday for next day
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	ss.data.LastActiveDate = yesterday
	ss.data.SessionsToday = 0
	ss.data.Date = time.Now().Format("2006-01-02")

	// Second day: record should increment streak
	ss.RecordSessionComplete()
	stats2 := ss.GetStats()
	if stats2.Streak != 2 {
		t.Errorf("Day 2 (consecutive): expected streak 2, got %d", stats2.Streak)
	}
}

func TestStatsServiceStreakReset(t *testing.T) {
	tmpDir := t.TempDir()
	ss := &StatsService{filePath: filepath.Join(tmpDir, "stats.json")}
	ss.load()

	// Set a past active date (2 days ago)
	twoDaysAgo := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	ss.data.LastActiveDate = twoDaysAgo
	ss.data.Streak = 5

	// Record session today — should reset streak
	ss.RecordSessionComplete()
	stats := ss.GetStats()
	if stats.Streak != 1 {
		t.Errorf("Expected streak reset to 1 after gap, got %d", stats.Streak)
	}
}

func TestStatsServiceRollover(t *testing.T) {
	tmpDir := t.TempDir()
	ss := &StatsService{filePath: filepath.Join(tmpDir, "stats.json")}
	ss.load()

	// Manually set date to yesterday
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	ss.data.Date = yesterday
	ss.data.SessionsToday = 5
	ss.data.Streak = 3

	// GetStats on new day should trigger rollover
	stats := ss.GetStats()
	if stats.SessionsToday != 0 {
		t.Errorf("After rollover, SessionsToday should be 0, got %d", stats.SessionsToday)
	}
	if stats.Streak != 3 {
		t.Errorf("Rollover should not change streak, got %d", stats.Streak)
	}
	today := time.Now().Format("2006-01-02")
	if stats.Date != today {
		t.Errorf("After rollover, date should be today (%s), got %s", today, stats.Date)
	}
}

func TestStatsServicePersistence(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "stats.json")

	// Create and record
	ss1 := &StatsService{filePath: filePath}
	ss1.load()
	ss1.RecordSessionComplete()
	ss1.RecordSessionComplete()

	// Load in new instance
	ss2 := &StatsService{filePath: filePath}
	ss2.load()

	stats := ss2.GetStats()
	if stats.SessionsToday != 2 {
		t.Errorf("After persistence, expected 2 sessions today, got %d", stats.SessionsToday)
	}
}

func TestStatsServiceNoFile(t *testing.T) {
	tmpDir := t.TempDir()
	ss := &StatsService{filePath: filepath.Join(tmpDir, "nonexistent.json")}
	ss.load()

	// Should have default/empty stats
	stats := ss.GetStats()
	if stats.SessionsToday != 0 {
		t.Errorf("Expected 0 sessions with no file, got %d", stats.SessionsToday)
	}
}

func TestStatsServiceMultipleSessions(t *testing.T) {
	tmpDir := t.TempDir()
	ss := &StatsService{filePath: filepath.Join(tmpDir, "stats.json")}
	ss.load()

	// Record 10 sessions in one day
	for i := 0; i < 10; i++ {
		ss.RecordSessionComplete()
	}

	stats := ss.GetStats()
	if stats.SessionsToday != 10 {
		t.Errorf("Expected 10 sessions, got %d", stats.SessionsToday)
	}
	if stats.Streak != 1 {
		t.Errorf("Expected streak 1 (same day), got %d", stats.Streak)
	}
}
