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
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/markbates/pkger"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/forensicanalysis/forensicstore"
	"github.com/forensicanalysis/storeview/cobraserver"
)

//go:generate yarn install
//go:generate yarn build
//go:generate go get -u github.com/markbates/pkger/cmd/pkger
//go:generate pkger -o assets

func main() {
	listTables := &cobraserver.Command{
		Name:   "listTables",
		Route:  "/tables",
		Method: http.MethodGet,
		Handler: func(w io.Writer, _ io.Reader, flags *pflag.FlagSet) error {
			storeName := flags.Args()[0]
			store, teardown, err := forensicstore.Open(storeName)
			if err != nil {
				return err
			}
			defer teardown()

			conn := store.Connection()
			stmt := conn.Prep(
				"SELECT " +
					"json_extract(json, '$.type') as type, " +
					"count(json) as count " +
					"FROM elements " +
					"GROUP BY json_extract(json, '$.type')",
			)
			var filtered []forensicstore.Element
			for {
				if hasRow, err := stmt.Step(); err != nil {
					return err
				} else if !hasRow {
					break
				}
				filtered = append(filtered, forensicstore.Element{
					"name":  stmt.GetText("type"),
					"count": stmt.GetInt64("count"),
				})
			}
			err = stmt.Finalize()
			if err != nil {
				return err
			}

			return cobraserver.PrintAny(w, filtered)
		},
	}

	selectItems := &cobraserver.Command{
		Name:   "selectItems",
		Route:  "/items",
		Method: http.MethodGet,
		SetupFlags: func(f *pflag.FlagSet) {
			f.String("type", "", "item type")
			f.String("uid", "", "uid")
			f.StringArray("filter", nil, "")
			f.StringArray("sort", nil, "")
			f.Int("offset", 0, "")
			f.Int("limit", 30, "")
		},
		Handler: func(w io.Writer, _ io.Reader, flags *pflag.FlagSet) error {
			uid, err := flags.GetString("uid")
			if err != nil {
				return err
			}

			storeName := flags.Args()[0]
			store, teardown, err := forensicstore.Open(storeName)
			if err != nil {
				return err
			}
			defer teardown()

			// get single item
			if uid != "" {
				item, err := store.Get(uid)
				if err != nil {
					return err
				}
				return cobraserver.PrintJSON(w, item)
			}

			name, err := flags.GetString("type")
			if err != nil {
				return err
			}

			if name == "" {
				return errors.New("type or uid must be set")
			}

			opt := NewSelectOptions()

			sort, err := flags.GetStringArray("sort")
			if err != nil {
				return err
			}
			if len(sort) > 0 {
				for _, sorting := range sort {
					parts := strings.SplitN(sorting, ":", 2)
					if len(parts) != 2 || parts[0] == "" {
						return fmt.Errorf("sort parameter %s invalid", sorting)
					}
					switch parts[1] {
					case "ASC":
						opt.Sort[parts[0]] = SortAscending
					case "DESC":
						opt.Sort[parts[0]] = SortDescending
					case "":
						opt.Sort[parts[0]] = SortDefault
					default:
						return fmt.Errorf("sort direction %s invalid", sorting)
					}
				}
			}
			filter, err := flags.GetStringArray("filter")
			if err != nil {
				return err
			}
			if len(filter) > 0 {
				for _, filtering := range filter {
					parts := strings.SplitN(filtering, ":", 2)
					if len(parts) != 2 || parts[0] == "" {
						return fmt.Errorf("filte parameter %s invalid", filtering)
					}
					opt.Filter[parts[0]] = parts[1]
				}
			}

			offset, err := flags.GetInt("offset")
			if err != nil {
				return err
			}
			limit, err := flags.GetInt("limit")
			if err != nil {
				return err
			}

			opt.Limit = limit
			opt.Offset = offset

			items, err := QueryStore(store, name, opt)
			if err != nil {
				return err
			}
			return cobraserver.PrintJSONList(w, items)
		},
	}

	loadFile := &cobraserver.Command{
		Name:   "loadFile",
		Route:  "/file",
		Method: http.MethodGet,
		SetupFlags: func(f *pflag.FlagSet) {
			f.String("path", "", "path")
		},
		Handler: func(w io.Writer, _ io.Reader, flags *pflag.FlagSet) error {
			p, err := flags.GetString("path")
			if err != nil {
				return err
			}
			if p == "" {
				return errors.New("path must be set")
			}

			storeName := flags.Args()[0]
			store, teardown, err := forensicstore.Open(storeName)
			if err != nil {
				return err
			}
			defer teardown()

			f, err := store.LoadFile(path.Join(storeName, p)) // TODO: fix in forensicstore
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(w, f)
			return err
		},
	}

	directoryTree := &cobraserver.Command{
		Name:   "listTree",
		Route:  "/tree",
		Method: http.MethodGet,
		SetupFlags: func(f *pflag.FlagSet) {
			f.String("directory", "/", "current directory")
			f.String("type", "file", "item type")
		},
		Handler: func(w io.Writer, _ io.Reader, flags *pflag.FlagSet) error {
			storeName := flags.Args()[0]
			store, teardown, err := forensicstore.Open(storeName)
			if err != nil {
				return err
			}
			defer teardown()

			directory, err := flags.GetString("directory")
			if err != nil {
				return err
			}

			elementType, err := flags.GetString("type")
			if err != nil {
				return err
			}

			types := map[string]map[string]string{
				"file":                 {"separator": "/", "col": "json_extract(json, '$.origin.path')"},
				"directory":            {"separator": "/", "col": "json_extract(json, '$.path')"},
				"windows-registry-key": {"separator": "\\", "col": "json_extract(json, '$.key')"},
			}

			col := types[elementType]["col"]
			separator := types[elementType]["separator"]
			query := fmt.Sprintf("SELECT substr(%s, length('%s')+1, instr(substr(%s, 1+length('%s')), '%s')-1) as dir "+
				"FROM 'elements' "+
				"WHERE json_extract(json, '$.type') = '%s' "+
				"AND %s LIKE '%s%%' "+
				"GROUP BY dir",
				col, directory, col, directory, separator, elementType, col, directory)

			fmt.Println(query)

			var children []string

			conn := store.Connection()

			stmt := conn.Prep(query)
			for {
				if hasRow, err := stmt.Step(); err != nil {
					return err
				} else if !hasRow {
					break
				}
				children = append(children, stmt.GetText("dir"))
			}
			err = stmt.Finalize()
			if err != nil {
				return err
			}

			return cobraserver.PrintAny(w, children)
		},
	}

	var staticPath pkger.Dir = "/dist"
	rootCmd := cobraserver.Application("fstore", 800, 600, staticPath, false, listTables, selectItems, loadFile, directoryTree)

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

type Direction string

const (
	SortDefault    Direction = ""
	SortAscending            = "ASC"
	SortDescending           = "DESC"
)

type SelectOptions struct {
	Sort   map[string]Direction
	Filter map[string]string
	Limit  int
	Offset int
}

func NewSelectOptions() *SelectOptions {
	return &SelectOptions{
		Sort:   map[string]Direction{},
		Filter: map[string]string{},
		Limit:  30,
		Offset: 0,
	}
}

func QueryStore(store *forensicstore.ForensicStore, itemType string, options *SelectOptions) ([]forensicstore.JSONElement, error) {
	q := "SELECT json FROM elements"

	filters := []string{
		fmt.Sprintf("json_extract(json, '$.type') = '%s'", itemType),
	}

	if len(options.Filter) > 0 {
		for column, filtering := range options.Filter {
			if filtering != "" {
				filters = append(filters, fmt.Sprintf("json_extract(json, '$.%s') LIKE '%%%s%%'", column, filtering))
			}
		}
	}

	q += " WHERE " + strings.Join(filters, " AND ")

	if len(options.Sort) > 0 {
		var sorts []string
		for column, sorting := range options.Sort {
			if sorting != "" {
				sorts = append(sorts, fmt.Sprintf("json_extract(json, '$.%s') %s", column, sorting))
			}
		}
		if len(sorts) > 0 {
			q += " ORDER BY " + strings.Join(sorts, ", ")
		}
	}

	q += fmt.Sprintf(" LIMIT %d", options.Limit)
	q += fmt.Sprintf(" OFFSET %d", options.Offset)
	q += ";"
	fmt.Println(q)
	return store.Query(q)
}
