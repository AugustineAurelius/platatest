package config

import "errors"

var (
	ReadErr      = errors.New("error while try to read configuration file")
	UnmarshalErr = errors.New("error while try to marshal configuration file")
)
