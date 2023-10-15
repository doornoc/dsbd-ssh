package tool

import (
	"log"
)

func Deb(debugType string, deb string) {
	if Debug {
		log.Println(debugType, deb)
	}
}

func Err(errorType string, error error) {
	if Debug {
		log.Println(errorType, error)
	}
}
