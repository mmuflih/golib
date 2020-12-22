package response

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
)

/**
 * Created by Muhammad Muflih Kholidin
 * at 2020-12-22 09:13:06
 * https://github.com/mmuflih
 * muflic.24@gmail.com
 **/

type SuccessResponse struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
}

type ErrorResponse struct {
	DeveloperMessage string `json:"developer_message"`
	ErrorCode        int    `json:"error_code"`
	MoreInfo         string `json:"more_info"`
	Status           int    `json:"status"`
	UserMessage      string `json:"user_message"`
}

type PaginateResponse struct {
	Data       interface{} `json:"data"`
	Additional interface{} `json:"additional,omitempty"`
	Paginate   struct {
		Count int `json:"total"`
		Page  int `json:"page"`
		Size  int `json:"size"`
	} `json:"paginate"`
	Code int `json:"code"`
}

func NewPaginateResponse(data interface{}, count, page, size int) PaginateResponse {
	dp := PaginateResponse{
		Data: data,
		Paginate: struct {
			Count int `json:"total"`
			Page  int `json:"page"`
			Size  int `json:"size"`
		}{
			count, page, size,
		},
		Code: 0,
	}
	dp.Code = http.StatusOK
	return dp
}

func ResponseException(w http.ResponseWriter, err error, code int) {
	/** sentry */
	go sendSentry(err)

	pc, fn, line, _ := runtime.Caller(1)
	log.Printf("[error] %s:%d %v on %s", fn, line, err, pc)
	exception := ErrorResponse{
		err.Error() + " on " + fn + ":" + strconv.Itoa(line),
		code,
		"Contact developer or administrator",
		code,
		err.Error(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err = json.NewEncoder(w).Encode(
		exception,
	)
	return
}

func ResponseData(w http.ResponseWriter, data interface{}) {
	exception := SuccessResponse{
		data,
		http.StatusOK,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(
		exception,
	)
	return
}

func ResponsePaginate(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
	return
}

func Response() {

}

/** local func */
func sendSentry(err error) {
	defer sentry.Flush(time.Second * 2)
	sentry.CaptureException(err)
}