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

package commands

import (
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/pflag"

	"github.com/forensicanalysis/forensicstore"
	"github.com/forensicanalysis/storeview/cobraserver"
)

func ListTree() *cobraserver.Command {
	return &cobraserver.Command{
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

			stmt, err := conn.Prepare(query)
			if err != nil {
				return err
			}
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
}

func SelectItems() *cobraserver.Command {
	return &cobraserver.Command{
		Name:   "selectItems",
		Route:  "/items",
		Method: http.MethodGet,
		SetupFlags: func(f *pflag.FlagSet) {
			f.String("type", "", "item type")
			f.String("uid", "", "uid")
			f.StringArray("filter", nil, "")
			f.StringArray("sort", nil, "")
			f.StringArray("labels", nil, "")
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

			/*
				if name == "" {
					return errors.New("type or uid must be set")
				}
			*/

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
			labels, err := flags.GetStringArray("labels")
			if err != nil {
				return err
			}
			opt.Labels = labels

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

			count, items, err := queryStore(store, name, opt)
			if err != nil {
				return err
			}
			return cobraserver.PrintJSONList(w, count, items)
		},
	}
}

func ListTables() *cobraserver.Command {
	return &cobraserver.Command{
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

			q := "SELECT " +
				"json_extract(json, '$.type') as type, " +
				"count(json) as count " +
				"FROM elements " +
				"GROUP BY json_extract(json, '$.type')"
			fmt.Println(q)
			stmt := conn.Prep(q)
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
}

func Label() *cobraserver.Command {
	return &cobraserver.Command{
		Name:   "label",
		Route:  "/label",
		Method: http.MethodGet,
		SetupFlags: func(f *pflag.FlagSet) {
			f.String("id", "", "id")
			f.StringArray("label", nil, "label")
		},
		Handler: func(w io.Writer, _ io.Reader, flags *pflag.FlagSet) error {
			storeName := flags.Args()[0]
			store, teardown, err := forensicstore.Open(storeName)
			if err != nil {
				return err
			}
			defer teardown()

			id, err := flags.GetString("id")
			if err != nil {
				return err
			}

			labels, err := flags.GetStringArray("label")
			if err != nil {
				return err
			}

			conn := store.Connection()

			b, err := json.Marshal(labels)
			if err != nil {
				return err
			}

			stmt, err := conn.Prepare(fmt.Sprintf(
				"UPDATE elements "+
					"SET json = json_patch(json,'{\"labels\": %s}') "+
					"WHERE id = $id", string(b),
			))
			if err != nil {
				return err
			}

			stmt.SetText("$id", id)

			_, err = stmt.Step()
			if err != nil {
				return err
			}

			err = stmt.Finalize()
			if err != nil {
				return err
			}

			return cobraserver.PrintAny(w, true)
		},
	}
}

func Labels() *cobraserver.Command {
	return &cobraserver.Command{
		Name:   "labels",
		Route:  "/labels",
		Method: http.MethodGet,
		Handler: func(w io.Writer, _ io.Reader, flags *pflag.FlagSet) error {
			storeName := flags.Args()[0]
			store, teardown, err := forensicstore.Open(storeName)
			if err != nil {
				return err
			}
			defer teardown()

			conn := store.Connection()

			q := "SELECT json_extract(json, '$.labels') AS labels " +
				"FROM elements " +
				"WHERE json_extract(json, '$.labels') != ''"
			fmt.Println(q)

			stmt, err := conn.Prepare(q)
			if err != nil {
				return err
			}

			seen := map[string]bool{}
			labels := []string{}

			for {
				if hasRow, err := stmt.Step(); err != nil {
					return err
				} else if !hasRow {
					break
				}

				var elabels []string
				label := stmt.GetText("labels")
				err = json.Unmarshal([]byte(label), &elabels)
				if err != nil {
					return err
				}

				for _, elabel := range elabels {
					if _, ok := seen[elabel]; !ok {
						labels = append(labels, elabel)
						seen[elabel] = true
					}
				}
			}
			err = stmt.Finalize()
			if err != nil {
				return err
			}

			return cobraserver.PrintAny(w, labels)
		},
	}
}

func Query() *cobraserver.Command {
	return &cobraserver.Command{
		Name:   "query",
		Route:  "/query",
		Method: http.MethodGet,
		SetupFlags: func(f *pflag.FlagSet) {
			f.String("query", "", "query")
		},
		Handler: func(w io.Writer, _ io.Reader, flags *pflag.FlagSet) error {
			storeName := flags.Args()[0]
			store, teardown, err := forensicstore.Open(storeName)
			if err != nil {
				return err
			}
			defer teardown()

			query, err := flags.GetString("query")
			if err != nil {
				return err
			}

			q, err := expandQuery("SELECT json FROM elements WHERE " + query)
			if err != nil {
				return err
			}
			elements, err := storequery(store, q)
			if err != nil {
				return err
			}

			return cobraserver.PrintJSONList(w, int64(len(elements)), elements)
		},
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
	Labels []string
	Limit  int
	Offset int
}

func NewSelectOptions() *SelectOptions {
	return &SelectOptions{
		Sort:   map[string]Direction{},
		Filter: map[string]string{},
		Labels: []string{},
		Limit:  30,
		Offset: 0,
	}
}

func queryStore(store *forensicstore.ForensicStore, itemType string, options *SelectOptions) (int64, []forensicstore.JSONElement, error) {
	q := ""

	var filters []string

	if itemType != "" {
		filters = append(filters, fmt.Sprintf("json_extract(json, '$.type') = '%s'", itemType))
	}

	if len(options.Filter) > 0 {
		for column, filtering := range options.Filter {
			if filtering != "" {
				if column == "elements" {
					filters = append(filters, fmt.Sprintf("%s MATCH '%s'", column, filtering))
				} else {
					filters = append(filters, fmt.Sprintf("json_extract(json, '$.%s') LIKE '%%%s%%'", column, filtering))
				}
			}
		}
	}

	if len(options.Labels) > 0 {
		for _, label := range options.Labels {
			filters = append(filters, fmt.Sprintf("json_extract(json, '$.labels.%s')", label))
		}
	}

	if len(filters) > 0 {
		q += " WHERE " + strings.Join(filters, " AND ")
	}

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

	countQuery := "SELECT count(json) as count FROM elements" + q

	var count int64
	countCached, found := queryCache.Get(countQuery)
	if found {
		count = countCached.(int64)
	} else {
		conn := store.Connection()

		fmt.Println(countQuery)
		stmt, err := conn.Prepare(countQuery)
		if err != nil {
			return 0, nil, err
		}

		_, err = stmt.Step()
		if err != nil {
			return 0, nil, err
		}

		count = stmt.GetInt64("count")
		queryCache.Set(countQuery, count, cache.DefaultExpiration)

		err = stmt.Finalize()
		if err != nil {
			return 0, nil, err
		}
	}

	q += fmt.Sprintf(" LIMIT %d", options.Limit)
	q += fmt.Sprintf(" OFFSET %d", options.Offset)
	q += ";"
	elements, err := storequery(store, "SELECT json FROM elements"+q)
	return count, elements, err
}
