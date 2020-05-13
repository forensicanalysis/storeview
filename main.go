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

	"github.com/forensicanalysis/forensicstore/goforensicstore"
	"github.com/forensicanalysis/forensicstore/gojsonlite"
	"github.com/forensicanalysis/forensicstore/gostore"
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
			store, err := goforensicstore.NewJSONLite(storeName)
			if err != nil {
				return err
			}

			items, err := store.Store.(*gojsonlite.JSONLite).Query("SELECT name FROM sqlite_master")
			if err != nil {
				return err
			}

			var filtered []gostore.Item
			for _, item := range items {
				if name, ok := item["name"]; ok {
					if sname, ok := name.(string); ok {
						if !strings.HasPrefix(sname, "_") && !strings.HasPrefix(sname, "sqlite") {

							items, err := store.Store.(*gojsonlite.JSONLite).Query("SELECT count(uid) as count FROM `" + sname + "`")
							if err != nil {
								return errors.Wrap(err, sname)
							}

							if len(items) > 0 {
								if count, ok := items[0]["count"]; ok {
									if cint, ok := count.(float64); ok {
										item["count"] = int(cint)
									}
								}
							}

							filtered = append(filtered, item)
						}
					}
				}
			}

			return cobraserver.PrintJSON(w, filtered)
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
			store, err := goforensicstore.NewJSONLite(storeName)
			if err != nil {
				return err
			}

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

			items, err := QueryStore(store.Store.(*gojsonlite.JSONLite), name, opt)
			if err != nil {
				return err
			}
			return cobraserver.PrintJSON(w, items)
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
			store, err := goforensicstore.NewJSONLite(storeName)
			if err != nil {
				return err
			}

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
		Name:   "listTables",
		Route:  "/tree",
		Method: http.MethodGet,
		SetupFlags: func(f *pflag.FlagSet) {
			f.String("directory", "/", "current directory")
			f.String("type", "file", "item type")
		},
		Handler: func(w io.Writer, _ io.Reader, flags *pflag.FlagSet) error {
			storeName := flags.Args()[0]
			store, err := goforensicstore.NewJSONLite(storeName)
			if err != nil {
				return err
			}

			directory, err := flags.GetString("directory")
			if err != nil {
				return err
			}

			elementType, err := flags.GetString("type")
			if err != nil {
				return err
			}

			types := map[string]map[string]string{
				"file": {"separator": "/", "col": "origin.path"},
				"directory": {"separator": "/", "col": "path"},
				"windows-registry-key": {"separator": "\\", "col": "key"},
			}

			col := types[elementType]["col"]
			separator := types[elementType]["separator"]
			query := fmt.Sprintf("SELECT substr(`%s`, length(\"%s\")+1, instr(substr(`%s`, 1+length(\"%s\")), \"%s\")-1) as dir "+
				"FROM '%s' "+
				"WHERE `%s` LIKE \"%s%%\" "+
				"GROUP BY dir;",
				col, directory, col, directory, separator, elementType, col, directory)

			println(query)

			directories, err := store.Store.(*gojsonlite.JSONLite).Query(query)
			if err != nil {
				return err
			}

			var children []string
			for _, dir := range directories {
				children = append(children, dir["dir"].(string))
			}

			return cobraserver.PrintJSON(w, children)
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

func QueryStore(store *gojsonlite.JSONLite, itemType string, options *SelectOptions) ([]gostore.Item, error) {
	q := fmt.Sprintf("SELECT * FROM \"%s\"", itemType)

	if len(options.Filter) > 0 {
		var filters []string
		for column, filtering := range options.Filter {
			if filtering != "" {
				filters = append(filters, fmt.Sprintf("%s LIKE '%%%s%%'", column, filtering))
			}
		}
		if len(filters) > 0 {
			q += " WHERE " + strings.Join(filters, " AND ")
		}
	}

	if len(options.Sort) > 0 {
		var sorts []string
		for column, sorting := range options.Sort {
			if sorting != "" {
				sorts = append(sorts, fmt.Sprintf("%s %s", column, sorting))
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
