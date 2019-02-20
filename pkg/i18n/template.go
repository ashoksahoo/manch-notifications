package i18n

import (
	"bytes"
)
import "text/template"

type DataModel struct {
	Name         string `json:"name" bson:"name"`
	Name2        string `json:"name2" bson:"name2"`
	Name3        string `json:"name3" bson:"name3"`
	Count        int    `json:"count" bson:"count"`
	Post         string `json:"post" bson:"post"`
	Comment      string `json:"comment" bson:"comment"`
	DeleteReason string `json:"deletereason" bson:"deletereason"`
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

func GetHtmlString(lang string, s string, d DataModel) string {
	var output bytes.Buffer
	var tpl string
	if HtmlStrings[lang] == nil {
		lang = "en"
	}
	tpl = HtmlStrings[lang][s]
	if tpl == "" {
		tpl = HtmlStrings["en"][s]
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
