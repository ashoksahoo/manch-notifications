package i18n

var Strings = map[string]map[string]string{
	"en": {
		"comment_multi": "{{.Name}} & {{.Count}} others commented on Your Post \"{{.Post}}\"",
		"comment_three": "{{.Name}}, {{.Name2}} & {{.Name3}} commented on Your Post \"{{.Post}}\"",
		"comment_two":   "{{.Name}} & {{.Name2}} commented on Your Post \"{{.Post}}\"",
		"comment_one":   "{{.Name}} commented on Your Post \"{{.Post}}\"",

		"comment_reply_multi": "{{.Name}} and {{.Count}} others have replied to your comment \"{{.Comment}}\"",
		"comment_reply_three": "{{.Name}}, {{.Name2}} and {{.Name3}} have replied to your comment \"{{.Comment}}\"",
		"comment_reply_two": "{{.Name}} and {{.Name2}} have replied to your comment \"{{.Comment}}\"",
		"comment_reply_one": "{{.Name}} has replied to your comment \"{{.Comment}}\"",

		"post_like_multi": "{{.Name}} & {{.Count}} others liked Your Post \"{{.Post}}\"",
		"post_like_three": "{{.Name}}, {{.Name2}} & {{.Name3}} liked Your Post \"{{.Post}}\"",
		"post_like_two":   "{{.Name}} & {{.Name2}} liked Your Post \"{{.Post}}\"",
		"post_like_one":   "{{.Name}} liked Your Post \"{{.Post}}\"",

		"comment_like_multi": "{{.Name}} & {{.Count}} others liked Your Comment {{.Comment}}",
		"comment_like_three": "{{.Name}}, {{.Name2}} & {{.Name3}} liked Your Comment {{.Comment}}",
		"comment_like_two":   "{{.Name}} & {{.Name2}} liked Your Comment {{.Comment}}",
		"comment_like_one":   "{{.Name}} liked Your Comment {{.Comment}}",

		"share_multi": "{{.Name}} & {{.Count}} others shared Your Post \"{{.Post}}\"",
		"share_three": "{{.Name}}, {{.Name2}} & {{.Name3}} shared Your Post \"{{.Post}}\"",
		"share_two":   "{{.Name}} & {{.Name2}} shared Your Post \"{{.Post}}\"",
		"share_one":   "{{.Name}} shared Your Post \"{{.Post}}\"",

		"follow_user_multi": "{{.Name}} and {{.Count}} others have started following you on Manch",
		"follow_user_three": "{{.Name}}, {{.Name2}} and {{.Name3}} have started following you on Manch", 
		"follow_user_two": "{{.Name}} and {{.Name2}} have started following you on Manch",
		"follow_user_one": "{{.Name}} has started following you on Manch",
		
		"post_removed": "Dear {{.Name}}, Your post \"{{.Post}}\" ",
		
		"comment_on_same_post_one": "{{.Name}} has also commented on the post  \"{{.Post}}\"",
		"comment_on_same_post_two": "{{.Name}} and {{.Name2}} have also commented on the post \"{{.Post}}\"",
		"comment_on_same_post_three": "{{.Name}}, {{.Name2}} and {{.Name3}} have also commented on the post \"{{.Post}}\"",
		"comment_on_same_post_multi": "{{.Name}} and {{.Count}} other have also commented on the post  \"{{.Post}}\"",

		"share_post_one": "{{.Name}}  ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_two": "{{.Name}} और {{.Name2}}  ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_three": "{{.Name}}, {{.Name2}} और {{.Name3}} ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_multi": "{{.Name}} और {{.Count}} अन्य सदस्य ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",

		"reply_on_same_comment_one": "{{.Name}} has also replied to the comment \"{{.Comment}}\"",
		"reply_on_same_comment_two": "{{.Name}} and {{.Name2}} have also replied to the comment \"{{.Comment}}\"",
		"reply_on_same_comment_three": "{{.Name}}, {{.Name2}} and {{.Name3}} नhave also replied to the comment \"{{.Comment}}\"",
		"reply_on_same_comment_multi": "{{.Name}} and {{.Count}} other have also replied to the comment \"{{.Comment}}\"",
	},

	"hi": {
		"comment_multi": "{{.Name}} और {{.Count}} लोगों ने आपकी पोस्ट \"{{.Post}}\" पर कमेंट किया है ",
		"comment_three": "{{.Name}}, {{.Name2}} और {{.Name3}} ने आपकी पोस्ट \"{{.Post}}\" पर कमेंट किया है ",
		"comment_two":   "{{.Name}} और {{.Name2}} ने आपकी पोस्ट \"{{.Post}}\" पर कमेंट किया है ",
		"comment_one":   "{{.Name}} ने आपकी पोस्ट \"{{.Post}}\" पर कमेंट किया है ",

		"comment_reply_multi": "{{.Name}} और {{.Count}} अन्य सदस्य ने आपके कमेंट \"{{.Comment}}\" का जवाब दिया है",
		"comment_reply_three": "{{.Name}}, {{.Name2}} और {{.Name3}} ने आपके कमेंट \"{{.Comment}}\" का जवाब दिया है",
		"comment_reply_two": "{{.Name}} और {{.Name2}} ने आपके कमेंट \"{{.Comment}}\" का जवाब दिया है",
		"comment_reply_one": "{{.Name}} ने आपके कमेंट \"{{.Comment}}\" का जवाब दिया है",

		"post_like_multi": "{{.Name}} और {{.Count}} लोगों ने आपकी पोस्ट \"{{.Post}}\" को पसंद किया है",
		"post_like_three": "{{.Name}}, {{.Name2}} और {{.Name3}} ने आपकी पोस्ट \"{{.Post}}\" को पसंद किया है",
		"post_like_two":   "{{.Name}} और {{.Name2}} ने आपकी पोस्ट \"{{.Post}}\" को पसंद किया है",
		"post_like_one":   "{{.Name}} ने आपकी पोस्ट \"{{.Post}}\" को पसंद किया है",

		"comment_like_multi": "आपके कमेंट {{.Comment}} को {{.Count}} लोगों ने पसंद किया है",
		"comment_like_three": "{{.Name}}, {{.Name2}} और {{.Name3}} ने आपकी कमेंट \"{{.Comment}}\" को पसंद किया है",
		"comment_like_two":   "{{.Name}} और {{.Name2}} ने आपकी कमेंट \"{{.Comment}}\" को पसंद किया है",
		"comment_like_one":   "आपके कमेंट \"{{.Comment}}\" को {{.Name}} ने पसंद किया है",

		"share_multi": "{{.Name}} और {{.Count}} लोगों ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है",
		"share_three": "{{.Name}}, {{.Name2}} और {{.Name3}} ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है",
		"share_two":   "{{.Name}} और {{.Name2}} ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है",
		"share_one":   "{{.Name}} ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है",

		"follow_user_multi": "{{.Name}} और {{.Count}} अन्य सदस्य अब आपको फॉलो कर रहे हैं",
		"follow_user_three": "{{.Name}}, {{.Name2}} और {{.Name3}} अब आपको फॉलो कर रहे हैं", 
		"follow_user_two": "{{.Name}} और {{.Name2}} अब आपको फॉलो कर रहे हैं",
		"follow_user_one": "{{.Name}} अब आपको फॉलो कर रहे हैं",

		"post_removed": "{{.Name}} जी आपके पोस्ट \"{{.Post}}\" में {{.DeleteReason}} होने के कारण, वह अब लोकप्रिय मंच पर नहीं दिखेगा",

		"comment_on_same_post_one": "{{.Name}} ने भी  पोस्ट \"{{.Post}}\" पर कमेंट किया है",
		"comment_on_same_post_two": "{{.Name}} और {{.Name2}} ने भी  पोस्ट \"{{.Post}}\" पर कमेंट किया है",
		"comment_on_same_post_three": "{{.Name}}, {{.Name2}} और {{.Name3}} ने भी  पोस्ट \"{{.Post}}\" पर कमेंट किया है",
		"comment_on_same_post_multi": "{{.Name}} और {{.Count}} अन्य सदस्य ने भी  पोस्ट \"{{.Post}}\" पर कमेंट किया है",

		"reply_on_same_comment_one": "{{.Name}} ने भी \"{{.Comment}}\" का जवाब दिया है",
		"reply_on_same_comment_two": "{{.Name}} और {{.Name2}} ने भी \"{{.Comment}}\" का जवाब दिया है",
		"reply_on_same_comment_three": "{{.Name}}, {{.Name2}} और {{.Name3}} ने भी \"{{.Comment}}\" का जवाब दिया है",
		"reply_on_same_comment_multi": "{{.Name}} और {{.Count}} अन्य सदस्य ने भी \"{{.Comment}}\" का जवाब दिया है",

		"share_post_one": "{{.Name}}  ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_two": "{{.Name}} और {{.Name2}}  ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_three": "{{.Name}}, {{.Name2}} और {{.Name3}} ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_multi": "{{.Name}} और {{.Count}} अन्य सदस्य ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",

		"tenth_follower_image_1": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Hindi.jpeg",
		"tenth_follower_image_2": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Hindi.jpeg",
		"tenth_follower_image_3": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Hindi.jpeg",
		"tenth_follower_title": "बधाई हो {{.Name}} जी आपके {{.Count}} फ़ॉलोअर्स हो गए हैं !! 🎉",
		"tenth_follower_text_1": "और फ़ॉलोअर्स बनाने के लिए मंच पर चर्चा में शामिल हों ",
		"tenth_follower_text_2": "आपके फ़ॉलोअर्स चाहते है आप मंच पर चर्चा करें ",
		"tenth_follower_text_3": "वे चाहते हैं की आप कुछ पोस्ट या कमेंट करें  ",
	},

	"te": {
		"comment_multi": "{{.Name}} +{{.Count}} వ్యక్తులు మీ పోస్ట్ \"{{.Post}}\" పైన కామెంట్ చేసారు.",
		"comment_three": "{{.Name}}, {{.Name2}} & {{.Name3}} మీ పోస్ట్ \"{{.Post}}\" పైన కామెంట్ చేసారు.",
		"comment_two":   "{{.Name}} & {{.Name2}} మీ పోస్ట్ \"{{.Post}}\" పైన కామెంట్ చేసారు.",
		"comment_one":   "{{.Name}} మీ పోస్ట్ \"{{.Post}}\" పైన కామెంట్ చేసారు.",

		"comment_reply_multi": "{{.Name}} మరియు {{.Count}} వ్యక్తులు మీ కామెంట్  \"{{Titile}}\" కి రిప్లయ్ ఇచ్చారు",
		"comment_reply_three": "{{.Name}}, {{.Name2}} మరియు {{.Name3}}  మీ కామెంట్  \"{{Comment}}\" కి రిప్లయ్ ఇచ్చారు",
		"comment_reply_two": "{{.Name}} మరియు {{.Name2}} మీ కామెంట్  \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు",
		"comment_reply_one": "{{.Name}} మీ కామెంట్ \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు ",

		"post_like_multi": "{{.Name}} +{{.Count}} వ్యక్తులు మీ పోస్ట్ \"{{.Post}}\" ని లైక్ చేసారు.",
		"post_like_three": "{{.Name}}, {{.Name2}} & {{.Name3}} మీ పోస్ట్ \"{{.Post}}\" ని లైక్ చేసారు.",
		"post_like_two":   "{{.Name}} & {{.Name2}} మీ పోస్ట్ \"{{.Post}}\" ని లైక్ చేసారు.",
		"post_like_one":   "{{.Name}} మీ పోస్ట్ \"{{.Post}}\" ని లైక్ చేసారు.",

		"comment_like_multi": "{{.Name}} & {{.Count}} వ్యక్తులు మీ కామెంట్ \"{{.Comment}}\" ని లైక్ చేసారు.",
		"comment_like_three": "{{.Name}}, {{.Name2}} & {{.Name3}} మీ కామెంట్ \"{{.Comment}}\" ని లైక్ చేసారు.",
		"comment_like_two":   "{{.Name}} & {{.Name2}} మీ కామెంట్ \"{{.Comment}}\" ని లైక్ చేసారు.",
		"comment_like_one":   "{{.Name}} మీ కామెంట్ \"{{.Comment}}\" ని లైక్ చేసారు",

		"share_multi": "{{.Name}} & {{.Count}} వ్యక్తులు మీ పోస్ట్ \"{{.Post}}\" ని షేర్ చేసారు.",
		"share_three": "{{.Name}}, {{.Name2}} & {{.Name3}} మీ పోస్ట్ \"{{.Post}}\" ని షేర్ చేసారు.",
		"share_two":   "{{.Name}} & {{.Name2}} మీ పోస్ట్ \"{{.Post}}\" ని షేర్ చేసారు.",
		"share_one":   "{{.Name}} మీ పోస్ట్ \"{{.Post}}\" ని షేర్ చేసారు.",

		"follow_user_multi": "{{.Name}} +{{.Count}} వ్యక్తులు మంచ్ లో మిమ్మల్ని ఫాలో చేస్తున్నారు",
		"follow_user_three": "{{.Name}}, {{.Name2}} మరియు {{.Name3}} మంచ్ లో మిమ్మల్ని ఫాలో చేస్తున్నారు", 
		"follow_user_two": "{{.Name}} మరియు {{.Name2}} మంచ్ లో మిమ్మల్ని ఫాలో చేస్తున్నారు",
		"follow_user_one": "{{.Name}} మంచ్ లో మిమ్మల్ని ఫాలో చేస్తున్నారు",

		"post_removed": "{{.Name}} మీ పోస్టు\"{{.Post}}\" లో{{.DeleteReason}}  కలిగి ఉంది, అందుకే పాపులర్  ఫీడ్ లో కనిపించవు.",

		"comment_on_same_post_one": "{{.Name}} కూడా \"{{.Post}}\" పోస్టుపై కామెంట్ చేసారు",
		"comment_on_same_post_two": "{{.Name}} మరియు {{.Name2}} కూడా \"{{.Post}}\" పోస్టుపై కామెంట్ చేసారు",
		"comment_on_same_post_three": "{{.Name}}, {{.Name2}} మరియు {{.Name3}} కూడా \"{{.Post}}\" పోస్టుపై కామెంట్ చేసారు",
		"comment_on_same_post_multi": "{{.Name}} మరియు ఇంకో {{.Count}} వ్యక్తులు కూడా \"{{.Post}}\" పోస్టుపై కామెంట్ చేసారు",

		"reply_on_same_comment_one": "{{.Name}} గారు మీ కామెంట్ \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు",
		"reply_on_same_comment_two": "{{.Name}} మరియు {{.Name2}} గారు కూడా మీ కామెంట్ \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు",
		"reply_on_same_comment_three": "{{.Name}}, {{.Name2}} మరియు {{.Name3}} గారు కూడా మీ కామెంట్ \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు",
		"reply_on_same_comment_multi": "{{.Name}} మరియు {{.Count}} వ్యక్తులు కూడా మీ కామెంట్ \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు",

		"share_post_one": "{{.Name}} మీ పోస్ట్  \"{{.Post}}\" ని వారి ఫ్రండ్స్ తో  షేర్ చేశారు . మీరు కుడా మీ ఫ్రెండ్స్ తో షేర్ చేయండి మరియూ మీ ఫాలోవర్స్ ని పెంచుకోండి",
        "share_post_two": "{{.Name}} మరియు {{.Name2}}  మీ పోస్ట్\"{{.Post}}\"ని వారి ఫ్రండ్స్ తో  షేర్ చేశారు . మీరు కుడా మీ ఫ్రెండ్స్ తో షేర్ చేయండి మరియూ మీ ఫాలోవర్స్ ని పెంచుకోండి. ",
        "share_post_three": "{{.Name1}}, {{.Name2}} మరియు {{.Name3}}  మీ పోస్ట్\"{{.Post}}\"ని వారి ఫ్రండ్స్ తో  షేర్ చేశారు . మీరు కుడా మీ ఫ్రెండ్స్ తో షేర్ చేయండి మరియూ మీ ఫాలోవర్స్ ని పెంచుకోండి. ",
		"share_post_multi": "{{.Name}}, మరియు +{{.Count}} ఇతరులు   మీ పోస్ట్\"{{.Post}}\"ని వారి ఫ్రండ్స్ తో  షేర్ చేశారు . మీరు కుడా మీ ఫ్రెండ్స్ తో షేర్ చేయండి మరియూ మీ ఫాలోవర్స్ ని పెంచుకోండి.",

		"tenth_follower_title": "కంగ్రాట్స్ {{.Name}} గారు",
		"tenth_follower_text_1": "అప్పుడే {{.Count}} మంది మిమల్ని ఫాల్లో చేస్తున్నారు!! 🎉",
		"tenth_follower_text_2": "అప్పుడే {{.Count}} మంది మిమల్ని ఫాల్లో చేస్తున్నారు!! 🎉",
		"tenth_follower_text_3": "అప్పుడే {{.Count}} మంది మిమల్ని ఫాల్లో చేస్తున్నారు!! 🎉",
		"tenth_follower_image_1": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Telugu+1.jpg",
		"tenth_follower_image_2": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Telugu+2.jpg",
		"tenth_follower_image_3": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Telugu+3.jpg",
	},
}


var DeleteReason = map[string]map[string]string {
	"en": {
		"selfie": "selfi",
		"ads": "ads",
		"abusive-language": "abusive language",
		"obscenity": "Obscenity",
		"other": "Other",
	},
	"hi": {
		"selfie": "सेल्फी",
		"ads": "विज्ञापन",
		"abusive-language": "अभद्र भाषा",
		"obscenity": "अश्लीलता",
		"other": "Other",
	},
	"te": {
		"selfie": "సెల్ఫీ",
		"ads": "యాడ్స్",
		"abusive-language": "దుర్భాష",
		"obscenity": "అసభ్యత",
		"other": "Other",
	},
}