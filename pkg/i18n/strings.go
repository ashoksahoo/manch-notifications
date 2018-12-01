package i18n

var Strings = map[string]map[string]string{
	"en": {
		"comment_multi": "{{.Name}} & {{.Count}} others commented on Your Post {{.Post}}",
		"comment_three": "{{.Name}}, {{.Name2}} & {{.Name3}} commented on Your Post {{.Post}}",
		"comment_two":   "{{.Name}} & {{.Name2}} commented on Your Post {{.Post}}",
		"comment_one":   "{{.Name}} commented on Your Post {{.Post}}",

		"post_like_multi": "{{.Name}} & {{.Count}} others liked Your Post {{.Post}}",
		"post_like_three": "{{.Name}}, {{.Name2}} & {{.Name3}} liked Your Post {{.Post}}",
		"post_like_two":   "{{.Name}} & {{.Name2}} liked Your Post {{.Post}}",
		"post_like_one":   "{{.Name}} liked Your Post {{.Post}}",

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

		"post_like_multi": "{{.Name}} और {{.Count}} लोगों ने आपकी पोस्ट {{.Post}} को पसंद किया है",
		"post_like_three": "{{.Name}}, {{.Name2}} और {{.Name3}} ने आपकी पोस्ट {{.Post}} को पसंद किया है",
		"post_like_two":   "{{.Name}} और {{.Name2}} ने आपकी पोस्ट {{.Post}} को पसंद किया है",
		"post_like_one":   "{{.Name}} ने आपकी पोस्ट {{.Post}} को पसंद किया है",

		"comment_like_multi": "आपके कमेंट को {{.Count}} लोगों ने पसंद किया है",
		"comment_like_three": "{{.Name}}, {{.Name2}} और {{.Name3}} ने आपकी कमेंट {{.Post}} को पसंद किया है",
		"comment_like_two":   "{{.Name}} और {{.Name2}} ने आपकी कमेंट {{.Post}} को पसंद किया है",
		"comment_like_one":   "आपके कमेंट को {{.Name}} ने पसंद किया है",

		"share_multi": "{{.Name}} और {{.Count}} लोगों ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है",
		"share_three": "{{.Name}}, {{.Name2}} और {{.Name3}} ने आपकी पोस्ट {{.Post}} को शेयर किया है",
		"share_two":   "{{.Name}} और {{.Name2}} ने आपकी पोस्ट {{.Post}} को शेयर किया है",
		"share_one":   "{{.Name}} ने आपकी पोस्ट {{.Post}} को शेयर किया है",
	},
	"te": {
		"comment_multi": "{{.Name}} +{{.Count}} వ్యక్తులు మీ పోస్ట్ \"{{.Post}}\" పైన కామెంట్ చేసారు.",
		"comment_three": "{{.Name}}, {{.Name2}} & {{.Name3}} మీ పోస్ట్ \"{{.Post}}\" పైన కామెంట్ చేసారు.",
		"comment_two":   "{{.Name}} & {{.Name2}} మీ పోస్ట్ \"{{.Post}}\" పైన కామెంట్ చేసారు.",
		"comment_one":   "{{.Name}} మీ పోస్ట్ \"{{.Post}}\" పైన కామెంట్ చేసారు.",

		"post_like_multi": "{{.Name}} +{{.Count}} వ్యక్తులు మీ పోస్ట్ \"{{.Post}}\" ని లైక్ చేసారు.",
		"post_like_three": "{{.Name}}, {{.Name2}} & {{.Name3}} మీ పోస్ట్ \"{{.Post}}\" ని లైక్ చేసారు.",
		"post_like_two":   "{{.Name}} & {{.Name2}} మీ పోస్ట్ \"{{.Post}}\" ని లైక్ చేసారు.",
		"post_like_one":   "{{.Name}} మీ {{.Post}} మీ పోస్ట్ \"{{.Post}}\" ని లైక్ చేసారు.",

		"comment_like_multi": "{{.Name}} & {{.Count}} వ్యక్తులు మీ కామెంట్ \"{{.Comment}}\" ని లైక్ చేసారు.",
		"comment_like_three": "{{.Name}}, {{.Name2}} & {{.Name3}} మీ కామెంట్ \"{{.Comment}}\" ని లైక్ చేసారు.",
		"comment_like_two":   "{{.Name}} & {{.Name2}} మీ కామెంట్ \"{{.Comment}}\" ని లైక్ చేసారు.",
		"comment_like_one":   "{{.Name}} మీ కామెంట్ \"{{.Comment}}\" ని లైక్ చేసారు",

		"share_multi": "{{.Name}} &మరియు ఇంకో {{.Count}} వ్యక్తులు మీ పోస్ట్ \"{{.Post}}\" ని షేర్ చేసారు.",
		"share_three": "{{.Name}}, {{.Name2}} & {{.Name3}} మీ పోస్ట్ \"{{.Post}}\" ని షేర్ చేసారు.",
		"share_two":   "{{.Name}} & {{.Name2}} మీ పోస్ట్ \"{{.Post}}\" ని షేర్ చేసారు.",
		"share_one":   "{{.Name}} మీ పోస్ట్ \"{{.Post}}\" ని షేర్ చేసారు.",
	},
}