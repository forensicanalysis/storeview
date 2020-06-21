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
	"errors"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"strconv"

	"github.com/spf13/pflag"

	forensicworkflows "github.com/forensicanalysis/forensicworkflows/cmd"
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
			runCmd := forensicworkflows.Run()
			commands := runCmd.Commands()
			var children []string
			for _, command := range commands {
				children = append(children, command.Name())
			}

			return cobraserver.PrintAny(w, children)
		},
	}
}

func getTaskSchema() *cobraserver.Command {
	return &cobraserver.Command{
		Name:   "task",
		Route:  "/task",
		Method: http.MethodGet,
		SetupFlags: func(f *pflag.FlagSet) {
			f.String("name", "", "command name")
		},
		Handler: func(w io.Writer, _ io.Reader, flags *pflag.FlagSet) error {
			name, err := flags.GetString("name")
			if err != nil {
				return err
			}

			runCmd := forensicworkflows.Run()
			commands := runCmd.Commands()
			for _, command := range commands {
				if command.Name() == name {
					schema := flagsToSchema(command.Flags())
					return cobraserver.PrintAny(w, schema)
				}
			}

			return errors.New("command not found")
		},
	}
}

func flagsToSchema(flags *pflag.FlagSet) forensicworkflows.JSONSchema {
	schema := forensicworkflows.JSONSchema{
		Properties: map[string]forensicworkflows.Property{},
		Required:   []string{},
	}

	flags.VisitAll(func(f *pflag.Flag) {
		typeMapping := map[string]string{
			"string": "string",
			"int":    "integer",
			"bool":   "boolean",
			"float":  "number",
		}

		property := forensicworkflows.Property{
			Type:        typeMapping[f.Value.Type()],
			Description: f.Usage,
		}

		if f.DefValue != "" {
			var defaultValue interface{}
			var err error
			switch f.Value.Type() {
			case "string":
				property.Default = f.DefValue
			case "int":
				defaultValue, err = strconv.ParseInt(f.DefValue, 10, 64)
			case "bool":
				defaultValue, err = strconv.ParseBool(f.DefValue)
			case "float":
				defaultValue, err = strconv.ParseFloat(f.DefValue, 64)
			}
			if err == nil {
				property.Default = defaultValue
			}
		}
		schema.Properties[f.Name] = property
		if _, ok := f.Annotations[cobra.BashCompOneRequiredFlag]; ok {
			schema.Required = append(schema.Required, f.Name)
		}
	})
	return schema
}
