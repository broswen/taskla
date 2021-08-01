package middleware

import (
	"context"
	"net/http"
	"strconv"
)

const (
	DefaultLimit  = 50
	DefaultOffset = 0
)

func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
		if err != nil {
			limit = DefaultLimit
		}
		offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
		if err != nil {
			offset = DefaultOffset
		}

		ctx := context.WithValue(r.Context(), "limit", limit)
		ctx = context.WithValue(r.Context(), "offset", offset)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
