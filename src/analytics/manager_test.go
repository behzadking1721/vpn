package analytics

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/history"
	"c:/Users/behza/OneDrive/Documents/vpn/src/managers"
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"testing"
	"time"
)

func TestAnalyticsManager(t *testing.T) {
	// Create mock managers
	configManager := managers.NewConfigManager("./test_config.json")
	serverManager := managers.NewServerManager()
	historyManager := history.NewHistoryManager("./test_history")
	
	// Create analytics manager
	analyticsManager := NewAnalyticsManager(serverManager, historyManager)
	
	// Test ping statistics calculation
	t.Run("CalculatePingStats", func(t *testing.T) {
		// Add servers with different ping values
		server1 := core.Server{
			ID:       utils.GenerateID(),
			Name:     "Server 1",
			Host:     "server1.example.com",
			Port:     443,
			Protocol: core.ProtocolVMess,
			Ping:     50,
			Enabled:  true,
		}
		
		server2 := core.Server{
			ID:       utils.GenerateID(),
			Name:     "Server 2",
			Host:     "server2.example.com",
			Port:     443,
			Protocol: core.ProtocolShadowsocks,
			Ping:     100,
			Enabled:  true,
		}
		
		server3 := core.Server{
			ID:       utils.GenerateID(),
			Name:     "Server 3",
			Host:     "server3.example.com",
			Port:     443,
			Protocol: core.ProtocolTrojan,
			Ping:     150,
			Enabled:  true,
		}
		
		// Add servers to manager
		serverManager.AddServer(server1)
		serverManager.AddServer(server2)
		serverManager.AddServer(server3)
		
		// Calculate ping stats
		stats, err := analyticsManager.CalculatePingStats("7d")
		if err != nil {
			t.Fatalf("Failed to calculate ping stats: %v", err)
		}
		
		// Check results
		if stats.Samples != 3 {
			t.Errorf("Expected 3 samples, got %d", stats.Samples)
		}
		
		if stats.Average != 100.0 {
			t.Errorf("Expected average of 100.0, got %f", stats.Average)
		}
		
		if stats.Min != 50 {
			t.Errorf("Expected min of 50, got %d", stats.Min)
		}
		
		if stats.Max != 150 {
			t.Errorf("Expected max of 150, got %d", stats.Max)
		}
		
		if stats.P95 != 150 {
			t.Errorf("Expected p95 of 150, got %d", stats.P95)
		}
	})
	
	// Test data usage statistics calculation
	t.Run("CalculateDataUsageStats", func(t *testing.T) {
		// Add data usage records
		record1 := history.DataUsageRecord{
			ID:           utils.GenerateID(),
			Timestamp:    time.Now().Add(-24 * time.Hour),
			ServerID:     "server1",
			ServerName:   "Server 1",
			DataSent:     1024 * 1024 * 100, // 100 MB
			DataReceived: 1024 * 1024 * 200, // 200 MB
			TotalSent:    1024 * 1024 * 100,
			TotalReceived: 1024 * 1024 * 200,
		}
		
		record2 := history.DataUsageRecord{
			ID:           utils.GenerateID(),
			Timestamp:    time.Now().Add(-12 * time.Hour),
			ServerID:     "server2",
			ServerName:   "Server 2",
			DataSent:     1024 * 1024 * 50, // 50 MB
			DataReceived: 1024 * 1024 * 150, // 150 MB
			TotalSent:    1024 * 1024 * 50,
			TotalReceived: 1024 * 1024 * 150,
		}
		
		// Add records to history manager
		historyManager.AddDataUsageRecord(record1)
		historyManager.AddDataUsageRecord(record2)
		
		// Calculate data usage stats
		stats, err := analyticsManager.CalculateDataUsageStats("7d")
		if err != nil {
			t.Fatalf("Failed to calculate data usage stats: %v", err)
		}
		
		// Check results
		expectedTotalSent := int64(1024 * 1024 * 150)     // 150 MB
		expectedTotalReceived := int64(1024 * 1024 * 350) // 350 MB
		
		if stats.TotalSent != expectedTotalSent {
			t.Errorf("Expected total sent of %d, got %d", expectedTotalSent, stats.TotalSent)
		}
		
		if stats.TotalReceived != expectedTotalReceived {
			t.Errorf("Expected total received of %d, got %d", expectedTotalReceived, stats.TotalReceived)
		}
		
		// Check average per day (2 days of data)
		expectedAvgPerDay := (expectedTotalSent + expectedTotalReceived) / 2
		if stats.AveragePerDay != expectedAvgPerDay {
			t.Errorf("Expected average per day of %d, got %d", expectedAvgPerDay, stats.AveragePerDay)
		}
	})
	
	// Test time patterns calculation
	t.Run("CalculateTimePatterns", func(t *testing.T) {
		// Add connection records at different hours
		record1 := history.ConnectionRecord{
			ID:             utils.GenerateID(),
			ServerID:       "server1",
			ServerName:     "Server 1",
			StartTime:      time.Now().Add(-3 * time.Hour),
			EndTime:        time.Now().Add(-2 * time.Hour),
			Duration:       3600,
			DataSent:       1024 * 1024 * 50,
			DataReceived:   1024 * 1024 * 100,
			Protocol:       "VMess",
			Status:         "connected",
			DisconnectReason: "",
		}
		
		record2 := history.ConnectionRecord{
			ID:             utils.GenerateID(),
			ServerID:       "server2",
			ServerName:     "Server 2",
			StartTime:      time.Now().Add(-2 * time.Hour),
			EndTime:        time.Now().Add(-1 * time.Hour),
			Duration:       3600,
			DataSent:       1024 * 1024 * 30,
			DataReceived:   1024 * 1024 * 70,
			Protocol:       "Shadowsocks",
			Status:         "disconnected",
			DisconnectReason: "timeout",
		}
		
		// Add records to history manager
		historyManager.AddConnectionRecord(record1)
		historyManager.AddConnectionRecord(record2)
		
		// Calculate time patterns
		patterns, err := analyticsManager.CalculateTimePatterns("7d")
		if err != nil {
			t.Fatalf("Failed to calculate time patterns: %v", err)
		}
		
		// Check results
		if len(patterns) == 0 {
			t.Error("Expected at least one time pattern")
		}
		
		// Find the pattern for the current hour
		currentHour := time.Now().Hour()
		var foundPattern *TimePattern
		for i := range patterns {
			if patterns[i].Hour == currentHour {
				foundPattern = &patterns[i]
				break
			}
		}
		
		if foundPattern == nil {
			t.Errorf("Expected pattern for hour %d", currentHour)
		} else {
			// Check usage count
			if foundPattern.UsageCount < 1 {
				t.Errorf("Expected at least 1 usage count for hour %d, got %d", currentHour, foundPattern.UsageCount)
			}
			
			// Check disconnects
			if foundPattern.Disconnects < 1 {
				t.Errorf("Expected at least 1 disconnect for hour %d, got %d", currentHour, foundPattern.Disconnects)
			}
		}
	})
	
	// Test report generation
	t.Run("GenerateReport", func(t *testing.T) {
		// Generate weekly report
		report, err := analyticsManager.GenerateReport(ReportPeriodWeekly)
		if err != nil {
			t.Fatalf("Failed to generate weekly report: %v", err)
		}
		
		// Check report properties
		if report.Period != ReportPeriodWeekly {
			t.Errorf("Expected period %s, got %s", ReportPeriodWeekly, report.Period)
		}
		
		if report.StartDate.After(report.EndDate) {
			t.Error("Expected start date to be before end date")
		}
		
		// Generate monthly report
		report, err = analyticsManager.GenerateReport(ReportPeriodMonthly)
		if err != nil {
			t.Fatalf("Failed to generate monthly report: %v", err)
		}
		
		// Check report properties
		if report.Period != ReportPeriodMonthly {
			t.Errorf("Expected period %s, got %s", ReportPeriodMonthly, report.Period)
		}
		
		if report.StartDate.After(report.EndDate) {
			t.Error("Expected start date to be before end date")
		}
		
		// Try to generate report with invalid period
		_, err = analyticsManager.GenerateReport("invalid")
		if err == nil {
			t.Error("Expected error for invalid report period")
		}
	})
	
	// Test daily data usage
	t.Run("GetDailyDataUsage", func(t *testing.T) {
		// Get daily data usage
		dailyUsage, err := analyticsManager.GetDailyDataUsage("7d")
		if err != nil {
			t.Fatalf("Failed to get daily data usage: %v", err)
		}
		
		// Check results
		if len(dailyUsage) == 0 {
			t.Error("Expected at least one daily usage record")
		}
	})
	
	// Test weekly data usage
	t.Run("GetWeeklyDataUsage", func(t *testing.T) {
		// Get weekly data usage
		weeklyUsage, err := analyticsManager.GetWeeklyDataUsage("30d")
		if err != nil {
			t.Fatalf("Failed to get weekly data usage: %v", err)
		}
		
		// Check results
		if len(weeklyUsage) == 0 {
			t.Error("Expected at least one weekly usage record")
		}
	})
	
	// Test monthly data usage
	t.Run("GetMonthlyDataUsage", func(t *testing.T) {
		// Get monthly data usage
		monthlyUsage, err := analyticsManager.GetMonthlyDataUsage("365d")
		if err != nil {
			t.Fatalf("Failed to get monthly data usage: %v", err)
		}
		
		// Check results
		if len(monthlyUsage) == 0 {
			t.Error("Expected at least one monthly usage record")
		}
	})
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
		hasError bool
	}{
		{"7d", 7 * 24 * time.Hour, false},
		{"30d", 30 * 24 * time.Hour, false},
		{"1h", 1 * time.Hour, false},
		{"30m", 30 * time.Minute, false},
		{"45s", 45 * time.Second, false},
		{"invalid", 0, true},
		{"", 0, true},
		{"d", 0, true},
	}
	
	for _, test := range tests {
		duration, err := parseDuration(test.input)
		
		if test.hasError {
			if err == nil {
				t.Errorf("Expected error for input '%s', but got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			} else if duration != test.expected {
				t.Errorf("Expected duration %v for input '%s', got %v", test.expected, test.input, duration)
			}
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		input    time.Duration
		expected string
	}{
		{7 * 24 * time.Hour, "7d"},
		{30 * 24 * time.Hour, "30d"},
		{1 * time.Hour, "1h"},
		{30 * time.Minute, "0h"},
		{45 * time.Second, "0h"},
	}
	
	for _, test := range tests {
		result := formatDuration(test.input)
		if result != test.expected {
			t.Errorf("Expected '%s' for duration %v, got '%s'", test.expected, test.input, result)
		}
	}
}

func TestBeginningOfWeek(t *testing.T) {
	// Test with a known date
	testDate := time.Date(2023, 10, 25, 15, 30, 0, 0, time.UTC) // Wednesday
	expected := time.Date(2023, 10, 23, 0, 0, 0, 0, time.UTC)    // Monday
	
	result := beginningOfWeek(testDate)
	if !result.Equal(expected) {
		t.Errorf("Expected beginning of week to be %v, got %v", expected, result)
	}
}

func TestBeginningOfMonth(t *testing.T) {
	// Test with a known date
	testDate := time.Date(2023, 10, 25, 15, 30, 0, 0, time.UTC) // October 25th
	expected := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)     // October 1st
	
	result := beginningOfMonth(testDate)
	if !result.Equal(expected) {
		t.Errorf("Expected beginning of month to be %v, got %v", expected, result)
	}
}