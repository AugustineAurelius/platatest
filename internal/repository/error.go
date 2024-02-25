package repository

import "errors"

var (
	TxErr        = errors.New("error while trying to begin transaction")
	TxCommitErr  = errors.New("error while commit creansaction")
	CreateErr    = errors.New("error while creating currency")
	GetByIdErr   = errors.New("error while get by id")
	GetByNameErr = errors.New("error while get by name")
	UpdateErr    = errors.New("error while upgrate price of currency")
	FetchErr     = errors.New("error while fetch all currencies")
	ScanErr      = errors.New("error while scan temp currency")
)
