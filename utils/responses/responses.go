package responses

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var errorMsg = map[int]string{
	422: "Unprocessable Entity",
	409: "Conflict",
	500: "Internal Server Error",
}

func ErrorJSON(w http.ResponseWriter, code int, desc, remoteAddr string) {
	msg := errorMsg[code]
	log.Printf(">>> response: %d, %s, %s - %s\n", code, msg, desc, remoteAddr)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	response := FailedResponse(msg, desc)
	json.NewEncoder(w).Encode(response)
}

func FailedResponse(msg, desc string) map[string]interface{} {
	return map[string]interface{}{
		"error_message": msg,
		"result":        "failure",
		"value":         "",
		"description":   desc,
		"execute_at":    time.Now().UTC().Add(time.Hour * 9).Format("2006/01/02 15:04:05.000"),
	}
}

func SuccessWithDataResponse(data interface{}, msg string) map[string]interface{} {
	log.Printf(">>> response: %s, %d\n", msg, http.StatusOK)
	return map[string]interface{}{
		"error_message": "",
		"result":        "success",
		"value":         data,
		"description":   msg,
		"execute_at":    time.Now().UTC().Add(time.Hour * 9).Format("2006/01/02 15:04:05.000"),
	}
}

func SuccessWithMultipleDataResponse(data []interface{}, msg string) map[string]interface{} {
	return map[string]interface{}{
		"error_message": "",
		"result":        "success",
		"value":         data,
		"description":   msg,
		"execute_at":    time.Now().UTC().Add(time.Hour * 9).Format("2006/01/02 15:04:05.000"),
	}
}

func SuccessWithDataResponsePagination(data interface{}, currentPage, totalPage int, msg string) map[string]interface{} {
	return map[string]interface{}{
		"error_message": "",
		"result":        "success",
		"value":         data,
		"pagination": map[string]int{
			"current_page": currentPage,
			"total_pages":  totalPage,
		},
		"description": msg,
		"execute_at":  time.Now().UTC().Add(time.Hour * 9).Format("2006/01/02 15:04:05.000"),
	}
}

func SuccessResponse(msg string) map[string]interface{} {
	return map[string]interface{}{
		"error_message": "",
		"result":        "success",
		"value":         "",
		"description":   msg,
		"execute_at":    time.Now().UTC().Add(time.Hour * 9).Format("2006/01/02 15:04:05.000"),
	}
}
