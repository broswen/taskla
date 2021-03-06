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
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			limit = DefaultLimit
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			offset = DefaultOffset
		}

		ctx := context.WithValue(r.Context(), "limit", limit)
		ctx = context.WithValue(ctx, "offset", offset)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
