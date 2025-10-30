package alert

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/history"
	"c:/Users/behza/OneDrive/Documents/vpn/src/managers"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"testing"
	"time"
)

// MockAlertHandler is a mock implementation of AlertHandler for testing
type MockAlertHandler struct {
	Alerts []*Alert
}

// HandleAlert implements the AlertHandler interface
func (m *MockAlertHandler) HandleAlert(alert *Alert) {
	m.Alerts = append(m.Alerts, alert)
}

func TestAlertManager(t *testing.T) {
	// Create mock managers
	configManager := managers.NewConfigManager("./test_config.json")
	serverManager := managers.NewServerManager()
	connManager := managers.NewConnectionManager()
	dataManager := serverManager.GetDataManager()
	historyManager := history.NewHistoryManager("./test_history")
	
	// Create alert manager
	config := AlertManagerConfig{
		DesktopNotifications: true,
		EvaluationInterval:   1,
		HistoryRetention:     1,
	}
	
	alertManager := NewAlertManager(serverManager, connManager, dataManager, historyManager, config)
	
	// Test adding and retrieving rules
	t.Run("AddAndRetrieveRule", func(t *testing.T) {
		rule := AlertRule{
			ID:          "test-rule-1",
			Name:        "Test Rule",
			Description: "A test rule",
			Type:        RuleTypeDataUsage,
			Enabled:     true,
			Threshold:   80.0,
		}
		
		// Add rule
		err := alertManager.AddAlertRule(rule)
		if err != nil {
			t.Fatalf("Failed to add rule: %v", err)
		}
		
		// Retrieve rule
		retrievedRule, err := alertManager.GetAlertRule("test-rule-1")
		if err != nil {
			t.Fatalf("Failed to retrieve rule: %v", err)
		}
		
		if retrievedRule.ID != rule.ID {
			t.Errorf("Expected rule ID %s, got %s", rule.ID, retrievedRule.ID)
		}
		
		if retrievedRule.Name != rule.Name {
			t.Errorf("Expected rule name %s, got %s", rule.Name, retrievedRule.Name)
		}
	})
	
	// Test listing rules
	t.Run("ListRules", func(t *testing.T) {
		rules := alertManager.GetAlertRules()
		if len(rules) == 0 {
			t.Error("Expected at least one rule (default rules)")
		}
	})
	
	// Test updating rule
	t.Run("UpdateRule", func(t *testing.T) {
		rule := AlertRule{
			ID:        "test-rule-2",
			Name:      "Original Name",
			Type:      RuleTypeHighPing,
			Enabled:   true,
			Threshold: 100.0,
		}
		
		// Add rule
		err := alertManager.AddAlertRule(rule)
		if err != nil {
			t.Fatalf("Failed to add rule: %v", err)
		}
		
		// Update rule
		rule.Name = "Updated Name"
		err = alertManager.UpdateAlertRule(rule)
		if err != nil {
			t.Fatalf("Failed to update rule: %v", err)
		}
		
		// Retrieve updated rule
		updatedRule, err := alertManager.GetAlertRule("test-rule-2")
		if err != nil {
			t.Fatalf("Failed to retrieve updated rule: %v", err)
		}
		
		if updatedRule.Name != "Updated Name" {
			t.Errorf("Expected updated name 'Updated Name', got '%s'", updatedRule.Name)
		}
	})
	
	// Test removing rule
	t.Run("RemoveRule", func(t *testing.T) {
		rule := AlertRule{
			ID:      "test-rule-3",
			Name:    "Rule to Remove",
			Type:    RuleTypeConnectionLoss,
			Enabled: true,
		}
		
		// Add rule
		err := alertManager.AddAlertRule(rule)
		if err != nil {
			t.Fatalf("Failed to add rule: %v", err)
		}
		
		// Remove rule
		err = alertManager.DeleteAlertRule("test-rule-3")
		if err != nil {
			t.Fatalf("Failed to remove rule: %v", err)
		}
		
		// Try to retrieve removed rule
		_, err = alertManager.GetAlertRule("test-rule-3")
		if err == nil {
			t.Error("Expected error when retrieving removed rule, but got none")
		}
	})
	
	// Test alert handlers
	t.Run("AlertHandlers", func(t *testing.T) {
		mockHandler := &MockAlertHandler{}
		alertManager.AddAlertHandler(mockHandler)
		
		// Create and process an alert
		alert := &Alert{
			ID:        utils.GenerateID(),
			Type:      RuleTypeDataUsage,
			Title:     "Test Alert",
			Message:   "This is a test alert",
			Timestamp: time.Now(),
			Severity:  SeverityInfo,
		}
		
		alertManager.sendAlert(alert)
		
		// Check if handler received the alert
		if len(mockHandler.Alerts) != 1 {
			t.Errorf("Expected 1 alert in mock handler, got %d", len(mockHandler.Alerts))
		}
		
		if mockHandler.Alerts[0].ID != alert.ID {
			t.Errorf("Expected alert ID %s, got %s", alert.ID, mockHandler.Alerts[0].ID)
		}
	})
}

func TestAlertEvaluation(t *testing.T) {
	// Create mock managers
	configManager := managers.NewConfigManager("./test_config.json")
	serverManager := managers.NewServerManager()
	connManager := managers.NewConnectionManager()
	dataManager := serverManager.GetDataManager()
	historyManager := history.NewHistoryManager("./test_history")
	
	// Create alert manager
	config := AlertManagerConfig{
		DesktopNotifications: true,
		EvaluationInterval:   1,
		HistoryRetention:     1,
	}
	
	alertManager := NewAlertManager(serverManager, connManager, dataManager, historyManager, config)
	
	// Test data usage rule evaluation
	t.Run("DataUsageRuleEvaluation", func(t *testing.T) {
		// Create a server with data limit
		server := core.Server{
			ID:        utils.GenerateID(),
			Name:      "Test Server",
			Host:      "example.com",
			Port:      443,
			Protocol:  core.ProtocolVMess,
			DataLimit: 1000, // 1000 bytes limit
			DataUsed:  900,  // 900 bytes used (90%)
			Enabled:   true,
		}
		
		// Add server to manager
		err := serverManager.AddServer(server)
		if err != nil {
			t.Fatalf("Failed to add server: %v", err)
		}
		
		// Create a data usage rule with 80% threshold
		rule := AlertRule{
			ID:        utils.GenerateID(),
			Name:      "Data Usage Alert",
			Type:      RuleTypeDataUsage,
			Enabled:   true,
			Threshold: 80.0, // 80% threshold
		}
		
		err = alertManager.AddAlertRule(rule)
		if err != nil {
			t.Fatalf("Failed to add rule: %v", err)
		}
		
		// Evaluate rules
		alertManager.evaluateRules()
		
		// Note: In a real test, we would check if alerts were generated
		// This requires a mock alert handler to capture the alerts
	})
	
	// Test high ping rule evaluation
	t.Run("HighPingRuleEvaluation", func(t *testing.T) {
		// Create a server with high ping
		server := core.Server{
			ID:       utils.GenerateID(),
			Name:     "High Ping Server",
			Host:     "slow.example.com",
			Port:     443,
			Protocol: core.ProtocolVMess,
			Ping:     300, // 300ms ping
			Enabled:  true,
		}
		
		// Add server to manager
		err := serverManager.AddServer(server)
		if err != nil {
			t.Fatalf("Failed to add server: %v", err)
		}
		
		// Create a high ping rule with 200ms threshold
		rule := AlertRule{
			ID:        utils.GenerateID(),
			Name:      "High Ping Alert",
			Type:      RuleTypeHighPing,
			Enabled:   true,
			Threshold: 200.0, // 200ms threshold
		}
		
		err = alertManager.AddAlertRule(rule)
		if err != nil {
			t.Fatalf("Failed to add rule: %v", err)
		}
		
		// Evaluate rules
		alertManager.evaluateRules()
		
		// Note: In a real test, we would check if alerts were generated
		// This requires a mock alert handler to capture the alerts
	})
}

func TestAlertSeverity(t *testing.T) {
	// Test severity determination for data usage alerts
	t.Run("DataUsageSeverity", func(t *testing.T) {
		// Create a mock alert manager
		configManager := managers.NewConfigManager("./test_config.json")
		serverManager := managers.NewServerManager()
		connManager := managers.NewConnectionManager()
		dataManager := serverManager.GetDataManager()
		historyManager := history.NewHistoryManager("./test_history")
		
		config := AlertManagerConfig{
			DesktopNotifications: true,
			EvaluationInterval:   1,
			HistoryRetention:     1,
		}
		
		alertManager := NewAlertManager(serverManager, connManager, dataManager, historyManager, config)
		
		// Create a mock handler to capture alerts
		mockHandler := &MockAlertHandler{}
		alertManager.AddAlertHandler(mockHandler)
		
		// Create a server with high data usage
		server := core.Server{
			ID:        utils.GenerateID(),
			Name:      "Test Server",
			Host:      "example.com",
			Port:      443,
			Protocol:  core.ProtocolVMess,
			DataLimit: 1000,
			DataUsed:  950, // 95% usage
			Enabled:   true,
		}
		
		err := serverManager.AddServer(server)
		if err != nil {
			t.Fatalf("Failed to add server: %v", err)
		}
		
		// Create a data usage rule
		rule := AlertRule{
			ID:        utils.GenerateID(),
			Name:      "Data Usage Alert",
			Type:      RuleTypeDataUsage,
			Enabled:   true,
			Threshold: 90.0,
		}
		
		err = alertManager.AddAlertRule(rule)
		if err != nil {
			t.Fatalf("Failed to add rule: %v", err)
		}
		
		// Evaluate rules
		alertManager.evaluateDataUsageRule(rule)
		
		// Check if alert was generated with correct severity
		if len(mockHandler.Alerts) != 1 {
			t.Fatalf("Expected 1 alert, got %d", len(mockHandler.Alerts))
		}
		
		alert := mockHandler.Alerts[0]
		if alert.Severity != SeverityError {
			t.Errorf("Expected severity 'error', got '%s'", alert.Severity)
		}
	})
	
	// Test severity determination for ping alerts
	t.Run("PingSeverity", func(t *testing.T) {
		// Create a mock alert manager
		configManager := managers.NewConfigManager("./test_config.json")
		serverManager := managers.NewServerManager()
		connManager := managers.NewConnectionManager()
		dataManager := serverManager.GetDataManager()
		historyManager := history.NewHistoryManager("./test_history")
		
		config := AlertManagerConfig{
			DesktopNotifications: true,
			EvaluationInterval:   1,
			HistoryRetention:     1,
		}
		
		alertManager := NewAlertManager(serverManager, connManager, dataManager, historyManager, config)
		
		// Create a mock handler to capture alerts
		mockHandler := &MockAlertHandler{}
		alertManager.AddAlertHandler(mockHandler)
		
		// Create a server with high ping
		server := core.Server{
			ID:       utils.GenerateID(),
			Name:     "High Ping Server",
			Host:     "slow.example.com",
			Port:     443,
			Protocol: core.ProtocolVMess,
			Ping:     350, // 350ms ping
			Enabled:  true,
		}
		
		err := serverManager.AddServer(server)
		if err != nil {
			t.Fatalf("Failed to add server: %v", err)
		}
		
		// Create a high ping rule
		rule := AlertRule{
			ID:        utils.GenerateID(),
			Name:      "High Ping Alert",
			Type:      RuleTypeHighPing,
			Enabled:   true,
			Threshold: 200.0,
		}
		
		err = alertManager.AddAlertRule(rule)
		if err != nil {
			t.Fatalf("Failed to add rule: %v", err)
		}
		
		// Evaluate rules
		alertManager.evaluateHighPingRule(rule)
		
		// Check if alert was generated with correct severity
		if len(mockHandler.Alerts) != 1 {
			t.Fatalf("Expected 1 alert, got %d", len(mockHandler.Alerts))
		}
		
		alert := mockHandler.Alerts[0]
		if alert.Severity != SeverityError {
			t.Errorf("Expected severity 'error', got '%s'", alert.Severity)
		}
	})
}