package managers

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// DataManager مدیریت داده‌های مصرفی
type DataManager struct {
	dataFile string
	mutex    sync.RWMutex
	data     map[string]*ServerData
}

// ServerData داده‌های مصرفی یک سرور
type ServerData struct {
	ServerID    string `json:"server_id"`
	TotalSent   int64  `json:"total_sent"`
	TotalRecv   int64  `json:"total_received"`
	LastUpdate  time.Time `json:"last_update"`
}

// NewDataManager ایجاد یک مدیر داده جدید
func NewDataManager(dataPath string) *DataManager {
	// اطمینان از وجود دایرکتوری داده
	dir := filepath.Dir(dataPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	dm := &DataManager{
		dataFile: dataPath,
		data:     make(map[string]*ServerData),
	}

	// بارگذاری داده‌های موجود
	dm.loadData()
	return dm
}

// loadData بارگذاری داده‌ها از فایل
func (dm *DataManager) loadData() error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// بررسی وجود فایل
	if _, err := os.Stat(dm.dataFile); os.IsNotExist(err) {
		return nil // فایل وجود ندارد، داده‌های خالی استفاده می‌شود
	}

	// خواندن فایل
	data, err := os.ReadFile(dm.dataFile)
	if err != nil {
		return fmt.Errorf("خطا در خواندن فایل داده: %v", err)
	}

	// تجزیه JSON
	var serverDataList []*ServerData
	err = json.Unmarshal(data, &serverDataList)
	if err != nil {
		return fmt.Errorf("خطا در تجزیه داده‌ها: %v", err)
	}

	// تبدیل به map برای دسترسی سریع
	dm.data = make(map[string]*ServerData)
	for _, sd := range serverDataList {
		dm.data[sd.ServerID] = sd
	}

	return nil
}

// saveData ذخیره داده‌ها در فایل
func (dm *DataManager) saveData() error {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	// تبدیل map به slice
	serverDataList := make([]*ServerData, 0, len(dm.data))
	for _, sd := range dm.data {
		serverDataList = append(serverDataList, sd)
	}

	// تبدیل به JSON
	data, err := json.MarshalIndent(serverDataList, "", "  ")
	if err != nil {
		return fmt.Errorf("خطا در تبدیل داده‌ها به JSON: %v", err)
	}

	// نوشتن در فایل
	err = os.WriteFile(dm.dataFile, data, 0644)
	if err != nil {
		return fmt.Errorf("خطا در نوشتن فایل داده: %v", err)
	}

	return nil
}

// RecordDataUsage ثبت مصرف داده برای یک سرور
func (dm *DataManager) RecordDataUsage(serverID string, sent, received int64) error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// گرفتن یا ایجاد داده سرور
	sd, exists := dm.data[serverID]
	if !exists {
		sd = &ServerData{
			ServerID:   serverID,
			TotalSent:  0,
			TotalRecv:  0,
			LastUpdate: time.Now(),
		}
		dm.data[serverID] = sd
	}

	// به‌روزرسانی داده‌ها
	sd.TotalSent += sent
	sd.TotalRecv += received
	sd.LastUpdate = time.Now()

	// ذخیره داده‌ها
	return dm.saveData()
}

// GetServerData دریافت داده‌های مصرفی یک سرور
func (dm *DataManager) GetServerData(serverID string) (*ServerData, error) {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	sd, exists := dm.data[serverID]
	if !exists {
		return nil, fmt.Errorf("داده‌ای برای سرور %s یافت نشد", serverID)
	}

	// بازگشت یک کپی برای جلوگیری از تغییرات خارجی
	dataCopy := *sd
	return &dataCopy, nil
}

// GetAllData دریافت تمام داده‌های مصرفی
func (dm *DataManager) GetAllData() map[string]*ServerData {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	// ایجاد کپی از داده‌ها
	dataCopy := make(map[string]*ServerData)
	for k, v := range dm.data {
		dataCopy[k] = &ServerData{
			ServerID:    v.ServerID,
			TotalSent:   v.TotalSent,
			TotalRecv:   v.TotalRecv,
			LastUpdate:  v.LastUpdate,
		}
	}

	return dataCopy
}

// ResetServerData بازنشانی داده‌های یک سرور
func (dm *DataManager) ResetServerData(serverID string) error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	delete(dm.data, serverID)
	return dm.saveData()
}

// ResetAllData بازنشانی تمام داده‌ها
func (dm *DataManager) ResetAllData() error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dm.data = make(map[string]*ServerData)
	return dm.saveData()
}