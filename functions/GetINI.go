package functions

import "strings"

func GetINI(iniText string, key string) string {

	keyLoc := strings.Index(iniText, key)
	if keyLoc == -1 {
		return ""
	}
	newLineLoc := strings.Index(iniText[keyLoc:], "\n")
	if newLineLoc == -1 {
		return ""
	}
	keyLocAdjusted := keyLoc + len(key)
	return iniText[keyLocAdjusted+1 : keyLoc+newLineLoc]
}
