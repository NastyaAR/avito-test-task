package middleware

import (
	"avito-test-task/internal/delivery/handlers"
	"net/http"
)

func AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var respBoby []byte
		_, err := r.Cookie("token")
		if err != nil {
			respBoby = handlers.CreateErrorResponse(r.Context(), handlers.ReadCookieError, handlers.ReadCookieErrorMsg)
			w.Write(respBoby)
			w.WriteHeader(http.StatusInternalServerError)
		}
		if err == http.ErrNoCookie {
			respBoby = handlers.CreateErrorResponse(r.Context(), handlers.NotAuthorizedError, handlers.NotAuthorizedErrorMsg)
			w.Write(respBoby)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write([]byte("Response in middleware "))
		handler.ServeHTTP(w, r)

	})
}
