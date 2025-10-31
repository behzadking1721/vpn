package monitoring

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/database"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// HealthCheckResult represents the result of a health check
type HealthCheckResult struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Checks    map[string]CheckResult `json:"checks"`
}

// CheckResult represents the result of a single check
type CheckResult struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// HealthChecker defines the interface for health checking
type HealthChecker interface {
	Check() CheckResult
	Name() string
}

// DatabaseHealthChecker checks the database health
type DatabaseHealthChecker struct {
	dbManager *database.Manager
}

// NewDatabaseHealthChecker creates a new database health checker
func NewDatabaseHealthChecker(dbManager *database.Manager) *DatabaseHealthChecker {
	return &DatabaseHealthChecker{
		dbManager: dbManager,
	}
}

// Check performs the database health check
func (dhc *DatabaseHealthChecker) Check() CheckResult {
	result := CheckResult{
		Timestamp: time.Now(),
	}

	// Try to ping the database
	if dhc.dbManager != nil {
		// Get the underlying sql.DB
		db := dhc.dbManager.DB

		// Ping the database
		ctx, cancel := timeoutContext(5 * time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			result.Status = "error"
			result.Message = fmt.Sprintf("Database ping failed: %v", err)
		} else {
			result.Status = "ok"
			result.Message = "Database is accessible"
		}
	} else {
		result.Status = "error"
		result.Message = "Database manager is not initialized"
	}

	return result
}

// Name returns the name of the checker
func (dhc *DatabaseHealthChecker) Name() string {
	return "database"
}

// WebSocketHealthChecker checks the WebSocket health
type WebSocketHealthChecker struct {
	addr string
}

// NewWebSocketHealthChecker creates a new WebSocket health checker
func NewWebSocketHealthChecker(addr string) *WebSocketHealthChecker {
	return &WebSocketHealthChecker{
		addr: addr,
	}
}

// Check performs the WebSocket health check
func (whc *WebSocketHealthChecker) Check() CheckResult {
	result := CheckResult{
		Timestamp: time.Now(),
	}

	// Try to connect to the WebSocket endpoint
	// For simplicity, we'll just check if we can connect to the address
	conn, err := net.DialTimeout("tcp", whc.addr, 5*time.Second)
	if err != nil {
		result.Status = "error"
		result.Message = fmt.Sprintf("WebSocket connection failed: %v", err)
	} else {
		conn.Close()
		result.Status = "ok"
		result.Message = "WebSocket endpoint is accessible"
	}

	return result
}

// Name returns the name of the checker
func (whc *WebSocketHealthChecker) Name() string {
	return "websocket"
}

// SystemHealthChecker checks the system health
type SystemHealthChecker struct {
}

// NewSystemHealthChecker creates a new system health checker
func NewSystemHealthChecker() *SystemHealthChecker {
	return &SystemHealthChecker{}
}

// Check performs the system health check
func (shc *SystemHealthChecker) Check() CheckResult {
	result := CheckResult{
		Timestamp: time.Now(),
	}

	// Get system stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Create a human-readable message
	result.Status = "ok"
	result.Message = fmt.Sprintf(
		"System healthy - Goroutines: %d, Memory allocated: %d KB, GC cycles: %d",
		runtime.NumGoroutine(),
		m.Alloc/1024,
		m.NumGC,
	)

	return result
}

// Name returns the name of the checker
func (shc *SystemHealthChecker) Name() string {
	return "system"
}

// HTTPHealthChecker checks HTTP endpoints
type HTTPHealthChecker struct {
	client  *http.Client
	url     string
	name    string
	timeout time.Duration
}

// NewHTTPHealthChecker creates a new HTTP health checker
func NewHTTPHealthChecker(url, name string, timeout time.Duration) *HTTPHealthChecker {
	return &HTTPHealthChecker{
		client: &http.Client{
			Timeout: timeout,
		},
		url:     url,
		name:    name,
		timeout: timeout,
	}
}

// Check performs the HTTP health check
func (hhc *HTTPHealthChecker) Check() CheckResult {
	result := CheckResult{
		Timestamp: time.Now(),
	}

	// Make an HTTP request to the endpoint
	resp, err := hhc.client.Get(hhc.url)
	if err != nil {
		result.Status = "error"
		result.Message = fmt.Sprintf("HTTP request failed: %v", err)
	} else {
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			result.Status = "ok"
			result.Message = fmt.Sprintf("HTTP endpoint is accessible (status: %d)", resp.StatusCode)
		} else {
			result.Status = "error"
			result.Message = fmt.Sprintf("HTTP endpoint returned error status: %d", resp.StatusCode)
		}
	}

	return result
}

// Name returns the name of the checker
func (hhc *HTTPHealthChecker) Name() string {
	return hhc.name
}

// HealthManager manages health checks
type HealthManager struct {
	checkers []HealthChecker
	mutex    sync.RWMutex
}

// NewHealthManager creates a new health manager
func NewHealthManager() *HealthManager {
	return &HealthManager{
		checkers: make([]HealthChecker, 0),
	}
}

// AddChecker adds a health checker
func (hm *HealthManager) AddChecker(checker HealthChecker) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	hm.checkers = append(hm.checkers, checker)
}

// RemoveChecker removes a health checker by name
func (hm *HealthManager) RemoveChecker(name string) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	for i, checker := range hm.checkers {
		if checker.Name() == name {
			hm.checkers = append(hm.checkers[:i], hm.checkers[i+1:]...)
			return
		}
	}
}

// Check performs all health checks
func (hm *HealthManager) Check() *HealthCheckResult {
	hm.mutex.RLock()
	checkers := make([]HealthChecker, len(hm.checkers))
	copy(checkers, hm.checkers)
	hm.mutex.RUnlock()

	result := &HealthCheckResult{
		Timestamp: time.Now(),
		Checks:    make(map[string]CheckResult),
	}

	// Perform all checks concurrently
	var wg sync.WaitGroup
	results := make(chan struct {
		name   string
		result CheckResult
	}, len(checkers))

	for _, checker := range checkers {
		wg.Add(1)
		go func(c HealthChecker) {
			defer wg.Done()
			results <- struct {
				name   string
				result CheckResult
			}{c.Name(), c.Check()}
		}(checker)
	}

	// Close the results channel when all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	allHealthy := true
	for res := range results {
		result.Checks[res.name] = res.result
		if res.result.Status != "ok" {
			allHealthy = false
		}
	}

	// Set overall status
	if allHealthy {
		result.Status = "healthy"
	} else {
		result.Status = "unhealthy"
	}

	return result
}

// GetChecker returns a checker by name
func (hm *HealthManager) GetChecker(name string) HealthChecker {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	for _, checker := range hm.checkers {
		if checker.Name() == name {
			return checker
		}
	}

	return nil
}

// timeoutContext creates a context with timeout
func timeoutContext(timeout time.Duration) (cancel func()) {
	// This is a simplified version. In a real implementation,
	// you would use context.WithTimeout

	// For now, we just return a no-op cancel function
	return func() {}
}
