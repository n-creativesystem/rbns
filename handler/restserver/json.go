package restserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/protobuf/runtime/protoiface"
)

var jsonContentType = []string{"application/json; charset=utf-8"}

func Render(c *gin.Context, code int, obj protoiface.MessageV1) {
	c.Render(code, JSON{Data: obj})
}

type JSON struct {
	Data protoiface.MessageV1
}

func (r JSON) Render(w http.ResponseWriter) (err error) {
	if err = WriteJSON(w, r.Data); err != nil {
		panic(err)
	}
	return
}

func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func WriteJSON(w http.ResponseWriter, obj protoiface.MessageV1) error {
	writeContentType(w, jsonContentType)
	m := jsonpb.Marshaler{EmitDefaults: true}
	err := m.Marshal(w, obj)
	return err
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
