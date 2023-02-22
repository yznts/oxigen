package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"strings"

	"github.com/kyoto-framework/zen/v3/slice"
)

var FuncMap = template.FuncMap{
	"title": strings.Title,
	"fonts": func() []string {
		files, err := os.ReadDir("dist/fonts")
		if err != nil {
			panic(fmt.Errorf("fonts directory not found: %v", err))
		}
		return slice.Filter(slice.Map(files, func(e fs.DirEntry) string {
			return e.Name()
		}), func(name string) bool {
			return strings.HasSuffix(name, ".ttf")
		})
	},
}
