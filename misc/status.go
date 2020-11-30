package misc

import (
	"bytes"
	"encoding/json"
)

type Status struct {
	Version     StatusVersion     `json:"version"`
	Players     StatusPlayers     `json:"players"`
	Description StatusDescription `json:"description"` //TODO chat component, favicon
}

type StatusVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type StatusPlayers struct {
	Max     int            `json:"max"`
	Online  int            `json:"online"`
	Players []StatusPlayer `json:"sample"`
}

type StatusPlayer struct {
	Name string `json:"name"`
	Uuid string `json:"id"`
}

type StatusDescription struct {
	Text string `json:"text"`
}

func (s *Status) ToJson() (string, error) {
	b, err := json.Marshal(&s)
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	if err := json.Compact(&buff, b); err != nil {
		return "", err
	}
	return string(buff.Bytes()), nil
}
