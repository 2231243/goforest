package response

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

type Response struct {
	Code    int         `json:"code" form:"code"`
	Message string      `json:"message" form:"message"`
	Ts      string      `json:"ts" form:"ts"`
	Data    interface{} `json:"data" form:"data"`
}

type ErrResponse struct {
	Code    int         `json:"code" form:"code"`
	Message string      `json:"message" form:"message"`
	Ts      string      `json:"ts" form:"ts"`
	Reason  string      `json:"reason" form:"reason"`
	Data    interface{} `json:"data" form:"data"`
}

func ErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	se := errors.FromError(err)
	reply := &ErrResponse{}
	reply.Code = int(se.Code)
	reply.Data = nil
	reply.Message = se.Message
	reply.Reason = se.Reason
	reply.Ts = time.Now().Format("2006-01-02 15:04:05.00000")

	codec, _ := khttp.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", strings.Join([]string{"application", codec.Name()}, "/"))
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func ResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	reply := &Response{}
	reply.Code = 0
	reply.Data = v
	reply.Message = "ok"
	reply.Ts = time.Now().Format("2006-01-02 15:04:05.00000")

	codec, _ := khttp.CodecForRequest(r, "Accept")
	data, err := codec.Marshal(reply)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", strings.Join([]string{"application", codec.Name()}, "/"))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return nil
}
