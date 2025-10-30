package analytics

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/history"
	"c:/Users/behza/OneDrive/Documents/vpn/src/managers"
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"sort"
	"time"
	"math"
	"fmt"
)

// AnalyticsManager handles analytics calculations
type AnalyticsManager struct {
	serverManager  *managers.ServerManager
	historyManager *history.HistoryManager
}

// NewAnalyticsManager creates a new analytics manager
func NewAnalyticsManager(
	serverMgr *managers.ServerManager,
	historyMgr *history.HistoryManager) *AnalyticsManager {
	
	return &AnalyticsManager{
		serverManager:  serverMgr,
		historyManager: historyMgr,
	}
}

// CalculatePingStats calculates ping statistics for a given time window
func (am *AnalyticsManager) CalculatePingStats(window string) (*PingStats, error) {
	// Parse window (e.g., "7d", "30d")
	duration, err := parseDuration(window)
	if err != nil {
		return nil, err
	}
	
	// Calculate time range
	endTime := time.Now()
	startTime := endTime.Add(-duration)
	
	// Get server data within the time range
	servers := am.serverManager.GetAllServers()
	
	// Collect all ping values
	var pingValues []int
	
	for _, server := range servers {
		// We're using the current ping value as a simplification
		// In a real implementation, you would track ping history over time
		if server.Ping > 0 {
			pingValues = append(pingValues, server.Ping)
		}
	}
	
	if len(pingValues) == 0 {
		return &PingStats{}, nil
	}
	
	// Sort ping values for percentile calculation
	sort.Ints(pingValues)
	
	// Calculate statistics
	var sum int
	for _, ping := range pingValues {
		sum += ping
	}
	
	avg := float64(sum) / float64(len(pingValues))
	max := pingValues[len(pingValues)-1]
	min := pingValues[0]
	
	// Calculate 95th percentile
	p95Index := int(math.Ceil(0.95*float64(len(pingValues)))) - 1
	if p95Index >= len(pingValues) {
		p95Index = len(pingValues) - 1
	}
	p95 := pingValues[p95Index]
	
	return &PingStats{
		Average: avg,
		Max:     max,
		Min:     min,
		P95:     p95,
		Samples: len(pingValues),
	}, nil
}

// CalculateDataUsageStats calculates data usage statistics for a given time window
func (am *AnalyticsManager) CalculateDataUsageStats(window string) (*DataUsageStats, error) {
	// Parse window (e.g., "7d", "30d")
	duration, err := parseDuration(window)
	if err != nil {
		return nil, err
	}
	
	// Calculate time range
	endTime := time.Now()
	startTime := endTime.Add(-duration)
	
	// Get data usage records within the time range
	records, err := am.historyManager.GetDataUsageRecords("", 0, 0)
	if err != nil {
		return nil, err
	}
	
	// Filter records within time range
	var filteredRecords []history.DataUsageRecord
	for _, record := range records {
		if record.Timestamp.After(startTime) && record.Timestamp.Before(endTime) {
			filteredRecords = append(filteredRecords, record)
		}
	}
	
	if len(filteredRecords) == 0 {
		return &DataUsageStats{}, nil
	}
	
	// Group by day
	dailyUsage := make(map[string]*DailyDataUsage)
	for _, record := range filteredRecords {
		dateKey := record.Timestamp.Format("2006-01-02")
		if _, exists := dailyUsage[dateKey]; !exists {
			dailyUsage[dateKey] = &DailyDataUsage{
				Date: record.Timestamp,
			}
		}
		dailyUsage[dateKey].DataSent += record.DataSent
		dailyUsage[dateKey].DataReceived += record.DataReceived
	}
	
	// Calculate totals
	var totalSent, totalReceived int64
	for _, usage := range dailyUsage {
		totalSent += usage.DataSent
		totalReceived += usage.DataReceived
	}
	
	// Calculate daily stats
	var maxPerDay int64
	var dailyTotals []int64
	
	for _, usage := range dailyUsage {
		dailyTotal := usage.DataSent + usage.DataReceived
		dailyTotals = append(dailyTotals, dailyTotal)
		
		if dailyTotal > maxPerDay {
			maxPerDay = dailyTotal
		}
	}
	
	// Calculate average per day
	var avgPerDay int64
	if len(dailyTotals) > 0 {
		var sum int64
		for _, total := range dailyTotals {
			sum += total
		}
		avgPerDay = sum / int64(len(dailyTotals))
	}
	
	return &DataUsageStats{
		TotalSent:     totalSent,
		TotalReceived: totalReceived,
		AveragePerDay: avgPerDay,
		MaxPerDay:     maxPerDay,
	}, nil
}

// CalculateTimePatterns calculates time-based patterns
func (am *AnalyticsManager) CalculateTimePatterns(window string) ([]TimePattern, error) {
	// Parse window (e.g., "7d", "30d")
	duration, err := parseDuration(window)
	if err != nil {
		return nil, err
	}
	
	// Calculate time range
	endTime := time.Now()
	startTime := endTime.Add(-duration)
	
	// Get connection records within the time range
	records, err := am.historyManager.GetConnectionRecords(0, 0)
	if err != nil {
		return nil, err
	}
	
	// Filter records within time range
	var filteredRecords []history.ConnectionRecord
	for _, record := range records {
		if record.StartTime.After(startTime) && record.StartTime.Before(endTime) {
			filteredRecords = append(filteredRecords, record)
		}
	}
	
	// Group by hour
	hourlyPatterns := make(map[int]*TimePattern)
	
	for _, record := range filteredRecords {
		hour := record.StartTime.Hour()
		if _, exists := hourlyPatterns[hour]; !exists {
			hourlyPatterns[hour] = &TimePattern{
				Hour: hour,
			}
		}
		hourlyPatterns[hour].UsageCount++
		
		// Count disconnects (simplified - in a real implementation, you would track disconnect reasons)
		if record.Status == "disconnected" || record.Status == "error" {
			hourlyPatterns[hour].Disconnects++
		}
	}
	
	// Convert to slice
	var patterns []TimePattern
	for _, pattern := range hourlyPatterns {
		patterns = append(patterns, *pattern)
	}
	
	// Sort by hour
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Hour < patterns[j].Hour
	})
	
	return patterns, nil
}

// GenerateReport generates a complete analytics report
func (am *AnalyticsManager) GenerateReport(period ReportPeriod) (*AnalyticsReport, error) {
	var startTime, endTime time.Time
	now := time.Now()
	
	switch period {
	case ReportPeriodWeekly:
		// Start of current week (Monday)
		startTime = beginningOfWeek(now)
		endTime = startTime.AddDate(0, 0, 7)
	case ReportPeriodMonthly:
		// Start of current month
		startTime = beginningOfMonth(now)
		endTime = startTime.AddDate(0, 1, 0)
	default:
		return nil, &InvalidPeriodError{Period: string(period)}
	}
	
	// Calculate ping stats for the period
	pingStats, err := am.CalculatePingStats(formatDuration(endTime.Sub(startTime)))
	if err != nil {
		return nil, err
	}
	
	// Calculate data usage stats for the period
	dataUsageStats, err := am.CalculateDataUsageStats(formatDuration(endTime.Sub(startTime)))
	if err != nil {
		return nil, err
	}
	
	// Calculate time patterns for the period
	timePatterns, err := am.CalculateTimePatterns(formatDuration(endTime.Sub(startTime)))
	if err != nil {
		return nil, err
	}
	
	report := &AnalyticsReport{
		ID:          utils.GenerateID(),
		Period:      period,
		StartDate:   startTime,
		EndDate:     endTime,
		PingStats:   *pingStats,
		DataUsage:   *dataUsageStats,
		TimePattern: timePatterns,
		CreatedAt:   time.Now(),
	}
	
	return report, nil
}

// GetDailyDataUsage returns daily data usage for a given time window
func (am *AnalyticsManager) GetDailyDataUsage(window string) ([]DailyDataUsage, error) {
	// Parse window (e.g., "7d", "30d")
	duration, err := parseDuration(window)
	if err != nil {
		return nil, err
	}
	
	// Calculate time range
	endTime := time.Now()
	startTime := endTime.Add(-duration)
	
	// Get data usage records within the time range
	records, err := am.historyManager.GetDataUsageRecords("", 0, 0)
	if err != nil {
		return nil, err
	}
	
	// Filter records within time range
	var filteredRecords []history.DataUsageRecord
	for _, record := range records {
		if record.Timestamp.After(startTime) && record.Timestamp.Before(endTime) {
			filteredRecords = append(filteredRecords, record)
		}
	}
	
	// Group by day
	dailyUsage := make(map[string]*DailyDataUsage)
	for _, record := range filteredRecords {
		dateKey := record.Timestamp.Format("2006-01-02")
		if _, exists := dailyUsage[dateKey]; !exists {
			dailyUsage[dateKey] = &DailyDataUsage{
				Date: record.Timestamp,
			}
		}
		dailyUsage[dateKey].DataSent += record.DataSent
		dailyUsage[dateKey].DataReceived += record.DataReceived
	}
	
	// Convert to slice
	var result []DailyDataUsage
	for _, usage := range dailyUsage {
		result = append(result, *usage)
	}
	
	// Sort by date
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Before(result[j].Date)
	})
	
	return result, nil
}

// GetWeeklyDataUsage returns weekly data usage for a given time window
func (am *AnalyticsManager) GetWeeklyDataUsage(window string) ([]WeeklyDataUsage, error) {
	dailyUsage, err := am.GetDailyDataUsage(window)
	if err != nil {
		return nil, err
	}
	
	// Group by week
	weeklyUsage := make(map[string]*WeeklyDataUsage)
	for _, daily := range dailyUsage {
		weekStart := beginningOfWeek(daily.Date)
		weekKey := weekStart.Format("2006-01-02")
		
		if _, exists := weeklyUsage[weekKey]; !exists {
			weeklyUsage[weekKey] = &WeeklyDataUsage{
				WeekStart: weekStart,
				WeekEnd:   weekStart.AddDate(0, 0, 7),
			}
		}
		
		weeklyUsage[weekKey].DataSent += daily.DataSent
		weeklyUsage[weekKey].DataReceived += daily.DataReceived
	}
	
	// Convert to slice
	var result []WeeklyDataUsage
	for _, usage := range weeklyUsage {
		result = append(result, *usage)
	}
	
	// Sort by week start
	sort.Slice(result, func(i, j int) bool {
		return result[i].WeekStart.Before(result[j].WeekStart)
	})
	
	return result, nil
}

// GetMonthlyDataUsage returns monthly data usage for a given time window
func (am *AnalyticsManager) GetMonthlyDataUsage(window string) ([]MonthlyDataUsage, error) {
	dailyUsage, err := am.GetDailyDataUsage(window)
	if err != nil {
		return nil, err
	}
	
	// Group by month
	monthlyUsage := make(map[string]*MonthlyDataUsage)
	for _, daily := range dailyUsage {
		monthStart := beginningOfMonth(daily.Date)
		monthKey := monthStart.Format("2006-01")
		
		if _, exists := monthlyUsage[monthKey]; !exists {
			monthlyUsage[monthKey] = &MonthlyDataUsage{
				MonthStart: monthStart,
				MonthEnd:   monthStart.AddDate(0, 1, 0),
			}
		}
		
		monthlyUsage[monthKey].DataSent += daily.DataSent
		monthlyUsage[monthKey].DataReceived += daily.DataReceived
	}
	
	// Convert to slice
	var result []MonthlyDataUsage
	for _, usage := range monthlyUsage {
		result = append(result, *usage)
	}
	
	// Sort by month start
	sort.Slice(result, func(i, j int) bool {
		return result[i].MonthStart.Before(result[j].MonthStart)
	})
	
	return result, nil
}

// GetServerPerformance returns performance metrics for all servers
func (am *AnalyticsManager) GetServerPerformance(window string) ([]ServerPerformance, error) {
	// Parse window (e.g., "7d", "30d")
	duration, err := parseDuration(window)
	if err != nil {
		return nil, err
	}
	
	// Calculate time range
	endTime := time.Now()
	startTime := endTime.Add(-duration)
	
	// Get all servers
	servers := am.serverManager.GetAllServers()
	
	var performance []ServerPerformance
	for _, server := range servers {
		// In a real implementation, you would track ping history over time
		// For now, we'll use the current ping value
		pingStats := PingStats{
			Average: float64(server.Ping),
			Max:     server.Ping,
			Min:     server.Ping,
			P95:     server.Ping,
			Samples: 1,
		}
		
		perf := ServerPerformance{
			ServerID:   server.ID,
			ServerName: server.Name,
			PingStats:  pingStats,
		}
		
		performance = append(performance, perf)
	}
	
	return performance, nil
}

// Helper functions

// parseDuration parses a duration string like "7d" or "30d"
func parseDuration(durationStr string) (time.Duration, error) {
	if len(durationStr) < 2 {
		return 0, &InvalidDurationError{Duration: durationStr}
	}
	
	unit := durationStr[len(durationStr)-1:]
	valueStr := durationStr[:len(durationStr)-1]
	
	// Convert value string to int
	var value int
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		return 0, &InvalidDurationError{Duration: durationStr}
	}
	
	switch unit {
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	case "h":
		return time.Duration(value) * time.Hour, nil
	case "m":
		return time.Duration(value) * time.Minute, nil
	case "s":
		return time.Duration(value) * time.Second, nil
	default:
		return 0, &InvalidDurationError{Duration: durationStr}
	}
}

// formatDuration formats a duration to a string like "7d" or "30d"
func formatDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	days := hours / 24
	
	if days > 0 {
		return fmt.Sprintf("%dd", days)
	}
	
	return fmt.Sprintf("%dh", hours)
}

// beginningOfWeek returns the beginning of the week (Monday) for a given time
func beginningOfWeek(t time.Time) time.Time {
	// Calculate days to subtract to get to Monday
	offset := int(t.Weekday() - time.Monday)
	if offset < 0 {
		offset += 7
	}
	
	// Return beginning of the week
	return time.Date(t.Year(), t.Month(), t.Day()-offset, 0, 0, 0, 0, t.Location())
}

// beginningOfMonth returns the beginning of the month for a given time
func beginningOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// InvalidPeriodError represents an invalid report period error
type InvalidPeriodError struct {
	Period string
}

func (e *InvalidPeriodError) Error() string {
	return fmt.Sprintf("invalid report period: %s", e.Period)
}

// InvalidDurationError represents an invalid duration error
type InvalidDurationError struct {
	Duration string
}

func (e *InvalidDurationError) Error() string {
	return fmt.Sprintf("invalid duration format: %s", e.Duration)
}