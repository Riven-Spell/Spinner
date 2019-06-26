package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/virepri/Spinner/common"
)

type SpinnerFlags struct {
	OAuthToken  string
	LogLevel    string
	LogLocation string
	LogfileOnly bool
}

func (s SpinnerFlags) Cook(args []string, lcm common.LifecycleManager) error {
	// Token must be explicit.
	if s.OAuthToken == "" {
		return errors.New("oauth token must be supplied")
	}

	if s.LogLocation != "" {
		fi, err := os.Stat(s.LogLocation)

		if err == nil && fi.IsDir() {
			fileName := time.Now().Format(time.RFC1123) + ".txt"
			fileName = filepath.Join(s.LogLocation, fileName)

			if f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666); err == nil {
				if !s.LogfileOnly {
					lcm.SetLogDestination(&common.RedirectWriter{
						UnderlyingWriters: []io.Writer{
							os.Stdout,
							f,
						},
					})
				} else {
					lcm.SetLogDestination(f)
				}
			} else {
				return errors.New(fmt.Sprintf("could not create file %s: %s", fileName, err))
			}
		} else {
			if f, err := os.OpenFile(s.LogLocation, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666); err == nil {
				if !s.LogfileOnly {
					lcm.SetLogDestination(&common.RedirectWriter{
						UnderlyingWriters: []io.Writer{
							os.Stdout,
							f,
						},
					})
				} else {
					lcm.SetLogDestination(f)
				}
			} else {
				return errors.New(fmt.Sprintf("could not create file %s: %s", s.LogLocation, err))
			}
		}
	} else if s.LogfileOnly {
		return errors.New(fmt.Sprintf("--logfile-only requires a specified logfile (--log-file)"))
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
