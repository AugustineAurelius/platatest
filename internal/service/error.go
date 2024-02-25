package service

import "errors"

var (
	IdErr                   = errors.New("no such id")
	CriteriaErr             = errors.New("no search criteria found, you can either search by id via /id/4 to search by code via /currency/EUR/USD}")
	UnmarshalIdErr          = errors.New("bad ID")
	NotSupportedCurrencyErr = errors.New("doesn't support such currency code")
	CreateErr               = errors.New("cannot decode currency data")
)
