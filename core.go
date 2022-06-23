package goutil

import (
	"fmt"
	"log"
	"strings"
)

func CheckErr(err error, msg string, action int) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err == nil {
		return
	}

	switch action {
	case 0:
		panic(fmt.Sprintf("%s: \n %v", msg, err))
	case 1:
		log.Printf("%s: \n %v", msg, err)
	}
}

func CleanString(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}
