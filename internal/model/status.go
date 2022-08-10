package model

import (
	"encoding/json"
	"errors"
)

type Status struct {
	string
}

var (
	ErrPlatformInvalidParam = errors.New("unknown value status")

	StatusUNKNOWN    = Status{"UNKNOWN"}
	StatusNEW        = Status{"NEW"}
	StatusPROCESSING = Status{"PROCESSING"}
	StatusINVALID    = Status{"INVALID"}
	StatusPROCESSED  = Status{"PROCESSED"}
)

func (s Status) String() string {
	return s.string
}

func GetStatus(s string) (Status, error) {
	switch s {
	case StatusNEW.string, StatusPROCESSING.string, StatusINVALID.string, StatusPROCESSED.string:
		return Status{s}, nil
	default:
		return StatusUNKNOWN, ErrPlatformInvalidParam
	}
}

func (s Status) MarshalJSON() ([]byte, error) {
	json, err := json.Marshal(s.String())
	return json, err
}
