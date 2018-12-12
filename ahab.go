package ahab

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"

	"mvdan.cc/sh/shell"
)

// Log is a logger with support for logging levels and pretty output
var Log = logrus.New()

func init() {
	Log.SetOutput(os.Stdout)
	formatter := logrus.TextFormatter{}
	formatter.DisableTimestamp = true
	formatter.DisableLevelTruncation = true
	Log.SetFormatter(&formatter)
}

// ParseFlags parses flags
func ParseFlags(data interface{}) ([]string, error) {
	f := flags.NewParser(data, flags.HelpFlag|flags.PassDoubleDash)
	args, err := f.Parse()

	if err != nil {
		// write usage to a byte buffer and then wrap the current err with it, then when you print it out it'll be there
		if strings.Contains(err.Error(), "unknown flag") {
			f.WriteHelp(os.Stdout)
			// Log.Fatal(err)
		}
	}
	return args, err

}

// Env gets an environment variable, errors if it doesn't exist
func Env(key string) (string, error) {
	s, b := os.LookupEnv(key)
	if !b {
		err := errors.New(key)
		return "", err
	}
	return s, nil
}

// EnvOr gets an environment variable, returns defaultValue if it doesn't exist
func EnvOr(key, defaultValue string) string {
	s, b := os.LookupEnv(key)
	if !b {
		return defaultValue
	}
	return s
}

// Read reads a file and returns a string, expanding any vars, if needed
func Read(path string) (string, error) {
	p, err := shell.Fields(path, nil)
	if err != nil {
		return "", nil
	}

	expanded := strings.Join(p, "")

	buf, err := ioutil.ReadFile(expanded)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// RecoverMe should be deferred at the top of main()
func RecoverMe() {
	if r := recover(); r != nil {
		Log.Error(r)

	}
}

// FatalIf errors and exits if `err` is nil, use with `RecoverMe` for ergonomic-ish error control
func FatalIf(err error) {
	if err != nil {
		Log.Fatal(err)
	}

}
