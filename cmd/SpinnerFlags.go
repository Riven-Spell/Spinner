package cmd

import (
	"errors"
)

type SpinnerFlags struct {
	OAuthToken string
}

func (s SpinnerFlags) Verify(args []string) error {
	if s.OAuthToken == "" && len(args) == 0 {
		return errors.New("oauth token must be supplied")
	}

	return nil
}

var sf SpinnerFlags