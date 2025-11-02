package database

import (
	"vpnclient/src/core"
)

// SubscriptionStore handles subscription database operations
// For now, JSONStore implements these methods directly
// This file is kept for future SQLite implementation

// NewSubscriptionStore creates a new subscription store
func NewSubscriptionStore(store Store) *SubscriptionStoreWrapper {
	return &SubscriptionStoreWrapper{store: store}
}

// SubscriptionStoreWrapper wraps Store to provide SubscriptionStore interface
type SubscriptionStoreWrapper struct {
	store Store
}

// AddSubscription adds a new subscription to the database
func (s *SubscriptionStoreWrapper) AddSubscription(sub *core.Subscription) error {
	return s.store.AddSubscription(sub)
}

// GetSubscription retrieves a subscription by ID
func (s *SubscriptionStoreWrapper) GetSubscription(id string) (*core.Subscription, error) {
	result, err := s.store.GetSubscription(id)
	if err != nil {
		return nil, err
	}
	return result.(*core.Subscription), nil
}

// GetAllSubscriptions retrieves all subscriptions
func (s *SubscriptionStoreWrapper) GetAllSubscriptions() ([]*core.Subscription, error) {
	results, err := s.store.GetAllSubscriptions()
	if err != nil {
		return nil, err
	}
	subs := make([]*core.Subscription, len(results))
	for i, r := range results {
		subs[i] = r.(*core.Subscription)
	}
	return subs, nil
}

// UpdateSubscription updates an existing subscription
func (s *SubscriptionStoreWrapper) UpdateSubscription(sub *core.Subscription) error {
	return s.store.UpdateSubscription(sub)
}

// DeleteSubscription deletes a subscription by ID
func (s *SubscriptionStoreWrapper) DeleteSubscription(id string) error {
	return s.store.DeleteSubscription(id)
}
