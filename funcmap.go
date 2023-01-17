package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"strings"

	"github.com/kyoto-framework/zen/v2"
)

var FuncMap = template.FuncMap{
	"title": func(v string) string {
		return strings.Title(v)
	},
	"fonts": func() []string {
		files, err := os.ReadDir("dist/fonts")
		if err != nil {
			panic(fmt.Errorf("fonts directory not found: %v", err))
		}
		return zen.Filter(zen.Map(files, func(e fs.DirEntry) string {
			return e.Name()
		}), func(name string) bool {
			return strings.HasSuffix(name, ".ttf")
		})
	},
}
