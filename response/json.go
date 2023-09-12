package response

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"reflect"
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

type SuccessData struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
}

type ErrorData struct {
	DeveloperMessage string `json:"developer_message"`
	ErrorCode        int    `json:"error_code"`
	MoreInfo         string `json:"more_info"`
	Status           int    `json:"status"`
	UserMessage      string `json:"user_message"`
}

type PaginateDataSvc struct {
	Data       interface{} `json:"data"`
	Additional interface{} `json:"additional,omitempty"`
	Paginate   struct {
		Total      int  `json:"total"`
		Page       int  `json:"page"`
		Size       int  `json:"size"`
		TotalPages int  `json:"total_pages"`
		NextPage   *int `json:"next_page"`
		PrevPage   *int `json:"prev_page"`
	} `json:"paginate"`
	Code int `json:"code"`
}

func NewPaginateFromSvc(data interface{}, count, page, size int) *PaginateDataSvc {
	var totalPages int = int(math.Ceil(float64(count) / float64(size)))
	var nextPage, prevPage *int

	if page > 1 {
		np := page - 1
		prevPage = &np
	}
	if page == totalPages {
	} else {
		np := page + 1
		nextPage = &np
	}

	dp := PaginateDataSvc{
		Data: data,
		Paginate: struct {
			Total      int  `json:"total"`
			Page       int  `json:"page"`
			Size       int  `json:"size"`
			TotalPages int  `json:"total_pages"`
			NextPage   *int `json:"next_page"`
			PrevPage   *int `json:"prev_page"`
		}{
			count, page, size, totalPages, nextPage, prevPage,
		},
		Code: 0,
	}
	dp.Code = http.StatusOK
	return &dp
}

func Exception(w http.ResponseWriter, err error, code int) {
	/** sentry */
	go sendSentry(err)

	pc, fn, line, _ := runtime.Caller(1)
	log.Printf("[error] %s:%d %v on %s", fn, line, err, pc)
	exception := ErrorData{
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
}

func ExceptionFormatted(w http.ResponseWriter, err error, validator map[string]string, code int) {
	/** sentry */
	go sendSentry(err)

	pc, fn, line, _ := runtime.Caller(1)
	log.Printf("[error] %s:%d %v on %s", fn, line, err, pc)
	moreInfo, err := json.Marshal(validator)
	if err != nil {
		moreInfo = []byte{}
	}
	exception := ErrorData{
		err.Error() + " on " + fn + ":" + strconv.Itoa(line),
		code,
		string(moreInfo),
		code,
		err.Error(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err = json.NewEncoder(w).Encode(
		exception,
	)
}

func Success(w http.ResponseWriter, data interface{}) {
	exception := SuccessData{
		data,
		http.StatusOK,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(
		exception,
	)
}

func Json(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		Exception(w, err, 422)
		return
	}
	r := parseStruct(resp)
	if r == "*Paginator" {
		Paginate(w, resp)
		return
	}
	if r == "*PaginateData" {
		Paginate(w, resp)
		return
	}
	if r == "*PaginatorResponse" {
		Paginate(w, resp)
		return
	}
	if r == "PaginatorSvc" {
		PaginateSvc(w, resp)
		return
	}
	Success(w, resp)
}

func Paginate(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func PaginateSvc(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var pgs PaginatorSvc
	bt, _ := json.Marshal(data)
	err := json.Unmarshal(bt, &pgs)
	if err != nil {
		Exception(w, err, 422)
		return
	}
	pg := NewPaginateFromSvc(pgs.Data, pgs.Total, pgs.Page, pgs.Size)
	json.NewEncoder(w).Encode(pg)
}

/** local func */
func sendSentry(err error) {
	defer sentry.Flush(time.Second * 2)
	sentry.CaptureException(err)
}

func parseStruct(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

/** local struct */
type PaginatorSvc struct {
	Data      []interface{} `json:"data"`
	Page      int           `json:"page"`
	Size      int           `json:"size"`
	Total     int           `json:"total"`
	PageCount int           `json:"page_count"`
}
