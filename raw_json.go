package svjson

import (
	"encoding/json"
	"github.com/bilibili/kratos/pkg/ecode"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

var jsonContentType = []string{"application/json; charset=utf-8"}

type SVJSON struct {
	Data interface{}
	Err  error
}

type SVError struct {
	Error string `json:"error"`
}

func HandleSVJSON(c *bm.Context, data interface{}, err error) {
	SVJSON{
		Data: data,
		Err:  err,
	}.Handle(c)
}

func (sj SVJSON) Render(w http.ResponseWriter) (err error) {
	var jsonBytes []byte
	if sj.Err != nil {
		if jsonBytes, err = json.Marshal(SVError{sj.Err.Error()}); err != nil {
			err = errors.WithStack(err)
			return
		}
	} else if sj.Data != nil {
		if jsonBytes, err = json.Marshal(sj.Data); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	if _, err = w.Write(jsonBytes); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// WriteContentType write content-type to http response writer.
func (SVJSON) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	header["Content-Type"] = jsonContentType
}

func (sj SVJSON) Handle(c *bm.Context) {
	//code := http.StatusOK
	c.Error = sj.Err
	bcode := ecode.Cause(sj.Err)
	writeStatusCode(c.Writer, bcode.Code())
	c.Render(bcode.Code(), sj)
}

func writeStatusCode(w http.ResponseWriter, ecode int) {
	header := w.Header()
	header.Set("kratos-status-code", strconv.FormatInt(int64(ecode), 10))
}
