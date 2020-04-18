package cobraserver

import (
	"encoding/json"
	"io"
)

func PrintJSON(w io.Writer, i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\n"))

	return err
}
