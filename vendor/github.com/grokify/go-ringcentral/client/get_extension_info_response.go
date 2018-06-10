/*
 * RingCentral Connect Platform API Explorer
 *
 * <p>This is a beta interactive API explorer for the RingCentral Connect Platform. To use this service, you will need to have an account with the proper credentials to generate an OAuth2 access token.</p><p><h2>Quick Start</h2></p><ol><li>1) Go to <b>Authentication > /oauth/token</b></li><li>2) Enter <b>app_key, app_secret, username, password</b> fields and then click \"Try it out!\"</li><li>3) Upon success, your access_token is loaded and you can access any form requiring authorization.</li></ol><h2>Links</h2><ul><li><a href=\"https://github.com/ringcentral\" target=\"_blank\">RingCentral SDKs on Github</a></li><li><a href=\"mailto:devsupport@ringcentral.com\">RingCentral Developer Support Email</a></li></ul>
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package ringcentral

type GetExtensionInfoResponse struct {

	// Internal identifier of an extension
	Id int64 `json:"id"`

	// Canonical URI of an extension
	Uri string `json:"uri"`

	// Contact detailed information
	Contact *ContactInfo `json:"contact,omitempty"`

	// Information on department extension(s), to which the requested extension belongs. Returned only for user extensions, members of department, requested by single extensionId
	Departments []DepartmentInfo `json:"departments,omitempty"`

	// Number of department extension
	ExtensionNumber string `json:"extensionNumber,omitempty"`

	// Extension user name
	Name string `json:"name,omitempty"`

	// For Partner Applications Internal identifier of an extension created by partner. The RingCentral supports the mapping of accounts and stores the corresponding account ID/extension ID for each partner ID of a client application. In request URIs partner IDs are accepted instead of regular RingCentral native IDs as path parameters using pid = XXX clause. Though in response URIs contain the corresponding account IDs and extension IDs. In all request and response bodies these values are reflected via partnerId attributes of account and extension
	PartnerId string `json:"partnerId,omitempty"`

	// Extension permissions, corresponding to the Service Web permissions 'Admin' and 'InternationalCalling'
	Permissions *ExtensionPermissions `json:"permissions,omitempty"`

	// Information on profile image
	ProfileImage *ProfileImageInfo `json:"profileImage"`

	// List of non-RC internal identifiers assigned to an extension
	References []ReferenceInfo `json:"references,omitempty"`

	// Extension region data (timezone, home country, language)
	RegionalSettings *RegionalSettings `json:"regionalSettings,omitempty"`

	// Extension service features returned in response only when the logged-in user requests his/her own extension info, see also Extension Service Features
	ServiceFeatures []ExtensionServiceFeatureInfo `json:"serviceFeatures,omitempty"`

	// Specifies extension configuration wizard state (web service setup). The default value is 'NotStarted'
	SetupWizardState string `json:"setupWizardState,omitempty"`

	// Extension current state. If the status is 'Unassigned'. Returned for all extensions
	Status string `json:"status"`

	// Status information (reason, comment). Returned for 'Disabled' status only
	StatusInfo *ExtensionStatusInfo `json:"statusInfo,omitempty"`

	// Extension type
	Type_ string `json:"type"`

	// For Department extension type only. Call queue settings
	CallQueueInfo *CallQueueInfo `json:"callQueueInfo,omitempty"`
}
