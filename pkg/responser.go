package pkg

import "net/http"

func CreateResponse(w http.ResponseWriter, body []byte, httpCode int) {
	w.Write(respBody)
	w.WriteHeader(http.StatusInternalServerError)
}
