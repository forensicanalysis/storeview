// Copyright (c) 2020 Siemens AG
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// Author(s): Jonas Plum

package main

import (
	"io"
	"net/http"

	"github.com/spf13/pflag"

	"github.com/forensicanalysis/forensicworkflows/cmd"
	"github.com/forensicanalysis/storeview/cobraserver"
)

func listTasks() *cobraserver.Command {
	return &cobraserver.Command{
		Name:   "listTasks",
		Route:  "/tasks",
		Method: http.MethodGet,
		SetupFlags: func(f *pflag.FlagSet) {
			// f.String("directory", "/", "current directory")
			// f.String("type", "file", "item type")
		},
		Handler: func(w io.Writer, _ io.Reader, flags *pflag.FlagSet) error {

			runCmd := cmd.Run()
			commands := runCmd.Commands()
			var children []string
			for _, command := range commands {
				children = append(children, command.Name())
			}

			return cobraserver.PrintAny(w, children)
		},
	}
}
