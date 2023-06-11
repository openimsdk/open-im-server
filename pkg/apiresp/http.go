package apiresp

import (
	"encoding/json"
	"net/http"
)

func HttpError(w http.ResponseWriter, err error) {
	data, err := json.Marshal(ParseError(err))
	if err != nil {
		panic(err)
	}
	_ = data

}

func HttpSuccess(w http.ResponseWriter, data any) {

}
