package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gitlab.com/daneofmanythings/chirpy/internal/auth"
	"gitlab.com/daneofmanythings/chirpy/internal/webhooks"
	"gitlab.com/daneofmanythings/chirpy/pkg/models"
)

func (m *Repository) ApiHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (m *Repository) ApiPostChirpsHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	subject, err := auth.ValidateUserAccess(r, m.App.JWTSECRET)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	chirp, err := m.App.DB.CreateChirp(censorChirp(params.Body), userID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusCreated, chirp)
	}
}

func (m *Repository) ApiGetChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := m.App.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get chirps from database")
		return
	}
	if s := r.URL.Query().Get("author_id"); s != "" {
		authorID, err := strconv.Atoi(s)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "couldn't parse author_id")
		}
		chirps = filterByAuthorID(&chirps, authorID)
	}
	if s := r.URL.Query().Get("sort"); s == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func (m *Repository) ApiGetChirpsByChirpIDHandler(w http.ResponseWriter, r *http.Request) {
	numStr := chi.URLParam(r, "chirpID")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't parse author_id")
		return
	}

	chirps, err := m.App.DB.GetChirps()
	if err != nil {
		log.Printf("Error accessing database: %s\n", err)
		w.WriteHeader(500)
		return
	}
	if num > len(chirps) || num < 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, chirps[num-1])
}

func (m *Repository) ApiDeleteChirpsByIDHandler(w http.ResponseWriter, r *http.Request) {
	// getting the requested chirp to delete
	chirpIDString := chi.URLParam(r, "chirpID")
	chirpIDNum, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to parse chirpID")
	}

	subject, err := auth.ValidateUserAccess(r, m.App.JWTSECRET)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	chirp, err := m.App.DB.GetChirpByChirpID(chirpIDNum)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// only allow deletion of the chirp if it is the author requesting it
	if chirp.AuthorID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	err = m.App.DB.DeleteChirpByID(chirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (m *Repository) ApiPostUsersHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	payload, err := m.App.DB.CreateUser(params.Password, params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusCreated, payload)
}

func (m *Repository) ApiPostLoginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	// find the user
	user, err := m.App.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// validate the request passworn
	if err := auth.AuthorizePasswordHash(user.Password, params.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Generate the tokens
	accessToken, err := auth.GenerateToken(auth.AccessJWTIssuer, m.App.JWTSECRET, auth.AccessTokenDuration, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	refreshToken, err := auth.GenerateToken(auth.RefreshJWTIssuer, m.App.JWTSECRET, auth.RefreshTokenDuration, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	// creating the payload
	payload := models.UserResource{
		ID:           user.ID,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, payload)
}

func (m *Repository) ApiPutUsersHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}

	subject, err := auth.ValidateUserAccess(r, m.App.JWTSECRET)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	passwordHashBytes, err := auth.GeneratePasswordHash(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password")
	}
	payload, err := m.App.DB.UpdateUser(userID, params.Email, string(passwordHashBytes), params.IsChirpyRed)
	if err != nil {
		log.Fatalf("user was not able to be updated: %s", err)
	}
	respondWithJSON(w, http.StatusOK, payload)
}

func (m *Repository) ApiRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
	}

	// requires a refresh token
	subject, err := auth.ValidateJWT(tokenString, m.App.JWTSECRET, auth.RefreshJWTIssuer)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	if m.App.DB.IsTokenRevoked(tokenString) {
		respondWithError(w, http.StatusUnauthorized, "token is revoked")
		return
	}
	id, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
	}

	newAccessToken, err := auth.GenerateToken(auth.AccessJWTIssuer, m.App.JWTSECRET, auth.AccessTokenDuration, id)
	if err != nil {
		log.Fatalf("unable to generate token: %s", err)
	}
	payload := struct {
		Token string `json:"token"`
	}{
		Token: newAccessToken,
	}
	log.Println("Token Refreshed")
	respondWithJSON(w, http.StatusOK, payload)
}

func (m *Repository) ApiRevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
	}

	// requires a refresh token
	_, err = auth.ValidateJWT(tokenString, m.App.JWTSECRET, auth.RefreshJWTIssuer)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	if m.App.DB.IsTokenRevoked(tokenString) {
		respondWithError(w, http.StatusUnauthorized, "Token is already revoked")
		return
	}
	m.App.DB.RevokeToken(tokenString)
	w.WriteHeader(http.StatusOK)
}

func (m *Repository) ApiPostPolkaWebhooksHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Data struct {
			UserID int `json:"user_id"`
		} `json:"data"`
		Event string `json:"event"`
	}

	apiKey, err := auth.GetApiKeyToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if apiKey != m.App.APIKEYPOLKA {
		log.Println(apiKey)
		log.Println(m.App.APIKEYPOLKA)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	if params.Event != webhooks.UserUpgraded {
		// placeholder for now
		w.WriteHeader(http.StatusOK)
		return
	}
	err = m.App.DB.ActivateChirpyRed(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to activate chirpy red")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (m *Repository) ApiResetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	m.App.ResetHits()
	w.Write([]byte("Metrics have been reset"))
}
