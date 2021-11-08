package yardimci

import (
	"log"
	"strings"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ReplaceTrChar(turkce string) string {
	replacer := strings.NewReplacer("ç", "c", "ğ", "g", "ı", "I", "ö", "o", "ş", "s", "ü", "u", "İ", "i", "Ö", "O", "Ü", "U", "Ş", "S", "Ç", "C", "Ğ", "G")
	engChar := replacer.Replace(turkce)
	return engChar
}
