package google

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/oauth2more/scim"
)

const (
	GoogleAPIUserinfoURL   = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"
	GoogleAPIPlusPeopleURL = "https://www.googleapis.com/plus/v1/people/me"
	GoogleAPIEmailURL      = "https://www.googleapis.com/userinfo/email"

	UserinfoEmailScope   = "https://www.googleapis.com/auth/userinfo#email"
	UserinfoProfileScope = "https://www.googleapis.com/auth/userinfo.profile"

	SpreadsheetsReadonly     = "https://www.googleapis.com/auth/spreadsheets.readonly"
	SpreadsheetsReadonlyDesc = "Allows read-only access to the user's sheets and their properties."
	Spreadsheets             = "https://www.googleapis.com/auth/spreadsheets"
	SpreadsheetsDesc         = "Allows read/write access to the user's sheets and their properties."

	DriveReadonly     = "https://www.googleapis.com/auth/drive.readonly"
	DriveReadonlyDesc = "Allows read-only access to the user's file metadata and file content."
	DriveFile         = "https://www.googleapis.com/auth/drive.file"
	DriveFileDesc     = "Per-file access to files created or opened by the app."
	Drive             = "https://www.googleapis.com/auth/drive"
	DriveDesc         = "Full, permissive scope to access all of a user's files. Request this scope only when it is strictly necessary."
)

// ClientUtil is a client library to retrieve the /userinfo
// endpoint which is not included in the Google API Go Client.
// For other endpoints, please consider using The Google API Go
// Client: https://github.com/google/google-api-go-client
type ClientUtil struct {
	Client *http.Client
	User   GoogleUserinfo `json:"user,omitempty"`
}

func NewClientUtil(client *http.Client) ClientUtil {
	return ClientUtil{Client: client}
}

func (apiutil *ClientUtil) SetClient(client *http.Client) {
	apiutil.Client = client
}

// GetUserinfoEmail retrieves the user's email from the
// https://www.googleapis.com/userinfo/email endpoint.
func (apiutil *ClientUtil) GetUserinfoEmail() (GoogleUserinfoEmail, error) {
	resp, err := apiutil.Client.Get(GoogleAPIEmailURL)
	if err != nil {
		return GoogleUserinfoEmail{}, err
	}

	bodyBytes, err := httputilmore.ResponseBody(resp)
	if err != nil {
		return GoogleUserinfoEmail{}, err
	}

	// parse user query string
	return ParseGoogleUserinfoEmail(string(bodyBytes))
}

type GoogleUserinfoEmail struct {
	Email      string `json:"email,omitempty"`
	IsVerified bool   `json:"isVerified,omitempty"`
}

func ParseGoogleUserinfoEmail(query string) (GoogleUserinfoEmail, error) {
	// parse email=johncwang@gmail.com&isVerified=true
	params, err := url.ParseQuery(query)
	googleUserinfoEmail := GoogleUserinfoEmail{}
	if err != nil {
		return googleUserinfoEmail, err
	}
	googleUserinfoEmail.Email = strings.TrimSpace(params.Get("email"))

	isVerified := strings.ToLower(strings.TrimSpace(params.Get("isVerified")))
	if isVerified == "true" {
		googleUserinfoEmail.IsVerified = true
	} else {
		googleUserinfoEmail.IsVerified = false
	}

	return googleUserinfoEmail, nil
}

// GetUserinfo retrieves the userinfo from the
// https://www.googleapis.com/oauth2/v1/userinfo?alt=json
// endpoint.
func (apiutil *ClientUtil) GetUserinfo() (GoogleUserinfo, error) {
	resp, err := apiutil.Client.Get(GoogleAPIUserinfoURL)
	if err != nil {
		return GoogleUserinfo{}, err
	}

	bodyBytes, err := httputilmore.ResponseBody(resp)
	if err != nil {
		return GoogleUserinfo{}, err
	}

	userinfo := GoogleUserinfo{}
	err = json.Unmarshal(bodyBytes, &userinfo)
	if err == nil {
		apiutil.User = userinfo
	}
	return userinfo, err
}

type GoogleUserinfo struct {
	FamilyName string `json:"family_name,omitempty"`
	Gender     string `json:"gender,omitempty"`
	GivenName  string `json:"given_name,omitempty"`
	ID         string `json:"id,omitempty"`
	Link       string `json:"link,omitempty"`
	Locale     string `json:"locale,omitempty"`
	Name       string `json:"name,omitempty"`
	PictureURL string `json:"picture,omitempty"`
}

// GetPlusPerson retrieves the userinfo from the
// https://www.googleapis.com/oauth2/v1/userinfo?alt=json
// endpoint.
func (apiutil *ClientUtil) GetPlusPerson() (GooglePlusPerson, error) {
	resp, err := apiutil.Client.Get(GoogleAPIPlusPeopleURL)
	if err != nil {
		return GooglePlusPerson{}, err
	}

	bodyBytes, err := httputilmore.ResponseBody(resp)
	if err != nil {
		return GooglePlusPerson{}, err
	}

	plusPerson := GooglePlusPerson{}
	err = json.Unmarshal(bodyBytes, &plusPerson)
	return plusPerson, err
}

type GooglePlusPerson struct {
	Kind        string                `json:"kind,omitempty"`
	Etag        string                `json:"etag,omitempty"`
	Gender      string                `json:"gender,omitempty"`
	ObjectType  string                `json:"objectType,omitempty"`
	ID          string                `json:"id,omitempty"`
	DisplayName string                `json:"displayName,omitempty"`
	Name        GooglePlusPersonName  `json:"name,omitempty"`
	URL         string                `json:"url,omitempty"`
	Image       GooglePlusPersonImage `json:"image,omitempty"`
	IsPlusUser  bool                  `json:"isPlusUser,omitempty"`
	Language    string                `json:"language,omitempty"`
	Verified    bool                  `json:"verified,omitempty"`
}

type GooglePlusPersonName struct {
	FamilyName string `json:"familyName,omitempty"`
	GivenName  string `json:"givenName,omitempty"`
}

type GooglePlusPersonImage struct {
	URL       string `json:"url,omitempty"`
	IsDefault bool   `json:"isDefault,omitempty"`
}

func (apiutil *ClientUtil) GetSCIMUser() (scim.User, error) {
	user := scim.User{}

	// Get Email
	googleUserinfoEmail, err := apiutil.GetUserinfoEmail()
	if err != nil {
		return user, err
	}

	err = user.AddEmail(googleUserinfoEmail.Email, true)
	if err != nil {
		return user, err
	}

	// Get Real Name
	googleUserinfo, err := apiutil.GetUserinfo()
	if err != nil {
		return user, err
	}
	user.Name = scim.Name{
		GivenName:  strings.TrimSpace(googleUserinfo.GivenName),
		FamilyName: strings.TrimSpace(googleUserinfo.FamilyName),
		Formatted:  strings.TrimSpace(googleUserinfo.Name)}

	return user, nil
}
