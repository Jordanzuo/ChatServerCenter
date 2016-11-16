package web

import (
	"encoding/json"
	"github.com/Jordanzuo/ChatServerModel/src/centerResponseObject"
	"net/http"
)

func responseResult(w http.ResponseWriter, responseObj *centerResponseObject.ResponseObject) {
	if responseBytes, err := json.Marshal(responseObj); err == nil {
		w.Write(responseBytes)
	}
}
