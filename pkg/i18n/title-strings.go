package i18n

func GetAppTitle(language string) string {
	if language == "hi" {
		return "मंच"
	} else if language == "te" {
		return "మంచ్"		
	} else if language == "en" {
		return "Manch"
	} else {
		return ""
	}
}