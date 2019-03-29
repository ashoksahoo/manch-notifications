package i18n

var Strings = map[string]map[string]string{
	"en": {
		"comment_multi": "{{.Name}} & {{.Count}} others commented on Your Post \"{{.Post}}\"",
		"comment_three": "{{.Name}}, {{.Name2}} & {{.Name3}} commented on Your Post \"{{.Post}}\"",
		"comment_two":   "{{.Name}} & {{.Name2}} commented on Your Post \"{{.Post}}\"",
		"comment_one":   "{{.Name}} commented on Your Post \"{{.Post}}\"",

		"comment_reply_multi": "{{.Name}} and {{.Count}} others have replied to your comment \"{{.Comment}}\"",
		"comment_reply_three": "{{.Name}}, {{.Name2}} and {{.Name3}} have replied to your comment \"{{.Comment}}\"",
		"comment_reply_two":   "{{.Name}} and {{.Name2}} have replied to your comment \"{{.Comment}}\"",
		"comment_reply_one":   "{{.Name}} has replied to your comment \"{{.Comment}}\"",

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
		"follow_user_two":   "{{.Name}} and {{.Name2}} have started following you on Manch",
		"follow_user_one":   "{{.Name}} has started following you on Manch",

		"post_removed": "Dear {{.Name}}, Your post \"{{.Post}}\" ",

		"comment_on_same_post_one":   "{{.Name}} has also commented on the post  \"{{.Post}}\"",
		"comment_on_same_post_two":   "{{.Name}} and {{.Name2}} have also commented on the post \"{{.Post}}\"",
		"comment_on_same_post_three": "{{.Name}}, {{.Name2}} and {{.Name3}} have also commented on the post \"{{.Post}}\"",
		"comment_on_same_post_multi": "{{.Name}} and {{.Count}} other have also commented on the post  \"{{.Post}}\"",

		"share_post_one":     "🔥 {{.Name}}  ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_two":     "🔥 {{.Name}} और {{.Name2}}  ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_three":   "🔥 {{.Name}}, {{.Name2}} और {{.Name3}} ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_multi":   "🔥 {{.Name}} और {{.Count}} अन्य सदस्य ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_multi_1": "🔥 अपनी पोस्ट को और भी अधिक वायरल व लोकप्रिय बनाने के लिए तीन दोस्तों के साथ Whatsapp पे शेयर करें",
		"share_post_multi_2": "🔥 पोस्ट वायरल करने के TIPS :पाँच दोस्तों के साथ  Whatsapp पे शेयर करो और मंच पे अधिक लोकप्रिय बनो",
		"share_post_multi_3": "🔥 आपके पोस्ट को {{.Count}} लोगो ने शेयर किया हैं ,आप भी उसको Whatsapp पे शेयर करके ट्रेंडिंग पोस्ट बना सकते हैं ",

		"reply_on_same_comment_one":   "{{.Name}} has also replied to the comment \"{{.Comment}}\"",
		"reply_on_same_comment_two":   "{{.Name}} and {{.Name2}} have also replied to the comment \"{{.Comment}}\"",
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
		"comment_reply_two":   "{{.Name}} और {{.Name2}} ने आपके कमेंट \"{{.Comment}}\" का जवाब दिया है",
		"comment_reply_one":   "{{.Name}} ने आपके कमेंट \"{{.Comment}}\" का जवाब दिया है",

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
		"follow_user_two":   "{{.Name}} और {{.Name2}} अब आपको फॉलो कर रहे हैं",
		"follow_user_one":   "{{.Name}} अब आपको फॉलो कर रहे हैं",

		"post_removed": "{{.Name}} जी आपके पोस्ट \"{{.Post}}\" में {{.DeleteReason}} होने के कारण, वह अब लोकप्रिय मंच पर नहीं दिखेगा",

		"comment_on_same_post_one":   "{{.Name}} ने भी  पोस्ट \"{{.Post}}\" पर कमेंट किया है",
		"comment_on_same_post_two":   "{{.Name}} और {{.Name2}} ने भी  पोस्ट \"{{.Post}}\" पर कमेंट किया है",
		"comment_on_same_post_three": "{{.Name}}, {{.Name2}} और {{.Name3}} ने भी  पोस्ट \"{{.Post}}\" पर कमेंट किया है",
		"comment_on_same_post_multi": "{{.Name}} और {{.Count}} अन्य सदस्य ने भी  पोस्ट \"{{.Post}}\" पर कमेंट किया है",

		"reply_on_same_comment_one":   "{{.Name}} ने भी \"{{.Comment}}\" का जवाब दिया है",
		"reply_on_same_comment_two":   "{{.Name}} और {{.Name2}} ने भी \"{{.Comment}}\" का जवाब दिया है",
		"reply_on_same_comment_three": "{{.Name}}, {{.Name2}} और {{.Name3}} ने भी \"{{.Comment}}\" का जवाब दिया है",
		"reply_on_same_comment_multi": "{{.Name}} और {{.Count}} अन्य सदस्य ने भी \"{{.Comment}}\" का जवाब दिया है",

		"share_post_one":     "🔥 {{.Name}}  ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_two":     "🔥 {{.Name}} और {{.Name2}}  ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_three":   "🔥 {{.Name}}, {{.Name2}} और {{.Name3}} ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_multi":   "🔥 {{.Name}} और {{.Count}} अन्य सदस्य ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_multi_1": "🔥 अपनी पोस्ट को और भी अधिक वायरल व लोकप्रिय बनाने के लिए तीन दोस्तों के साथ Whatsapp पे शेयर करें",
		"share_post_multi_2": "🔥 पोस्ट वायरल करने के TIPS :पाँच दोस्तों के साथ  Whatsapp पे शेयर करो और मंच पे अधिक लोकप्रिय बनो",
		"share_post_multi_3": "🔥 आपके पोस्ट को {{.Count}} लोगो ने शेयर किया हैं ,आप भी उसको Whatsapp पे शेयर करके ट्रेंडिंग पोस्ट बना सकते हैं ",

		"tenth_follower_image_1": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Hindi.jpeg",
		"tenth_follower_image_2": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Hindi.jpeg",
		"tenth_follower_image_3": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Hindi.jpeg",
		"tenth_follower_title":   "बधाई हो {{.Name}} जी आपके {{.Count}} फ़ॉलोअर्स हो गए हैं !! 🎉",
		"tenth_follower_text_1":  "और फ़ॉलोअर्स बनाने के लिए मंच पर चर्चा में शामिल हों ",
		"tenth_follower_text_2":  "आपके फ़ॉलोअर्स चाहते है आप मंच पर चर्चा करें ",
		"tenth_follower_text_3":  "वे चाहते हैं की आप कुछ पोस्ट या कमेंट करें  ",

		"live_topic_winners_title_1": "बधाई हो {{.Name}}, आपके ऊपर हो रही है Coins की वर्षा",
		"live_topic_winners_title_2": "{{.Name}}, बधाई हो आज के लाइव चर्चा के आप हैं विजेता",
		"live_topics_winner_text":    "आपको मिलें हैं {{.Count}} Coins",

		"live_topic_participants_title_1": "{{.Name}}, ये हैं आज की चर्चा के टॉप यूज़र्स",
		"live_topic_participants_title_2": "{{.Name}}, आज की चर्चा के टॉप यूज़र्स की लिस्ट",

		"welcome_message": `नमस्कार {{.Name}} 🙏🏻,

मंच से जुड़ने के लिए धन्यवाद् 💐💐।
*मंच* भारत में बना, भारतीयों के लिए बना पहला हिंदी एप है🇮🇳🇮🇳।
हमारा उद्देश्य सारे भारतीयों को जोड़ना है।

आप हैं हमारे लिए खास। इसलिए हम आपको दे रहें हैं *50 Coins* 💰💰।
Coins का लाभ लेने के लिए इस मैसेज का रिप्लाई करें।😇😇

- मंच परिवार`,

		"100_coin_milestone_title":   "Coins की Century 🤩",
		"100_coin_milestone_text":    "100 Coins तक पहुचन की बधाई 🎉",
		"100_coin_referral_title":    "{{.Name}}, आपको मिले है 100 Referral Coins 🤑",
		"100_coin_referral_text":     "Thanks बोलें अपने दोस्त {{.Name2}} को 😄",
		"100_coin_referral_image":    "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/Referral_100_Coins.jpg",
		"100_coin_milestone_image":   "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/100_Coins-min.jpg",
		"join_manch_request_private": "{{.Name}} आपके मंच {{.Community}} से जुड़ना चाहते है 🤩",
		"join_manch_request_public":  "{{.Name}} जी आपके मंच {{.Community}} से जुड़ गए हैं",
		"join_manch_approved":        "मुबारक हो अब आप {{.Community}} का हिस्सा हैं 😄",
		"manch_activation_title":     "आपका मंच {{.Community}} Activate हो गया है🔥",
		"manch_activation_text":      "पोस्ट करें और अपने मंच को ट्रेंडिंग बनाएं 🤩",
		"manch_100_members":          "शतक के लिए शुभकामना , 100 सदस्य अब आपके मंच से जुड़ चुके हैं  ⚡⚡",
	},

	"te": {
		"comment_multi": "{{.Name}} +{{.Count}} వ్యక్తులు మీ పోస్ట్ \"{{.Post}}\" పైన కామెంట్ చేసారు.",
		"comment_three": "{{.Name}}, {{.Name2}} & {{.Name3}} మీ పోస్ట్ \"{{.Post}}\" పైన కామెంట్ చేసారు.",
		"comment_two":   "{{.Name}} & {{.Name2}} మీ పోస్ట్ \"{{.Post}}\" పైన కామెంట్ చేసారు.",
		"comment_one":   "{{.Name}} మీ పోస్ట్ \"{{.Post}}\" పైన కామెంట్ చేసారు.",

		"comment_reply_multi": "{{.Name}} మరియు {{.Count}} వ్యక్తులు మీ కామెంట్  \"{{Titile}}\" కి రిప్లయ్ ఇచ్చారు",
		"comment_reply_three": "{{.Name}}, {{.Name2}} మరియు {{.Name3}}  మీ కామెంట్  \"{{Comment}}\" కి రిప్లయ్ ఇచ్చారు",
		"comment_reply_two":   "{{.Name}} మరియు {{.Name2}} మీ కామెంట్  \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు",
		"comment_reply_one":   "{{.Name}} మీ కామెంట్ \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు ",

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
		"follow_user_two":   "{{.Name}} మరియు {{.Name2}} మంచ్ లో మిమ్మల్ని ఫాలో చేస్తున్నారు",
		"follow_user_one":   "{{.Name}} మంచ్ లో మిమ్మల్ని ఫాలో చేస్తున్నారు",

		"post_removed": "{{.Name}} మీ పోస్టు\"{{.Post}}\" లో{{.DeleteReason}}  కలిగి ఉంది, అందుకే పాపులర్  ఫీడ్ లో కనిపించవు.",

		"comment_on_same_post_one":   "{{.Name}} కూడా \"{{.Post}}\" పోస్టుపై కామెంట్ చేసారు",
		"comment_on_same_post_two":   "{{.Name}} మరియు {{.Name2}} కూడా \"{{.Post}}\" పోస్టుపై కామెంట్ చేసారు",
		"comment_on_same_post_three": "{{.Name}}, {{.Name2}} మరియు {{.Name3}} కూడా \"{{.Post}}\" పోస్టుపై కామెంట్ చేసారు",
		"comment_on_same_post_multi": "{{.Name}} మరియు ఇంకో {{.Count}} వ్యక్తులు కూడా \"{{.Post}}\" పోస్టుపై కామెంట్ చేసారు",

		"reply_on_same_comment_one":   "{{.Name}} గారు మీ కామెంట్ \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు",
		"reply_on_same_comment_two":   "{{.Name}} మరియు {{.Name2}} గారు కూడా మీ కామెంట్ \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు",
		"reply_on_same_comment_three": "{{.Name}}, {{.Name2}} మరియు {{.Name3}} గారు కూడా మీ కామెంట్ \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు",
		"reply_on_same_comment_multi": "{{.Name}} మరియు {{.Count}} వ్యక్తులు కూడా మీ కామెంట్ \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు",

		"share_post_one":     "🔥 {{.Name}} మీ పోస్ట్  \"{{.Post}}\" ని వారి ఫ్రండ్స్ తో  షేర్ చేశారు . మీరు కుడా మీ ఫ్రెండ్స్ తో షేర్ చేయండి మరియూ మీ ఫాలోవర్స్ ని పెంచుకోండి",
		"share_post_two":     "🔥 {{.Name}} మరియు {{.Name2}}  మీ పోస్ట్\"{{.Post}}\"ని వారి ఫ్రండ్స్ తో  షేర్ చేశారు . మీరు కుడా మీ ఫ్రెండ్స్ తో షేర్ చేయండి మరియూ మీ ఫాలోవర్స్ ని పెంచుకోండి. ",
		"share_post_three":   "🔥 {{.Name1}}, {{.Name2}} మరియు {{.Name3}}  మీ పోస్ట్\"{{.Post}}\"ని వారి ఫ్రండ్స్ తో  షేర్ చేశారు . మీరు కుడా మీ ఫ్రెండ్స్ తో షేర్ చేయండి మరియూ మీ ఫాలోవర్స్ ని పెంచుకోండి. ",
		"share_post_multi":   "🔥 {{.Name}}, మరియు +{{.Count}} ఇతరులు   మీ పోస్ట్\"{{.Post}}\"ని వారి ఫ్రండ్స్ తో  షేర్ చేశారు . మీరు కుడా మీ ఫ్రెండ్స్ తో షేర్ చేయండి మరియూ మీ ఫాలోవర్స్ ని పెంచుకోండి.",
		"share_post_multi_1": "🔥 మీ పోస్ట్ ని మరింత వైరల్ లేదా పాపులర్  చేయడం కోసం ముగ్గరు స్నేహితులకి Whatsapp లో షేర్ చేయండి",
		"share_post_multi_2": "🔥 మీ పోస్ట్ పాపులర్ కావడానికి TIPS. Whatsapp లో 5 ఫ్రెండ్స్ కి షేర్ చేయండి.",
		"share_post_multi_3": "🔥 మీ పోస్ట్ ని {{count}} సభ్యులు షేర్ చేసారు, మీరు కూడా Whatsapp లో షేర్ చేయండి మీ పోస్ట్ ని ట్రెండింగ్ పోస్ట్ గ చేయండి.",

		"tenth_follower_title":   "కంగ్రాట్స్ {{.Name}} గారు",
		"tenth_follower_text_1":  "అప్పుడే {{.Count}} మంది మిమల్ని ఫాల్లో చేస్తున్నారు!! 🎉",
		"tenth_follower_text_2":  "అప్పుడే {{.Count}} మంది మిమల్ని ఫాల్లో చేస్తున్నారు!! 🎉",
		"tenth_follower_text_3":  "అప్పుడే {{.Count}} మంది మిమల్ని ఫాల్లో చేస్తున్నారు!! 🎉",
		"tenth_follower_image_1": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Telugu+1.jpg",
		"tenth_follower_image_2": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Telugu+2.jpg",
		"tenth_follower_image_3": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Telugu+3.jpg",

		"live_topic_winners_title_1": "congratulations {{.Name}} గారు, చిరుజల్లులలాంటి Coins మీకోసం",
		"live_topic_winners_title_2": "congratulations {{.Name}} గారు,  ఈరోజు LIVE చర్చ విజేతలు మీరే",
		"live_topics_winner_text":    "మీరు {{.Count}} Coins గెలుచుకున్నారు",

		"live_topic_participants_title_1": "{{.Name}} గారు ,  నేటి  LIVE చర్చ టాప్ యూజర్లు",
		"live_topic_participants_title_2": "{{.Name}} గారు , నేటి  LIVE చర్చ టాప్ యూజర్ల వివరాలు",
		"welcome_message": `Hello {{.Name}} 🙏🏻, 

Manch ఆప్ కి స్వాగతం 💐💐,
Manch, ఇండియా లో మొదటి తెలుగు ఆప్ 🇮🇳🇮🇳.  తెలుగు జాతిని ఒక ఆప్ పై కలపడం Manch ముఖ్య ఉద్దేశం.

మీరు మాకు ప్రత్యేకం. కావున మీకు * 50 Coins * 💰💰 ఇస్తున్నాము. 
Coins  పొందుటకు ఈ మెసేజ్ కి రిప్లై చేయండి. 

- Manch టీం`,

		"100_coin_milestone_title":   "Coins🤩 సెంచరీ",
		"100_coin_milestone_text":    "అభినందనలు 100 mark🎉 కొట్టినందుకు",
		"100_coin_referral_title":    "{{.Name}}, మీకు 100 రెఫరల్ Coins 🤩వచ్చాయి",
		"100_coin_referral_text":     "మీ స్నేహితులకు ధన్యవాదాలు 😃 {{.Name2}}",
		"100_coin_referral_image":    "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/Referral_100_Coins.jpg",
		"100_coin_milestone_image":   "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/100_Coins-min.jpg",
		"join_manch_request_private": "{{.Name}} మీ {{.Community}} మంచ్ లో చేరాలనుకుంటున్నారు. 🤩",
		"join_manch_request_public":  "{{.Name}} మీ మంచ్ లో చేరారు {{.Community}}",
		"join_manch_approved":        "{{.Community}} లో చేరడానికి మీ అభ్యర్థన ఆమోదించబడింది. 😄",
		"manch_activation_title":     "Congratulations, మీ మంచ్ {{.Community}} ఆక్టివేట్ చేయబడింది!!",
		"manch_activation_text":      "మీ అభిప్రాయాలను  పంచుకొని మీ పోస్ట్ ని ట్రేండింగ్ లో ఉండేలా చూసుకోండి.",
		"manch_100_members":          "సెంచరీ! మీ  {{.Community}} మంచ్ ఇప్పుడు 100 మంది సభ్యులను కలిగి ఉన్నారు ⚡⚡",
	},
}

var DeleteReason = map[string]map[string]string{
	"en": {
		"selfie":           "selfi",
		"ads":              "ads",
		"abusive-language": "abusive language",
		"obscenity":        "Obscenity",
		"other":            "Other",
	},
	"hi": {
		"selfie":           "सेल्फी",
		"ads":              "विज्ञापन",
		"abusive-language": "अभद्र भाषा",
		"obscenity":        "अश्लीलता",
		"other":            "Other",
	},
	"te": {
		"selfie":           "సెల్ఫీ",
		"ads":              "యాడ్స్",
		"abusive-language": "దుర్భాష",
		"obscenity":        "అసభ్యత",
		"other":            "Other",
	},
}

var HtmlStrings = map[string]map[string]string{
	"en": {
		"comment_multi": "<b>{{.Name}}</b> & <b>{{.Count}}</b> others commented on Your Post \"{{.Post}}\"",
		"comment_three": "<b>{{.Name}}</b>, <b>{{.Name2}}</b> & <b>{{.Name3}}</b> commented on Your Post \"{{.Post}}\"",
		"comment_two":   "<b>{{.Name}}</b> & <b>{{.Name2}}</b> commented on Your Post \"{{.Post}}\"",
		"comment_one":   "<b>{{.Name}}</b> commented on Your Post \"{{.Post}}\"",

		"comment_reply_multi": "<b>{{.Name}}</b> and <b>{{.Count}}</b> others have replied to your comment \"{{.Comment}}\"",
		"comment_reply_three": "<b>{{.Name}}</b>, <b>{{.Name2}}</b> and <b>{{.Name3}}</b> have replied to your comment \"{{.Comment}}\"",
		"comment_reply_two":   "<b>{{.Name}}</b> and <b>{{.Name2}}</b> have replied to your comment \"{{.Comment}}\"",
		"comment_reply_one":   "<b>{{.Name}}</b> has replied to your comment \"{{.Comment}}\"",

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
		"follow_user_two":   "{{.Name}} and {{.Name2}} have started following you on Manch",
		"follow_user_one":   "{{.Name}} has started following you on Manch",

		"post_removed": "Dear {{.Name}}, Your post \"{{.Post}}\" ",

		"comment_on_same_post_one":   "{{.Name}} has also commented on the post  \"{{.Post}}\"",
		"comment_on_same_post_two":   "{{.Name}} and {{.Name2}} have also commented on the post \"{{.Post}}\"",
		"comment_on_same_post_three": "{{.Name}}, {{.Name2}} and {{.Name3}} have also commented on the post \"{{.Post}}\"",
		"comment_on_same_post_multi": "{{.Name}} and {{.Count}} other have also commented on the post  \"{{.Post}}\"",

		"share_post_one":     "🔥 {{.Name}}  ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_two":     "🔥 {{.Name}} और {{.Name2}}  ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_three":   "🔥 {{.Name}}, {{.Name2}} और {{.Name3}} ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_multi":   "🔥 {{.Name}} और {{.Count}} अन्य सदस्य ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_multi_1": "🔥 अपनी पोस्ट को और भी अधिक वायरल व लोकप्रिय बनाने के लिए तीन दोस्तों के साथ Whatsapp पे शेयर करें",
		"share_post_multi_2": "🔥 पोस्ट वायरल करने के TIPS :पाँच दोस्तों के साथ  Whatsapp पे शेयर करो और मंच पे अधिक लोकप्रिय बनो",
		"share_post_multi_3": "🔥 आपके पोस्ट को {{.Count}} लोगो ने शेयर किया हैं ,आप भी उसको Whatsapp पे शेयर करके ट्रेंडिंग पोस्ट बना सकते हैं ",

		"reply_on_same_comment_one":   "{{.Name}} has also replied to the comment \"{{.Comment}}\"",
		"reply_on_same_comment_two":   "{{.Name}} and {{.Name2}} have also replied to the comment \"{{.Comment}}\"",
		"reply_on_same_comment_three": "{{.Name}}, {{.Name2}} and {{.Name3}} नhave also replied to the comment \"{{.Comment}}\"",
		"reply_on_same_comment_multi": "{{.Name}} and {{.Count}} other have also replied to the comment \"{{.Comment}}\"",
	},

	"hi": {
		"comment_multi": "<b>{{.Name}}</b> और <b>{{.Count}}</b> लोगों ने आपकी पोस्ट \"{{.Post}}\" पर कमेंट किया है ",
		"comment_three": "<b>{{.Name}}</b>, <b>{{.Name2}}</b> और <b>{{.Name3}}</b> ने आपकी पोस्ट \"{{.Post}}\" पर कमेंट किया है ",
		"comment_two":   "<b>{{.Name}}</b> और <b>{{.Name2}}</b> ने आपकी पोस्ट \"{{.Post}}\" पर कमेंट किया है ",
		"comment_one":   "<b>{{.Name}}</b> ने आपकी पोस्ट \"{{.Post}}\" पर कमेंट किया है ",

		"comment_reply_multi": "<b>{{.Name}}</b> और <b>{{.Count}}</b> अन्य सदस्य ने आपके कमेंट \"{{.Comment}}\" का जवाब दिया है",
		"comment_reply_three": "<b>{{.Name}}</b>, <b>{{.Name2}}</b> और <b>{{.Name3}}</b> ने आपके कमेंट \"{{.Comment}}\" का जवाब दिया है",
		"comment_reply_two":   "<b>{{.Name}}</b> और <b>{{.Name2}}</b> ने आपके कमेंट \"{{.Comment}}\" का जवाब दिया है",
		"comment_reply_one":   "<b>{{.Name}}</b> ने आपके कमेंट \"{{.Comment}}\" का जवाब दिया है",

		"post_like_multi": "<b>{{.Name}}</b> और <b>{{.Count}}</b> लोगों ने आपकी पोस्ट \"{{.Post}}\" को पसंद किया है",
		"post_like_three": "<b>{{.Name}}</b>, <b>{{.Name2}}</b> और <b>{{.Name3}}</b> ने आपकी पोस्ट \"{{.Post}}\" को पसंद किया है",
		"post_like_two":   "<b>{{.Name}}</b> और <b>{{.Name2}}</b> ने आपकी पोस्ट \"{{.Post}}\" को पसंद किया है",
		"post_like_one":   "<b>{{.Name}}</b> ने आपकी पोस्ट \"{{.Post}}\" को पसंद किया है",

		"comment_like_multi": "आपके कमेंट {{.Comment}} को <b>{{.Count}}</b> लोगों ने पसंद किया है",
		"comment_like_three": "<b>{{.Name}}</b>, <b>{{.Name2}}</b> और <b>{{.Name3}}</b> ने आपकी कमेंट \"{{.Comment}}\" को पसंद किया है",
		"comment_like_two":   "<b>{{.Name}}</b> और <b>{{.Name2}}</b> ने आपकी कमेंट \"{{.Comment}}\" को पसंद किया है",
		"comment_like_one":   "आपके कमेंट \"{{.Comment}}\" को <b>{{.Name}} ने पसंद किया है",

		"share_multi": "<b>{{.Name}}</b> और <b>{{.Count}}</b> लोगों ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है",
		"share_three": "<b>{{.Name}}</b>, <b>{{.Name2}}</b> और <b>{{.Name3}}</b> ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है",
		"share_two":   "<b>{{.Name}}</b> और <b>{{.Name2}}</b> ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है",
		"share_one":   "<b>{{.Name}}</b> ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है",

		"follow_user_multi": "<b>{{.Name}}</b> और <b>{{.Count}}</b> अन्य सदस्य अब आपको फॉलो कर रहे हैं",
		"follow_user_three": "<b>{{.Name}}</b>, <b>{{.Name2}}</b> और <b>{{.Name3}}</b> अब आपको फॉलो कर रहे हैं",
		"follow_user_two":   "<b>{{.Name}}</b> और <b>{{.Name2}}</b> अब आपको फॉलो कर रहे हैं",
		"follow_user_one":   "<b>{{.Name}}</b> अब आपको फॉलो कर रहे हैं",

		"post_removed": "<b>{{.Name}}</b> जी आपके पोस्ट \"{{.Post}}\" में {{.DeleteReason}} होने के कारण, वह अब लोकप्रिय मंच पर नहीं दिखेगा",

		"comment_on_same_post_one":   "<b>{{.Name}}</b> ने भी  पोस्ट \"{{.Post}}\" पर कमेंट किया है",
		"comment_on_same_post_two":   "<b>{{.Name}}</b> और <b>{{.Name2}}</b> ने भी  पोस्ट \"{{.Post}}\" पर कमेंट किया है",
		"comment_on_same_post_three": "<b>{{.Name}}</b>, <b>{{.Name2}}</b> और <b>{{.Name3}}</b> ने भी  पोस्ट \"{{.Post}}\" पर कमेंट किया है",
		"comment_on_same_post_multi": "<b>{{.Name}}</b> और <b>{{.Count}}</b> अन्य सदस्य ने भी  पोस्ट \"{{.Post}}\" पर कमेंट किया है",

		"reply_on_same_comment_one":   "<b>{{.Name}}</b> ने भी \"{{.Comment}}\" का जवाब दिया है",
		"reply_on_same_comment_two":   "<b>{{.Name}}</b> और <b>{{.Name2}}</b> ने भी \"{{.Comment}}\" का जवाब दिया है",
		"reply_on_same_comment_three": "<b>{{.Name}}</b>, <b>{{.Name2}}</b> और <b>{{.Name3}}</b> ने भी \"{{.Comment}}\" का जवाब दिया है",
		"reply_on_same_comment_multi": "<b>{{.Name}}</b> और <b>{{.Count}}</b> अन्य सदस्य ने भी \"{{.Comment}}\" का जवाब दिया है",

		"share_post_one":     "🔥 <b>{{.Name}}</b>  ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_two":     "🔥 <b>{{.Name}}</b> और <b>{{.Name2}}</b>  ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_three":   "🔥 <b>{{.Name}}</b>, <b>{{.Name2}}</b> और <b>{{.Name3}}</b> ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_multi":   "🔥 <b>{{.Name}}</b> और {{.Count}} अन्य सदस्य ने आपकी पोस्ट \"{{.Post}}\" को शेयर किया है. आप भी अपने पोस्ट को दोस्तों के साथ शेयर करें और अपने फॉलोअर्स बढ़ायें.",
		"share_post_multi_1": "🔥 अपनी पोस्ट को और भी अधिक वायरल व लोकप्रिय बनाने के लिए तीन दोस्तों के साथ Whatsapp पे शेयर करें",
		"share_post_multi_2": "🔥 पोस्ट वायरल करने के TIPS :पाँच दोस्तों के साथ  Whatsapp पे शेयर करो और मंच पे अधिक लोकप्रिय बनो",
		"share_post_multi_3": "🔥 आपके पोस्ट को <b>{{.Count}}</b> लोगो ने शेयर किया हैं ,आप भी उसको Whatsapp पे शेयर करके ट्रेंडिंग पोस्ट बना सकते हैं ",

		"tenth_follower_image_1": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Hindi.jpeg",
		"tenth_follower_image_2": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Hindi.jpeg",
		"tenth_follower_image_3": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Hindi.jpeg",
		"tenth_follower_title":   "बधाई हो <b>{{.Name}}</b> जी आपके <b>{{.Count}}</b> फ़ॉलोअर्स हो गए हैं !! 🎉",
		"tenth_follower_text_1":  "और फ़ॉलोअर्स बनाने के लिए मंच पर चर्चा में शामिल हों ",
		"tenth_follower_text_2":  "आपके फ़ॉलोअर्स चाहते है आप मंच पर चर्चा करें ",
		"tenth_follower_text_3":  "वे चाहते हैं की आप कुछ पोस्ट या कमेंट करें  ",

		"live_topic_winners_title_1": "बधाई हो <b>{{.Name}}</b>, आपके ऊपर हो रही है Coins की वर्षा",
		"live_topic_winners_title_2": "<b>{{.Name}}</b>, बधाई हो आज के लाइव चर्चा के आप हैं विजेता",
		"live_topics_winner_text":    "आपको मिलें हैं <b>{{.Count}}</b> Coins",

		"live_topic_participants_title_1": "<b>{{.Name}}</b>, ये हैं आज की चर्चा के टॉप यूज़र्स",
		"live_topic_participants_title_2": "<b>{{.Name}}</b>, आज की चर्चा के टॉप यूज़र्स की लिस्ट",

		"100_coin_milestone_title":   "Coins की Century 🤩",
		"100_coin_milestone_text":    "100 Coins तक पहुचन की बधाई 🎉",
		"100_coin_referral_title":    "<b>{{.Name}}</b>, आपको मिले है 100 Referral Coins 🤑",
		"100_coin_referral_text":     "Thanks बोलें अपने दोस्त <b>{{.Name2}}</b> को 😄",
		"join_manch_request_private": "<b>{{.Name}}</b> आपके मंच <b>{{.Community}}</b> से जुड़ना चाहते है 🤩",
		"join_manch_request_public":  "<b>{{.Name}}</b> जी आपके मंच <b>{{.Community}}</b> से जुड़ गए हैं",
		"join_manch_approved":        "मुबारक हो अब आप <b>{{.Community}}</b> का हिस्सा हैं 😄",
		"manch_activation_title":     "आपका मंच <b>{{.Community}}</b> Activate हो गया है🔥",
		"manch_activation_text":      "पोस्ट करें और अपने मंच को ट्रेंडिंग बनाएं 🤩",
		"manch_100_members":          "शतक के लिए शुभकामना , <b>100</b> सदस्य अब आपके मंच से जुड़ चुके हैं  ⚡⚡",
	},

	"te": {
		"comment_multi": "<b>{{.Name}}</b> <b>+{{.Count}}</b> వ్యక్తులు మీ పోస్ట్ \"{{.Post}}\" పైన కామెంట్ చేసారు.",
		"comment_three": "<b>{{.Name}}</b>, <b>{{.Name2}}</b> & <b>{{.Name3}}</b> మీ పోస్ట్ \"{{.Post}}\" పైన కామెంట్ చేసారు.",
		"comment_two":   "<b>{{.Name}}</b> & <b>{{.Name2}}</b> మీ పోస్ట్ \"{{.Post}}\" పైన కామెంట్ చేసారు.",
		"comment_one":   "<b>{{.Name}}</b> మీ పోస్ట్ \"{{.Post}}\" పైన కామెంట్ చేసారు.",

		"comment_reply_multi": "<b>{{.Name}}</b> మరియు {{.Count}} వ్యక్తులు మీ కామెంట్  \"{{Titile}}\" కి రిప్లయ్ ఇచ్చారు",
		"comment_reply_three": "<b>{{.Name}}</b>, {{.Name2}} మరియు {{.Name3}}  మీ కామెంట్  \"{{Comment}}\" కి రిప్లయ్ ఇచ్చారు",
		"comment_reply_two":   "<b>{{.Name}}</b> మరియు {{.Name2}} మీ కామెంట్  \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు",
		"comment_reply_one":   "<b>{{.Name}}</b> మీ కామెంట్ \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు ",

		"post_like_multi": "<b>{{.Name}}</b> <b>+{{.Count}}</b> వ్యక్తులు మీ పోస్ట్ \"{{.Post}}\" ని లైక్ చేసారు.",
		"post_like_three": "<b>{{.Name}}</b>, <b>{{.Name2}}</b> & <b>{{.Name3}}</b> మీ పోస్ట్ \"{{.Post}}\" ని లైక్ చేసారు.",
		"post_like_two":   "<b>{{.Name}}</b> & <b>{{.Name2}}</b> మీ పోస్ట్ \"{{.Post}}\" ని లైక్ చేసారు.",
		"post_like_one":   "<b>{{.Name}}</b> మీ పోస్ట్ \"{{.Post}}\" ని లైక్ చేసారు.",

		"comment_like_multi": "<b>{{.Name}}</b> & <b>{{.Count}}</b> వ్యక్తులు మీ కామెంట్ \"{{.Comment}}\" ని లైక్ చేసారు.",
		"comment_like_three": "<b>{{.Name}}</b>, <b>{{.Name2}}</b> & <b>{{.Name3}}</b> మీ కామెంట్ \"{{.Comment}}\" ని లైక్ చేసారు.",
		"comment_like_two":   "<b>{{.Name}}</b> & <b>{{.Name2}}</b> మీ కామెంట్ \"{{.Comment}}\" ని లైక్ చేసారు.",
		"comment_like_one":   "<b>{{.Name}}</b> మీ కామెంట్ \"{{.Comment}}\" ని లైక్ చేసారు",

		"share_multi": "<b>{{.Name}}</b> & <b>{{.Count}}</b> వ్యక్తులు మీ పోస్ట్ \"{{.Post}}\" ని షేర్ చేసారు.",
		"share_three": "<b>{{.Name}}</b>, <b>{{.Name2}}</b> & <b>{{.Name3}}</b> మీ పోస్ట్ \"{{.Post}}\" ని షేర్ చేసారు.",
		"share_two":   "<b>{{.Name}}</b> & <b>{{.Name2}}</b> మీ పోస్ట్ \"{{.Post}}\" ని షేర్ చేసారు.",
		"share_one":   "<b>{{.Name}}</b> మీ పోస్ట్ \"{{.Post}}\" ని షేర్ చేసారు.",

		"follow_user_multi": "<b>{{.Name}}</b> <b>+{{.Count}}</b> వ్యక్తులు మంచ్ లో మిమ్మల్ని ఫాలో చేస్తున్నారు",
		"follow_user_three": "<b>{{.Name}}</b>, <b>{{.Name2}}</b> మరియు <b>{{.Name3}}</b> మంచ్ లో మిమ్మల్ని ఫాలో చేస్తున్నారు",
		"follow_user_two":   "<b>{{.Name}}</b> మరియు <b>{{.Name2}}</b> మంచ్ లో మిమ్మల్ని ఫాలో చేస్తున్నారు",
		"follow_user_one":   "<b>{{.Name}}</b> మంచ్ లో మిమ్మల్ని ఫాలో చేస్తున్నారు",

		"post_removed": "{{.Name}} మీ పోస్టు\"{{.Post}}\" లో{{.DeleteReason}}  కలిగి ఉంది, అందుకే పాపులర్  ఫీడ్ లో కనిపించవు.",

		"comment_on_same_post_one":   "<b>{{.Name}}</b> కూడా \"{{.Post}}\" పోస్టుపై కామెంట్ చేసారు",
		"comment_on_same_post_two":   "<b>{{.Name}}</b> మరియు {{.Name2}} కూడా \"{{.Post}}\" పోస్టుపై కామెంట్ చేసారు",
		"comment_on_same_post_three": "<b>{{.Name}}</b>, {{.Name2}} మరియు {{.Name3}} కూడా \"{{.Post}}\" పోస్టుపై కామెంట్ చేసారు",
		"comment_on_same_post_multi": "<b>{{.Name}}</b> మరియు ఇంకో {{.Count}} వ్యక్తులు కూడా \"{{.Post}}\" పోస్టుపై కామెంట్ చేసారు",

		"reply_on_same_comment_one":   "<b>{{.Name}}</b> గారు మీ కామెంట్ \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు",
		"reply_on_same_comment_two":   "<b>{{.Name}}</b> మరియు {{.Name2}} గారు కూడా మీ కామెంట్ \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు",
		"reply_on_same_comment_three": "<b>{{.Name}}</b>, {{.Name2}} మరియు {{.Name3}} గారు కూడా మీ కామెంట్ \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు",
		"reply_on_same_comment_multi": "<b>{{.Name}}</b> మరియు {{.Count}} వ్యక్తులు కూడా మీ కామెంట్ \"{{.Comment}}\" కి రిప్లయ్ ఇచ్చారు",

		"share_post_one":     "🔥 <b>{{.Name}}</b> మీ పోస్ట్  \"{{.Post}}\" ని వారి ఫ్రండ్స్ తో  షేర్ చేశారు . మీరు కుడా మీ ఫ్రెండ్స్ తో షేర్ చేయండి మరియూ మీ ఫాలోవర్స్ ని పెంచుకోండి",
		"share_post_two":     "🔥 <b>{{.Name}}</b> మరియు <b>{{.Name2}}</b>  మీ పోస్ట్\"{{.Post}}\"ని వారి ఫ్రండ్స్ తో  షేర్ చేశారు . మీరు కుడా మీ ఫ్రెండ్స్ తో షేర్ చేయండి మరియూ మీ ఫాలోవర్స్ ని పెంచుకోండి. ",
		"share_post_three":   "🔥 <b>{{.Name1}}</b>, <b>{{.Name2}}</b> మరియు <b>{{.Name3}}</b>  మీ పోస్ట్\"{{.Post}}\"ని వారి ఫ్రండ్స్ తో  షేర్ చేశారు . మీరు కుడా మీ ఫ్రెండ్స్ తో షేర్ చేయండి మరియూ మీ ఫాలోవర్స్ ని పెంచుకోండి. ",
		"share_post_multi":   "🔥 <b>{{.Name}}</b>, మరియు <b>+{{.Count}}</b> ఇతరులు   మీ పోస్ట్\"{{.Post}}\"ని వారి ఫ్రండ్స్ తో  షేర్ చేశారు . మీరు కుడా మీ ఫ్రెండ్స్ తో షేర్ చేయండి మరియూ మీ ఫాలోవర్స్ ని పెంచుకోండి.",
		"share_post_multi_1": "🔥 మీ పోస్ట్ ని మరింత వైరల్ లేదా పాపులర్  చేయడం కోసం ముగ్గరు స్నేహితులకి Whatsapp లో షేర్ చేయండి",
		"share_post_multi_2": "🔥 మీ పోస్ట్ పాపులర్ కావడానికి TIPS. Whatsapp లో 5 ఫ్రెండ్స్ కి షేర్ చేయండి.",
		"share_post_multi_3": "🔥 మీ పోస్ట్ ని <b>{{count}}</b> సభ్యులు షేర్ చేసారు, మీరు కూడా Whatsapp లో షేర్ చేయండి మీ పోస్ట్ ని ట్రెండింగ్ పోస్ట్ గ చేయండి.",

		"tenth_follower_title":   "కంగ్రాట్స్ <b>{{.Name}}</b> గారు",
		"tenth_follower_text_1":  "అప్పుడే <b>{{.Count}}</b> మంది మిమల్ని ఫాల్లో చేస్తున్నారు!! 🎉",
		"tenth_follower_text_2":  "అప్పుడే <b>{{.Count}}</b> మంది మిమల్ని ఫాల్లో చేస్తున్నారు!! 🎉",
		"tenth_follower_text_3":  "అప్పుడే <b>{{.Count}}</b> మంది మిమల్ని ఫాల్లో చేస్తున్నారు!! 🎉",
		"tenth_follower_image_1": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Telugu+1.jpg",
		"tenth_follower_image_2": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Telugu+2.jpg",
		"tenth_follower_image_3": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/10+Followers+Telugu+3.jpg",

		"live_topic_winners_title_1": "congratulations <b>{{.Name}}</b> గారు, చిరుజల్లులలాంటి Coins మీకోసం",
		"live_topic_winners_title_2": "congratulations <b>{{.Name}}</b> గారు,  ఈరోజు LIVE చర్చ విజేతలు మీరే",
		"live_topics_winner_text":    "మీరు <b>{{.Count}}</b> Coins గెలుచుకున్నారు",

		"live_topic_participants_title_1": "<b>{{.Name}}</b> గారు ,  నేటి  LIVE చర్చ టాప్ యూజర్లు",
		"live_topic_participants_title_2": "<b>{{.Name}}</b> గారు , నేటి  LIVE చర్చ టాప్ యూజర్ల వివరాలు",

		"100_coin_milestone_title":   "Coins🤩 సెంచరీ",
		"100_coin_milestone_text":    "అభినందనలు 100 mark🎉 కొట్టినందుకు",
		"100_coin_referral_title":    "<b>{{.Name}}</b>, మీకు 100 రెఫరల్ Coins 🤩వచ్చాయి",
		"100_coin_referral_text":     "మీ స్నేహితులకు ధన్యవాదాలు 😃 <b>{{.Name2}}</b>",
		"join_manch_request_private": "<b>{{.Name}}</b> మీ <b>{{.Community}}</b> మంచ్ లో చేరాలనుకుంటున్నారు. 🤩",
		"join_manch_request_public":  "<b>{{.Name}}</b> మీ మంచ్ లో చేరారు <b>{{.Community}}</b>",
		"join_manch_approved":        "<b>{{.Community}}</b> లో చేరడానికి మీ అభ్యర్థన ఆమోదించబడింది. 😄",
		"manch_activation_title":     "Congratulations, మీ మంచ్ <b>{{.Community}}</b> ఆక్టివేట్ చేయబడింది!!",
		"manch_activation_text":      "మీ అభిప్రాయాలను  పంచుకొని మీ పోస్ట్ ని ట్రేండింగ్ లో ఉండేలా చూసుకోండి.",
		"manch_100_members":          "సెంచరీ! మీ  <b>{{.Community}}</b> మంచ్ ఇప్పుడు <b>100</b> మంది సభ్యులను కలిగి ఉన్నారు ⚡⚡",
	},
}
