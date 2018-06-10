package clientutil

import (
	"time"

	rc "github.com/grokify/go-ringcentral/client"
)

type SubscriptionManager struct {
	Client       *rc.APIClient
	EventFilters []string
	subscription SubscriptionInfo
}

func NewSubscriptionManager(apiClient *rc.APIClient) SubscriptionManager {
	sub := SubscriptionManager{
		Client:       apiClient,
		EventFilters: []string{},
		subscription: SubscriptionInfo{},
	}
	return sub
}

type SubscriptionInfo struct {
	EventFilters   []string
	SubscriptionId string
	DeliveryMode   DeliveryMode
	CreationTime   time.Time
	ExpirationTime time.Time
	ExpiresIn      int64
	Status         string
	URI            string
}

type DeliveryMode struct {
	TransportType string
	Encryption    bool
	Address       string
	SubscriberKey string
	SecretKey     string
}
