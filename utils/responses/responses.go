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
	400: "Bad Request",
	401: "Unauthorized",
}

func ErrorJSON(w http.ResponseWriter, code int, desc interface{}, remoteAddr string) {
	msg := errorMsg[code]
	log.Printf(">>> response: %d, %s, %s - %s\n", code, msg, desc, remoteAddr)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	response := FailedResponse(msg, desc)
	json.NewEncoder(w).Encode(response)
}

func FailedResponse(msg, desc interface{}) map[string]interface{} {
	return map[string]interface{}{
		"error_message": msg,
		"result":        "failure",
		"description":   desc,
		"execute_at":    time.Now().UTC().Add(time.Hour * 9).Format("2006/01/02 15:04:05.000"),
	}
}

func SuccessWithDataResponse(data interface{}, code int, msg string) map[string]interface{} {
	log.Printf(">>> %d, response: %s\n", code, msg)
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
