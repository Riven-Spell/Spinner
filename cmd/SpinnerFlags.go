package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/virepri/Spinner/common"
)

type SpinnerFlags struct {
	OAuthToken  string
	LogLevel    string
	LogLocation string
}

func (s SpinnerFlags) Cook(args []string, lcm common.LifecycleManager) error {
	// Token must be explicit.
	if s.OAuthToken == "" {
		return errors.New("oauth token must be supplied")
	}

	if s.LogLocation != "" {
		fi, err := os.Stat(s.LogLocation)

		if err != nil {
			return errors.New(fmt.Sprintf("invalid log location specified: %s", err))
		}

		if fi.IsDir() {
			// TODO: Open a file under this directory.
		} else {
			// TODO: Open and append to the file itself.
		}
	}

	if v, exists := common.MLogLevel[strings.ToLower(s.LogLevel)]; s.LogLevel != "" {
		if !exists {
			return errors.New("invalid log level. valid log levels are as follows: info, error, warning, fatal")
		}

		lcm.SetMinimumLogLevel(v)
	}

	return nil
}

var sf SpinnerFlags
