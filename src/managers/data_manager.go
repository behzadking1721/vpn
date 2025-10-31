package managers

import (
	"vpn-client/src/core"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// DataManager handles data persistence for servers and subscriptions
type DataManager struct {
	serversFile       string
	subscriptionsFile string
}

// NewDataManager creates a new data manager
func NewDataManager(serversFile, subscriptionsFile string) *DataManager {
	return &DataManager{
		serversFile:       serversFile,
		subscriptionsFile: subscriptionsFile,
	}
}

// SaveServers saves servers to a JSON file
func (dm *DataManager) SaveServers(servers []core.Server) error {
	data, err := json.MarshalIndent(servers, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal servers: %v", err)
	}

	err = ioutil.WriteFile(dm.serversFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write servers file: %v", err)
	}

	return nil
}

// LoadServers loads servers from a JSON file
func (dm *DataManager) LoadServers() ([]core.Server, error) {
	// Check if file exists
	if _, err := os.Stat(dm.serversFile); os.IsNotExist(err) {
		// Return empty slice if file doesn't exist
		return []core.Server{}, nil
	}

	data, err := ioutil.ReadFile(dm.serversFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read servers file: %v", err)
	}

	var servers []core.Server
	err = json.Unmarshal(data, &servers)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal servers: %v", err)
	}

	return servers, nil
}

// SaveSubscriptions saves subscriptions to a JSON file
func (dm *DataManager) SaveSubscriptions(subscriptions []core.Subscription) error {
	data, err := json.MarshalIndent(subscriptions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal subscriptions: %v", err)
	}

	err = ioutil.WriteFile(dm.subscriptionsFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write subscriptions file: %v", err)
	}

	return nil
}

// LoadSubscriptions loads subscriptions from a JSON file
func (dm *DataManager) LoadSubscriptions() ([]core.Subscription, error) {
	// Check if file exists
	if _, err := os.Stat(dm.subscriptionsFile); os.IsNotExist(err) {
		// Return empty slice if file doesn't exist
		return []core.Subscription{}, nil
	}

	data, err := ioutil.ReadFile(dm.subscriptionsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read subscriptions file: %v", err)
	}

	var subscriptions []core.Subscription
	err = json.Unmarshal(data, &subscriptions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscriptions: %v", err)
	}

	return subscriptions, nil
}
