package helpers

import (
	"fmt"
	"log"
	"os"
)

func GetSetting(name string) string {
	data := os.Getenv(name)

	if data == "" {
		log.Fatal(fmt.Sprintf("$%s must be set", name))
		return ""
	}

	return data
}
