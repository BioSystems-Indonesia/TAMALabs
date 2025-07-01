package common

import "github.com/kardianos/hl7"

func Encode(msg any) (string, error) {
	e := hl7.NewEncoder(nil)
	encoded, err := e.Encode(msg)
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}
