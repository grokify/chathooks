package config

import (
	"errors"
	"net/url"
)

const (
	DefaultIconFile = "icon_webhookrc_512x512.png"
)

var Icons = map[string]string{
	"circleci":   "icon_circleci_128x128.png",
	"codeship":   "icon_codeship_512x512.png",
	"deskdotcom": "icon_deskdotcom_400x400.png",
	"enchant":    "icon_enchant_400x400.png",
	"gosquared":  "icon_gosquared_128x128.png",
	"heroku":     "icon_heroku_512x512.png",
	"magnumci":   "icon_magnumci_400x400.png",
	"opsgenie":   "icon_opsgenie_128x128.png",
	"papertrail": "icon_papertrail_128x128.png",
	"pingdom":    "icon_pingdom_512x512.png",
	"statuspage": "icon_statuspage_512x512.png",
	"userlike":   "icon_userlike_512x512.png",
	"victorops":  "icon_victorops_225x225.png"}

func joinURL(baseURL string, pathPart string) (*url.URL, error) {
	u, err := url.Parse(pathPart)
	if err != nil {
		return &url.URL{}, nil
	}
	base, err := url.Parse("http://grokify.github.io/webhookproxy/images/icons/")
	if err != nil {
		return &url.URL{}, nil
	}
	return base.ResolveReference(u), nil
	//fmt.Println(base.ResolveReference(u).String())
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
