package middleware

import (
	"avito-test-task/internal/delivery/handlers"
	"avito-test-task/internal/domain"
	"avito-test-task/pkg"
	"net/http"
	"regexp"
)

func isModeratorOnly(path string) bool {
	return path == "/house/create" || path == "/flat/update"
}

func isClientOnly(path string) bool {
	matched, _ := regexp.MatchString("/house/[0-9]+/subscribe", path)
	return path == "/flat/create" || matched
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
		if isModeratorOnly(path) && role != domain.Moderator {
			respBody = handlers.CreateErrorResponse(r.Context(), handlers.NoAccessError, handlers.NoAccessErrorMsg)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(respBody)
			return
		}

		if isClientOnly(path) && role != domain.Client {
			respBody = handlers.CreateErrorResponse(r.Context(), handlers.NoAccessError, handlers.NoAccessErrorMsg)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(respBody)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
