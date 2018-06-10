/*
 * RingCentral Connect Platform API Explorer
 *
 * <p>This is a beta interactive API explorer for the RingCentral Connect Platform. To use this service, you will need to have an account with the proper credentials to generate an OAuth2 access token.</p><p><h2>Quick Start</h2></p><ol><li>1) Go to <b>Authentication > /oauth/token</b></li><li>2) Enter <b>app_key, app_secret, username, password</b> fields and then click \"Try it out!\"</li><li>3) Upon success, your access_token is loaded and you can access any form requiring authorization.</li></ol><h2>Links</h2><ul><li><a href=\"https://github.com/ringcentral\" target=\"_blank\">RingCentral SDKs on Github</a></li><li><a href=\"mailto:devsupport@ringcentral.com\">RingCentral Developer Support Email</a></li></ul>
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package ringcentral

import (
	"time"
)

type GlipApnsInfo struct {

	// Apple Push Notification Service Info
	Aps *ApsInfo `json:"aps,omitempty"`

	// Datetime of a call action in ISO 8601 format including timezone, for example 2016-03-10T18:07:52.534Z
	Timestamp time.Time `json:"timestamp,omitempty"`

	// Universally unique identifier of a notification
	Uuid string `json:"uuid,omitempty"`

	// Event filter URI
	Event string `json:"event,omitempty"`

	// Internal identifier of a subscription
	SubscriptionId string `json:"subscriptionId,omitempty"`

	// Unread messages data
	Body *GlipUnreadMessageCountInfo `json:"body,omitempty"`
}
