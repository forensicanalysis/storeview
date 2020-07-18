package main

import (
	"bytes"
	"encoding/json"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/forensicanalysis/storeview/cobraserver"
	"github.com/forensicanalysis/storeview/commands"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type resW struct {
	*bytes.Buffer
	header http.Header
}

func NewResW() *resW {
	return &resW{
		Buffer: &bytes.Buffer{},
		header: http.Header{},
	}
}

func (r resW) Header() http.Header        { return r.header }
func (r resW) WriteHeader(statusCode int) {}

type closingBuffer struct{ *bytes.Buffer }

func (c closingBuffer) Close() error { return nil }

var store = ""

func open(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	if m.Name == "open" {
		log.Println("XXXX open")
		store = strings.Trim(string(m.Payload), "\"")
		return
	}
	return
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	log.Println("using store", store)
	router := cobraserver.Router(commands.Commands, []string{store})

	u, err := url.Parse("/api" + m.Name)
	if err != nil {
		return nil, err
	}
	buffer := closingBuffer{bytes.NewBuffer(m.Payload)}
	req := &http.Request{Method: http.MethodGet, URL: u, Body: buffer}

	var match mux.RouteMatch
	if router.Match(req, &match) {
		log.Println("MATCHING ROUTE FOUND", u, match.Handler, match.Vars, store)
		res := NewResW()
		match.Handler.ServeHTTP(res, req)
		var payload interface{}
		log.Println("res", res.String())
		err = json.Unmarshal(res.Bytes(), &payload)
		return payload, err
	} else {
		log.Println("NO MATCHING ROUTE FOUND", u, match.MatchErr)
	}

	return payload, err
}

/*
// handleMessages handles messages
func handleMessages(store string) func(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	return func(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
		router := cobraserver.Router(commands.Commands, []string{store})

		u, err := url.Parse("/api" + m.Name)
		if err != nil {
			return nil, err
		}
		buffer := closingBuffer{bytes.NewBuffer(m.Payload)}
		req := &http.Request{Method: http.MethodGet, URL: u, Body: buffer}

		var match mux.RouteMatch
		if router.Match(req, &match) {
			res := NewResW()
			match.Handler.ServeHTTP(res, req)
			json.Unmarshal(res.Bytes(), &payload)
		} else {
			log.Println("NO MATCHING ROUTE FOUND", u, match.MatchErr)
		}

		return payload, nil
	}
}
*/
