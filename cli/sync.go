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

	// debug
	//if err := PrintAll(); err != nil {
	//	return err
	//}

	if err := Backup(); err != nil {
		return err
	}

	return nil
}
