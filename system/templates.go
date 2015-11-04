package system

//go:generate rice embed-go

import (
	"github.com/GeertJohan/go.rice"
	"github.com/denisbakhtin/blog/helpers"
	"html/template"
	"log"
	"os"
	"strings"
)

var tmpl *template.Template

func loadTemplates() {
	box := rice.MustFindBox("../views")
	tmpl = template.New("").Funcs(template.FuncMap{
		"isActive":      helpers.IsActive,
		"stringInSlice": helpers.StringInSlice,
		"dateTime":      helpers.DateTime,
		"recentPosts":   helpers.RecentPosts,
		"tags":          helpers.Tags,
		"archives":      helpers.Archives,
	})

	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".html") {
			var err error
			tmpl, err = tmpl.Parse(box.MustString(path))
			if err != nil {
				return err
			}
		}
		return nil
	}

	err := box.Walk("", fn)
	if err != nil {
		log.Panic(err)
	}
}
