package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type muxRequestReader struct{}

type RequestReader interface {
	GetRouteParam(r *http.Request, name string) string
	GetRouteParamInt(r *http.Request, name string) int
	GetJsonData(r *http.Request, data interface{}) (err error)
	GetQuery(r *http.Request, query string) string
	GetQueryInt(r *http.Request, query string) int
}

func NewMuxRequestReader() RequestReader {
	return &muxRequestReader{}
}

func (rr *muxRequestReader) GetRouteParam(r *http.Request, name string) string {
	return mux.Vars(r)[name]
}

func (rr *muxRequestReader) GetRouteParamInt(r *http.Request, name string) int {
	i, err := strconv.Atoi(mux.Vars(r)[name])
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return i
}

func (rr *muxRequestReader) GetJsonData(r *http.Request, data interface{}) (err error) {
	err = json.NewDecoder(r.Body).Decode(data)
	return
}

func (rr *muxRequestReader) GetQuery(r *http.Request, key string) string {
	qs := r.URL.Query()
	return qs.Get(key)
}

func (rr *muxRequestReader) GetQueryInt(r *http.Request, key string) int {
	qs := r.URL.Query()
	qi, err := strconv.Atoi(qs.Get(key))
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return qi
}
