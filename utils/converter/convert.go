package converter

import (
	"log"
	"time"
)

func ConvertStrToDate(s string) time.Time {
	date, err := time.Parse("2006-1-2", s)
	if err != nil {
		log.Println("error ConvertStrToDate(), msg:", err)
	}

	return date
}
