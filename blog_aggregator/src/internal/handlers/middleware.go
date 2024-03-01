package handlers

import (
	"net/http"

	"github.com/daneofmanythings/blog_aggregator/internal/auth"
	"github.com/daneofmanythings/blog_aggregator/internal/database"
	"github.com/daneofmanythings/blog_aggregator/internal/handlers/utils"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (m *Repository) MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKeyToken(r.Header)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error()+" (retrieving api_key)")
			return
		}

		user, err := m.App.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error()+" (retrieving user from db)")
			return
		}

		handler(w, r, user)
	})
}
