package cli

import (
	"net/url"
	"time"
)

func SyncHandler(u *url.URL) error {
	novusUrl := u.Query().Get("novus-url")
	return Sync(novusUrl)
}

func Sync(novusUrl string) error {

	syncStart := time.Now().Unix()

	// resetting errors before sync, to only track the latest iteration errors
	if err := ResetErrors(); err != nil {
		return err
	}

	if err := GetContent(); err != nil {
		return err
	}

	if err := Decode(); err != nil {
		return err
	}

	if err := MatchContent(); err != nil {
		return err
	}

	if err := Reduce(syncStart); err != nil {
		return err
	}

	if err := Diff(); err != nil {
		return err
	}

	if err := PublishAtom(novusUrl); err != nil {
		return err
	}

	if err := ResetChanges(); err != nil {
		return nil
	}

	if err := Backup(); err != nil {
		return err
	}

	return nil
}
