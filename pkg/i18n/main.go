package i18n

import (
	"bytes"
)
import "text/template"

var Strings = map[string]map[string]string{
	"en": {
		"comment_multi": "{{.Name}} & {{.Count}} others commented on Your Post {{.Post}}",
		"comment_three": "{{.Name}}, {{.Name2}} & {{.Name3}} commented on Your Post {{.Post}}",
		"comment_two":   "{{.Name}} & {{.Name2}} commented on Your Post {{.Post}}",
		"comment_one":   "{{.Name}} commented on Your Post {{.Post}}",


		"like_multi": "{{.Name}} & {{.Count}} others liked Your Post {{.Post}}",
		"like_three": "{{.Name}}, {{.Name2}} & {{.Name3}} liked Your Post {{.Post}}",
		"like_two":   "{{.Name}} & {{.Name2}} liked Your Post {{.Post}}",
		"like_one":   "{{.Name}} liked Your Post {{.Post}}",

		"comment_like_multi": "{{.Name}} & {{.Count}} others liked Your Comment {{.Comment}}",
		"comment_like_three": "{{.Name}}, {{.Name2}} & {{.Name3}} liked Your Comment {{.Comment}}",
		"comment_like_two":   "{{.Name}} & {{.Name2}} liked Your Comment {{.Comment}}",
		"comment_like_one":   "{{.Name}} liked Your Comment {{.Comment}}",

		"share_multi": "{{.Name}} & {{.Count}} others shared Your Post {{.Post}}",
		"share_three": "{{.Name}}, {{.Name2}} & {{.Name3}} shared Your Post {{.Post}}",
		"share_two":   "{{.Name}} & {{.Name2}} shared Your Post {{.Post}}",
		"share_one":   "{{.Name}} shared Your Post {{.Post}}",
	},
	"hi": {
		"comment_multi": "{{.Name}} और {{.Count}} लोगों ने आपकी पोस्ट {{.Post}} पर कमेंट किया है ",
		"comment_three": "{{.Name}}, {{.Name2}} और {{.Name3}} ने आपकी पोस्ट {{.Post}} पर कमेंट किया है ",
		"comment_two":   "{{.Name}} और {{.Name2}} ने आपकी पोस्ट {{.Post}} पर कमेंट किया है ",
		"comment_one":   "{{.Name}} ने आपकी पोस्ट {{.Post}} पर कमेंट किया है ",


		"like_multi": "{{.Name}} और {{.Count}} लोगों ने आपकी पोस्ट {{.Post}} को पसंद किया है",
		"like_three": "{{.Name}}, {{.Name2}} और {{.Name3}} ने आपकी पोस्ट {{.Post}}को पसंद किया है",
		"like_two":   "{{.Name}} और {{.Name2}} ने आपकी पोस्ट {{.Post}}को पसंद किया है",
		"like_one":   "{{.Name}} ने आपकी पोस्ट {{.Post}}को पसंद किया है",

		"comment_like_multi": "आपके कमेंट को {{.Count}} लोगों ने पसंद किया है",
		"comment_like_one":   "आपके कमेंट को {{.Name}} ने पसंद किया है",

		"share_multi": "{{.Name}} और {{.Count}} लोगों ने आपकी पोस्ट {{.Post}} को शेयर किया है",
		"share_three": "{{.Name}}, {{.Name2}} और {{.Name3}} ने आपकी पोस्ट {{.Post}} को शेयर किया है",
		"share_two":   "{{.Name}} और {{.Name2}} ने आपकी पोस्ट {{.Post}} को शेयर किया है",
		"share_one":   "{{.Name}} ने आपकी पोस्ट {{.Post}} को शेयर किया है",
	},
	"te": {
		"comment_multi": "{{.Name}} & {{.Count}} others commented on Your Post {{.Post}}",
		"comment_three": "{{.Name}}, {{.Name2}} & {{.Name3}} commented on Your Post {{.Post}}",
		"comment_two":   "{{.Name}} & {{.Name2}} commented on Your Post {{.Post}}",
		"comment_one":   "{{.Name}} commented on Your Post {{.Post}}",


		"like_multi": "{{.Name}} & {{.Count}} others liked Your Post {{.Post}}",
		"like_three": "{{.Name}}, {{.Name2}} & {{.Name3}} liked Your Post {{.Post}}",
		"like_two":   "{{.Name}} & {{.Name2}} liked Your Post {{.Post}}",
		"like_one":   "{{.Name}} liked Your Post {{.Post}}",

		"comment_like_multi": "{{.Name}} & {{.Count}} others liked Your Comment {{.Comment}}",
		"comment_like_three": "{{.Name}}, {{.Name2}} & {{.Name3}} liked Your Comment {{.Comment}}",
		"comment_like_two":   "{{.Name}} & {{.Name2}} liked Your Comment {{.Comment}}",
		"comment_like_one":   "{{.Name}} liked Your Comment {{.Comment}}",

		"share_multi": "{{.Name}} & {{.Count}} others shared Your Post {{.Post}}",
		"share_three": "{{.Name}}, {{.Name2}} & {{.Name3}} shared Your Post {{.Post}}",
		"share_two":   "{{.Name}} & {{.Name2}} shared Your Post {{.Post}}",
		"share_one":   "{{.Name}} shared Your Post {{.Post}}",
	},
}

type DataModel struct {
	Name    string
	Name2   string
	Name3   string
	Count   int
	Post    string
	Comment string
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