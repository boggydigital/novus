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

	if err := GetContent(); err != nil {
		return err
	}

	if err := Decode(); err != nil {
		return err
	}

	if err := MatchContent(syncStart); err != nil {
		return err
	}

	if err := ReduceContent(syncStart); err != nil {
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

	if err := ResetErrors(); err != nil {
		return err
	}

	return nil
}
