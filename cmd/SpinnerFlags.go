package cmd

import (
	"errors"
)

type SpinnerFlags struct {
	OAuthToken string
}

func (s SpinnerFlags) Verify(args []string) error {
	// Token must be explicit.
	if s.OAuthToken == "" {
		return errors.New("oauth token must be supplied")
	}

	return nil
}

var sf SpinnerFlags