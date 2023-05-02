package cli

import "net/url"

func GetNewsHandler(u *url.URL) error {
	return GetNews()
}

func GetNews() error {
	return nil
}
