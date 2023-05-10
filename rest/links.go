package rest

import (
	"net/url"
	"strings"
)

func Host(su string) string {
	if u, err := url.Parse(su); err == nil {
		return u.Scheme + "://" + u.Host
	}
	return ""
}

func AbsHref(element, host string) string {
	if host == "" {
		return element
	}
	return strings.Replace(element, "href=\"/", "href=\""+host+"/", -1)
}
