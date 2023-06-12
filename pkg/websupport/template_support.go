package websupport

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type Model struct {
	Map map[string]any
}

func ModelAndView(w http.ResponseWriter, resources *embed.FS, data Model, desired ...string) error {
	views := make([]string, 0)
	for _, v := range desired {
		views = append(views, fmt.Sprintf("resources/templates/%v.gohtml", v))
	}
	views = append(views, "resources/templates/template.gohtml")

	base := filepath.Base(views[0])
	return template.Must(template.New(base).Funcs(template.FuncMap{
		"formatFloat": func(value float64) string {
			return fmt.Sprintf("%.0f", value)
		},
		"formatPercentage": func(value float64) string {
			return fmt.Sprintf("%.0f", value*100)
		},
	}).ParseFS(resources, views...)).Execute(w, data)
}
