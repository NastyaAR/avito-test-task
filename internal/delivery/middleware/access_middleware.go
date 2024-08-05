package middleware

import (
	"avito-test-task/internal/delivery/handlers"
	"avito-test-task/internal/domain"
	"avito-test-task/pkg"
	"fmt"
	"net/http"
)

func isModeratorOnly(path string) bool {
	return path == "/house/create" || path == "/flat/update"
}

func AccessMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var respBody []byte
		token := r.Header.Get("authorization")
		role, err := pkg.ExtractPayloadFromToken(token, "role")
		if err != nil {
			respBody = handlers.CreateErrorResponse(r.Context(), handlers.ExtractRoleFromTokenError, handlers.ExtractRoleFromTokenErrorMsg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(respBody)
			return
		}

		path := r.URL.Path
		fmt.Println(path)
		if isModeratorOnly(path) && role != domain.Moderator {
			respBody = handlers.CreateErrorResponse(r.Context(), handlers.NoAccessError, handlers.NoAccessErrorMsg)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(respBody)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
