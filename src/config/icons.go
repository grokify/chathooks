package config

import (
	"errors"
	"net/url"
)

const (
	DefaultIconFile = "icon_webhookrc_512x512.png"
)

var Icons = map[string]string{
	"aha":        "icon_aha_256x256.png",
	"appsignal":  "icon_appsignal_400x400.png",
	"apteligent": "icon_apteligent_496x496.png",
	"circleci":   "icon_circleci_128x128.png",
	"codeship":   "icon_codeship_512x512.png",
	"confluence": "icon_confluence_256x256.png",
	"datadog":    "icon_datadog_512x512.png",
	"deskdotcom": "icon_deskdotcom_400x400.png",
	"enchant":    "icon_enchant_400x400.png",
	"gosquared":  "icon_gosquared_128x128.png",
	"heroku":     "icon_heroku_512x512.png",
	"librato":    "icon_librato_128x128.png",
	"magnumci":   "icon_magnumci_400x400.png",
	"marketo":    "icon_marketo_250x250.png",
	"opsgenie":   "icon_opsgenie_128x128.png",
	"papertrail": "icon_papertrail_128x128.png",
	"pingdom":    "icon_pingdom_512x512.png",
	"raygun":     "icon_raygun_512x512.png",
	"runscope":   "icon_runscope_400x400.png",
	"semaphore":  "icon_semaphore_512x512.png",
	"statuspage": "icon_statuspage_512x512.png",
	"travisci":   "icon_travisci_225x225.png",
	"userlike":   "icon_userlike_512x512.png",
	"victorops":  "icon_victorops_225x225.png"}

func joinURL(baseURL string, pathPart string) (*url.URL, error) {
	u, err := url.Parse(pathPart)
	if err != nil {
		return &url.URL{}, nil
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return &url.URL{}, nil
	}
	return base.ResolveReference(u), nil
}

func getAppIconFile(appSlug string) (string, error) {
	if file, ok := Icons[appSlug]; ok {
		return file, nil
	}
	return "", errors.New("E_NO_APP_ICON_FILE")
}

func getAppIconFileWithDefault(appSlug string) string {
	file, err := getAppIconFile(appSlug)
	if err != nil {
		file = DefaultIconFile
	}
	return file
}

func buildIconURL(baseURL string, appSlug string) (*url.URL, error) {
	iconFile := getAppIconFileWithDefault(appSlug)
	return joinURL(baseURL, iconFile)
}
