package clientutil

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/grokify/gotilla/net/urlutil"
	"github.com/joho/godotenv"

	rc "github.com/grokify/go-ringcentral/client"
	rs "github.com/grokify/go-scim-client"
	ro "github.com/grokify/oauth2more/ringcentral"
)

func LoadEnv() error {
	envPaths := []string{}
	if len(os.Getenv("ENV_PATH")) > 0 {
		envPaths = append(envPaths, os.Getenv("ENV_PATH"))
	}
	return godotenv.Load(envPaths...)
}

func NewApiClientHttpClientBaseURL(httpClient *http.Client, serverUrl string) (*rc.APIClient, error) {
	if len(strings.TrimSpace(serverUrl)) == 0 {
		return nil, fmt.Errorf("No RingCentral API Server URL provided")
	}
	apiConfig := rc.NewConfiguration()
	apiConfig.BasePath = strings.TrimSpace(serverUrl)
	apiConfig.HTTPClient = httpClient
	apiClient := rc.NewAPIClient(apiConfig)
	return apiClient, nil
}

func NewApiClientPassword(app ro.ApplicationCredentials, pwd ro.PasswordCredentials) (*rc.APIClient, error) {
	httpClient, err := ro.NewClientPassword(app, pwd)
	if err != nil {
		return nil, err
	}
	return NewApiClientHttpClientBaseURL(httpClient, app.ServerURL)
}

func NewApiClientPasswordSimple(app ro.ApplicationCredentials, user ro.UserCredentials) (*rc.APIClient, error) {
	httpClient, err := ro.NewClientPasswordSimple(app, user)
	if err != nil {
		return nil, err
	}
	return NewApiClientHttpClientBaseURL(httpClient, app.ServerURL)
}

func NewApiClientPasswordEnv() (*rc.APIClient, error) {
	return NewApiClientPassword(
		ro.NewApplicationCredentialsEnv(),
		ro.NewPasswordCredentialsEnv())
}

func NewScimApiClient(app ro.ApplicationCredentials, pwd ro.PasswordCredentials) (*rs.APIClient, error) {
	httpClient, err := ro.NewClientPassword(app, pwd)
	if err != nil {
		return nil, err
	}
	apiConfig := rs.NewConfiguration()
	apiConfig.BasePath = urlutil.JoinAbsolute(app.ServerURL, "/scim/v2")
	apiConfig.HTTPClient = httpClient
	apiClient := rs.NewAPIClient(apiConfig)
	return apiClient, nil
}
