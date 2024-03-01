package models

type Chirp struct {
	Body     string `json:"body"`
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
}

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type UserResource struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type UserUpdate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SanitizedUser struct {
	Email       string `json:"email"`
	ID          int    `json:"id"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type Chirps struct {
	Chirps []Chirp `json:"chirps"`
}

type ReturnErrorJSON struct {
	Error string `json:"error"`
}
