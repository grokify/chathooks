package clientutil

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	InstantMessageSMSExample   = "/restapi/v1.0/account/~/extension/12345678/message-store/instant?type=SMS"
	InstantMessageSMSPattern   = `message-store/instant.+type=SMS`
	InstantMessageSMSPatternX  = `message-store/instant`
	GlipPostEventFilterPattern = `/restapi/v1.0/glip/posts`
	SubscriptionRenewalFilter  = `/restapi/v1.0/subscription/.+\?threshold=`
)

type EventType int

const (
	AccountPresenceEvent EventType = iota
	ContactDirectoryEvent
	DetailedExtensionPresenceEvent
	DetailedExtensionPresenceWithSIPEvent
	ExtensionFavoritesEvent
	ExtensionFavoritesPresenceEvent
	ExtensionGrantListEvent
	ExtensionListEvent
	ExtensionInfoEvent
	ExtensionPresenceEvent
	ExtensionPresenceLineEvent
	GlipGroupsEvent
	GlipPostEvent
	GlipUnreadMessageCountEvent
	InboundMessageEvent
	IncomingCallEvent
	InstantMessageEvent
	MessageEvent
	MissedCallEvent
	RCVideoNotificationsEvent
	SubscriptionRenewalEvent
)

// Events is an array of event structs for reference.
var Events = []string{
	"AccountPresenceEvent",
	"ContactDirectoryEvent",
	"DetailedExtensionPresenceEvent",
	"DetailedExtensionPresenceWithSIPEvent",
	"ExtensionFavoritesEvent",
	"ExtensionFavoritesPresenceEvent",
	"ExtensionGrantListEvent",
	"ExtensionListEvent",
	"ExtensionInfoEvent",
	"ExtensionPresenceEvent",
	"ExtensionPresenceLineEvent",
	"GlipGroupsEvent",
	"GlipPostEvent",
	"GlipUnreadMessageCountEvent",
	"InboundMessageEvent",
	"IncomingCallEvent",
	"InstantMessageEvent",
	"MessageEvent",
	"MissedCallEvent",
	"RCVideoNotificationsEvent",
	"SubscriptionRenewalEvent",
}

func (d EventType) String() string { return Events[d] }

/*
func IsInstantMessageSMS(s string) bool {
	if strings.Index(s, InstantMessageSMSPattern) == -1 {
		return false
	}
	return true
}
*/
func ParseEventTypeForFilter(eventFilter string) (EventType, error) {
	if strings.Index(eventFilter, GlipPostEventFilterPattern) > -1 {
		return GlipPostEvent, nil
	}
	rx := regexp.MustCompile(SubscriptionRenewalFilter)
	m := rx.FindString(eventFilter)
	if len(m) > 0 {
		return SubscriptionRenewalEvent, nil
	}
	fmt.Printf("EVT_FILTER: %v\n", eventFilter)
	m2 := regexp.MustCompile(InstantMessageSMSPattern).FindString(eventFilter)
	if len(m2) > 0 {
		fmt.Println("HERE")
		return InstantMessageEvent, nil
	}
	fmt.Println("NODICE")

	return GlipPostEvent, fmt.Errorf("No Event found for filter %v", eventFilter)
}
