package alert

import (
	"fmt"
	"sync"
	"time"
	"c:/Users/behza/OneDrive/Documents/vpn/src/managers"
	"c:/Users/behza/OneDrive/Documents/vpn/src/history"
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
)

// AlertRepository defines the interface for alert data storage
type AlertRepository interface {
	AddAlertRecord(record Alert) error
	GetAlertRecords(unread, unresolved bool, limit int) ([]Alert, error)
	UpdateAlertRecord(record Alert) error
}

// AlertManager manages alert rules and notifications
type AlertManager struct {
	config        AlertManagerConfig
	serverManager *managers.ServerManager
	connManager   *managers.ConnectionManager
	dataManager   *managers.DataManager
	historyManager history.HistoryManager
	repository     AlertRepository
	rules         []AlertRule
	handlers      []AlertHandler
	mutex         sync.RWMutex
	running       bool
	stopChan      chan struct{}
}

// NewAlertManager creates a new alert manager
func NewAlertManager(
	serverMgr *managers.ServerManager,
	connMgr *managers.ConnectionManager,
	dataMgr *managers.DataManager,
	historyMgr history.HistoryManager,
	config AlertManagerConfig) *AlertManager {
	
	am := &AlertManager{
		config:        config,
		serverManager: serverMgr,
		connManager:   connMgr,
		dataManager:   dataMgr,
		historyManager: historyMgr,
		rules:         make([]AlertRule, 0),
		handlers:      make([]AlertHandler, 0),
		stopChan:      make(chan struct{}),
	}
	
	// Load default rules
	am.loadDefaultRules()
	
	return am
}

// SetAlertRepository sets the alert repository for data storage
func (am *AlertManager) SetAlertRepository(repo AlertRepository) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	am.repository = repo
}

// loadDefaultRules loads default alert rules
func (am *AlertManager) loadDefaultRules() {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	// Data usage rules
	am.rules = append(am.rules, AlertRule{
		ID:        utils.GenerateID(),
		Name:      "Data Usage 80%",
		Type:      RuleTypeDataUsage,
		Threshold: 80,
		Enabled:   true,
		NotifyDesktop: true,
		NotifyUI:      true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	
	am.rules = append(am.rules, AlertRule{
		ID:        utils.GenerateID(),
		Name:      "Data Usage 95%",
		Type:      RuleTypeDataUsage,
		Threshold: 95,
		Enabled:   true,
		NotifyDesktop: true,
		NotifyUI:      true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	
	// High ping rule
	am.rules = append(am.rules, AlertRule{
		ID:        utils.GenerateID(),
		Name:      "High Ping",
		Type:      RuleTypeHighPing,
		Threshold: 200,
		Enabled:   true,
		NotifyDesktop: true,
		NotifyUI:      true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	
	// Connection loss rule
	am.rules = append(am.rules, AlertRule{
		ID:        utils.GenerateID(),
		Name:      "Connection Loss",
		Type:      RuleTypeConnectionLoss,
		Duration:  30,
		Enabled:   true,
		NotifyDesktop: true,
		NotifyUI:      true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

// Start starts the alert manager
func (am *AlertManager) Start() {
	am.mutex.Lock()
	if am.running {
		am.mutex.Unlock()
		return
	}
	am.running = true
	am.mutex.Unlock()
	
	go func() {
		ticker := time.NewTicker(time.Duration(am.config.EvaluationInterval) * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				am.evaluateRules()
			case <-am.stopChan:
				return
			}
		}
	}()
}

// Stop stops the alert manager
func (am *AlertManager) Stop() {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	if !am.running {
		return
	}
	
	am.running = false
	close(am.stopChan)
}

// evaluateRules evaluates all enabled alert rules
func (am *AlertManager) evaluateRules() {
	am.mutex.RLock()
	rules := make([]AlertRule, len(am.rules))
	copy(rules, am.rules)
	am.mutex.RUnlock()
	
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		
		switch rule.Type {
		case RuleTypeDataUsage:
			am.evaluateDataUsageRule(rule)
		case RuleTypeHighPing:
			am.evaluateHighPingRule(rule)
		case RuleTypeConnectionLoss:
			am.evaluateConnectionLossRule(rule)
		}
	}
}

// evaluateDataUsageRule evaluates a data usage rule
func (am *AlertManager) evaluateDataUsageRule(rule AlertRule) {
	servers := am.serverManager.GetAllServers()
	
	for _, server := range servers {
		if server.DataLimit <= 0 {
			continue
		}
		
		usagePercent := float64(server.DataUsed) / float64(server.DataLimit) * 100
		if usagePercent >= rule.Threshold {
			// Determine severity based on threshold
			severity := SeverityWarning
			if rule.Threshold >= 90 {
				severity = SeverityError
			}
			
			alert := Alert{
				ID:         utils.GenerateID(),
				RuleID:     rule.ID,
				Type:       RuleTypeDataUsage,
				Title:      rule.Name,
				Message:    fmt.Sprintf("Data usage on server %s reached %.2f%%", server.Name, usagePercent),
				Timestamp:  time.Now(),
				Value:      usagePercent,
				Threshold:  rule.Threshold,
				Severity:   severity,
				Resolved:   false,
				Read:       false,
				ServerID:   server.ID,
				ServerName: server.Name,
			}
			
			am.sendAlert(&alert)
		}
	}
}

// evaluateHighPingRule evaluates a high ping rule
func (am *AlertManager) evaluateHighPingRule(rule AlertRule) {
	servers := am.serverManager.GetAllServers()
	
	for _, server := range servers {
		if server.Ping >= int(rule.Threshold) {
			// Determine severity based on ping value
			severity := SeverityWarning
			if server.Ping >= 300 {
				severity = SeverityError
			}
			
			alert := Alert{
				ID:         utils.GenerateID(),
				RuleID:     rule.ID,
				Type:       RuleTypeHighPing,
				Title:      rule.Name,
				Message:    fmt.Sprintf("Server %s ping is %d ms", server.Name, server.Ping),
				Timestamp:  time.Now(),
				Value:      float64(server.Ping),
				Threshold:  rule.Threshold,
				Severity:   severity,
				Resolved:   false,
				Read:       false,
				ServerID:   server.ID,
				ServerName: server.Name,
			}
			
			am.sendAlert(&alert)
		}
	}
}

// evaluateConnectionLossRule evaluates a connection loss rule
func (am *AlertManager) evaluateConnectionLossRule(rule AlertRule) {
	status := am.connManager.GetStatus()
	
	// Check if disconnected and for how long
	if status.State == core.ConnectionStateDisconnected {
		// In a real implementation, you would track disconnection time
		// For now, we'll just send an alert if disconnected
		if rule.Duration <= 0 {
			rule.Duration = 30 // Default to 30 seconds
		}
		
		alert := Alert{
			ID:        utils.GenerateID(),
			RuleID:    rule.ID,
			Type:      RuleTypeConnectionLoss,
			Title:     rule.Name,
			Message:   "VPN connection has been lost",
			Timestamp: time.Now(),
			Value:     float64(rule.Duration),
			Severity:  SeverityError,
			Resolved:  false,
			Read:      false,
		}
		
		am.sendAlert(&alert)
	}
}

// sendAlert sends an alert to all handlers and stores it
func (am *AlertManager) sendAlert(alert *Alert) {
	// Store alert in repository if available
	if am.repository != nil {
		if err := am.repository.AddAlertRecord(*alert); err != nil {
			// Log error in a real implementation
		}
	}
	
	// Send to handlers
	am.mutex.RLock()
	handlers := make([]AlertHandler, len(am.handlers))
	copy(handlers, am.handlers)
	am.mutex.RUnlock()
	
	for _, handler := range handlers {
		handler.HandleAlert(alert)
	}
}

// AddAlertHandler adds an alert handler
func (am *AlertManager) AddAlertHandler(handler AlertHandler) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	am.handlers = append(am.handlers, handler)
}

// GetAlerts returns alerts based on filters
func (am *AlertManager) GetAlerts(unread, unresolved bool, limit int) ([]Alert, error) {
	// If repository is available, use it
	if am.repository != nil {
		return am.repository.GetAlertRecords(unread, unresolved, limit)
	}
	
	// Fallback to in-memory alerts (in a real implementation)
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	
	// Return empty slice as fallback
	return []Alert{}, nil
}

// UpdateAlert updates an alert
func (am *AlertManager) UpdateAlert(alert Alert) error {
	// If repository is available, use it
	if am.repository != nil {
		return am.repository.UpdateAlertRecord(alert)
	}
	
	// Fallback implementation (in a real implementation)
	return nil
}

// GetAlertRules returns all alert rules
func (am *AlertManager) GetAlertRules() []AlertRule {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	
	rules := make([]AlertRule, len(am.rules))
	copy(rules, am.rules)
	
	return rules
}

// GetAlertRule returns a specific alert rule by ID
func (am *AlertManager) GetAlertRule(id string) (AlertRule, error) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	
	for _, rule := range am.rules {
		if rule.ID == id {
			return rule, nil
		}
	}
	
	return AlertRule{}, ErrRuleNotFound
}

// AddAlertRule adds a new alert rule
func (am *AlertManager) AddAlertRule(rule AlertRule) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	// Check if rule with same ID already exists
	for _, existingRule := range am.rules {
		if existingRule.ID == rule.ID {
			return ErrRuleAlreadyExists
		}
	}
	
	am.rules = append(am.rules, rule)
	return nil
}

// UpdateAlertRule updates an existing alert rule
func (am *AlertManager) UpdateAlertRule(rule AlertRule) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	for i, existingRule := range am.rules {
		if existingRule.ID == rule.ID {
			am.rules[i] = rule
			return nil
		}
	}
	
	return ErrRuleNotFound
}

// DeleteAlertRule deletes an alert rule by ID
func (am *AlertManager) DeleteAlertRule(id string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	for i, rule := range am.rules {
		if rule.ID == id {
			// Remove the rule
			am.rules = append(am.rules[:i], am.rules[i+1:]...)
			return nil
		}
	}
	
	return ErrRuleNotFound
}

// MarkAlertAsRead marks an alert as read
func (am *AlertManager) MarkAlertAsRead(alertID string) error {
	// In a real implementation, you would update the alert in the repository
	// For now, we'll just return nil
	return nil
}

// ResolveAlert resolves an alert
func (am *AlertManager) ResolveAlert(alertID string) error {
	// In a real implementation, you would update the alert in the repository
	// For now, we'll just return nil
	return nil
}

// ExportRules exports alert rules as JSON
func (am *AlertManager) ExportRules() ([]byte, error) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	
	return utils.ToJSON(am.rules)
}

// ImportRules imports alert rules from JSON
func (am *AlertManager) ImportRules(data []byte) error {
	var rules []AlertRule
	if err := utils.FromJSON(data, &rules); err != nil {
		return err
	}
	
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	// Replace all rules
	am.rules = rules
	return nil
}