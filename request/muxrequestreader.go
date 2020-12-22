package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type muxReader struct{}

type Reader interface {
	GetRouteParam(r *http.Request, name string) string
	GetRouteParamInt(r *http.Request, name string) int
	GetJsonData(r *http.Request, data interface{}) (err error)
	GetQuery(r *http.Request, query string) string
	GetQueryInt(r *http.Request, query string) int
}

func NewMuxReader() Reader {
	return &muxReader{}
}

func (rr *muxReade) GetRouteParam(r *http.Request, name string) string {
	return mux.Vars(r)[name]
}

func (rr *muxReader) GetRouteParamInt(r *http.Request, name string) int {
	i, err := strconv.Atoi(mux.Vars(r)[name])
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return i
}

func (rr *muxReader) GetJsonData(r *http.Request, data interface{}) (err error) {
	err = json.NewDecoder(r.Body).Decode(data)
	return
}

func (rr *muxReader) GetQuery(r *http.Request, key string) string {
	qs := r.URL.Query()
	return qs.Get(key)
}

func (rr *muxReader) GetQueryInt(r *http.Request, key string) int {
	qs := r.URL.Query()
	qi, err := strconv.Atoi(qs.Get(key))
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return qi
}
