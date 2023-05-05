package cli

import "net/url"

func SyncHandler(_ *url.URL) error {
	return Sync()
}

func Sync() error {

	if err := GetNews(); err != nil {
		return err
	}

	if err := Reduce(); err != nil {
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

	return nil
}
