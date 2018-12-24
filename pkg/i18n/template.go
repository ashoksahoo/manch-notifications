package i18n

import (
	"bytes"
)
import "text/template"

type DataModel struct {
	Name    string
	Name2   string
	Name3   string
	Count   int
	Post    string
	Comment string
	DeleteReason string
}

func GetString(lang string, s string, d DataModel) string {
	var output bytes.Buffer
	var tpl string
	if Strings[lang] == nil {
		lang = "en"
	}
	tpl = Strings[lang][s]
	if tpl == "" {
		tpl = Strings["en"][s]
	}
	if tpl == "" {
		return ""
	} else {
		tmpl, err := template.New(lang + s).Parse(tpl)
		if err != nil {
			return ""
		}
		err = tmpl.Execute(&output, d)
		if err != nil {
			return ""
		}
		return output.String()
	}

}
