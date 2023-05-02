package cli

import "net/url"

func SyncHandler(u *url.URL) error {
	return Sync()
}

func Sync() error {
	return nil
}
