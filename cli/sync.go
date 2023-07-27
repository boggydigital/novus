package cli

import (
	"net/url"
	"time"
)

func SyncHandler(_ *url.URL) error {
	return Sync()
}

func Sync() error {

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

	if err := PublishAtom(); err != nil {
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
