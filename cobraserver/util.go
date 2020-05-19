package cobraserver

import (
	"encoding/json"
	"github.com/forensicanalysis/forensicstore"
	"io"
)

func PrintAny(w io.Writer, i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}

	return PrintJSON(w, b)
}

func PrintJSONList(w io.Writer, elements []forensicstore.JSONElement) error {
	_, err := w.Write([]byte("["))
	if err != nil {
		return err
	}

	for i, element := range elements {
		if i != 0 {
			_, err := w.Write([]byte(","))
			if err != nil {
				return err
			}
		}
		_, err := w.Write(element)
		if err != nil {
			return err
		}
	}

	_, err = w.Write([]byte("]\n"))
	return err
}

func PrintJSON(w io.Writer, b []byte) error {
	_, err := w.Write(b)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\n"))

	return err
}
