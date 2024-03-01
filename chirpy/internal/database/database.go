package database

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"gitlab.com/daneofmanythings/chirpy/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotFound error = errors.New("user not found in database")

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps      map[int]models.Chirp `json:"chirps"`
	Users       map[int]models.User  `json:"users"`
	Revokations map[string]time.Time `json:"revokations"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}

	err := db.ensureDB()
	if err != nil {
		log.Println("ensureDB threw an error")
		return nil, err
	}

	return db, nil
}

func newDBStructure() *DBStructure {
	dbs := &DBStructure{}
	dbs.Chirps = make(map[int]models.Chirp)
	dbs.Users = make(map[int]models.User)
	dbs.Revokations = make(map[string]time.Time)
	return dbs
}

func (db *DB) ensureDB() error {
	log.Println("ensureDB fired...")
	_, err := os.Lstat(db.path)
	if os.IsNotExist(err) {
		log.Println("ensureDB is writing a new database")
		dat, err := json.Marshal(newDBStructure())
		if err != nil {
			return err
		}

		db.mu.Lock()
		os.WriteFile(db.path, dat, 0644)
		db.mu.Unlock()
	} else {
		return err
	}
	return nil
}

func (db *DB) CreateChirp(body string, author_id int) (models.Chirp, error) {
	if len(body) > 140 || len(body) == 0 {
		return models.Chirp{}, fmt.Errorf("invalid chirp length: %d", len(body))
	}

	dbs, err := db.loadDB()
	if err != nil {
		return models.Chirp{}, err
	}
	chirp := models.Chirp{
		ID:       len(dbs.Chirps) + 1,
		Body:     body,
		AuthorID: author_id,
	}
	dbs.Chirps[chirp.ID] = chirp

	err = db.writeDB(dbs)
	if err != nil {
		return models.Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]models.Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := []models.Chirp{}
	for _, c := range dbs.Chirps {
		chirps = append(chirps, c)
	}
	sort.Slice(chirps, func(i, j int) bool { return chirps[i].ID < chirps[j].ID })
	return chirps, nil
}

func (db *DB) GetChirpByChirpID(chirpID int) (models.Chirp, error) {
	chirps, err := db.GetChirps()
	if err != nil {
		return models.Chirp{}, err
	}
	for _, chirp := range chirps {
		if chirp.ID == chirpID {
			return chirp, nil
		}
	}
	return models.Chirp{}, errors.New("chirp could not be found")
}

func (db *DB) DeleteChirpByID(chirpID int) error {
	dbs, err := db.loadDB()
	if err != nil {
		return err
	}
	_, ok := dbs.Chirps[chirpID]
	if !ok {
		return errors.New("could not find chirp in database")
	}
	delete(dbs.Chirps, chirpID)

	err = db.writeDB(dbs)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) CreateUser(password, email string) (models.SanitizedUser, error) {
	log.Println("creating user")
	// TODO: validate the email address

	// TODO: allow these to fail gracefully
	dbs, err := db.loadDB()
	if err != nil {
		log.Fatalf("error loading database: %s", err)
	}

	for _, user := range dbs.Users {
		if email == user.Email {
			return models.SanitizedUser{}, errors.New("that email address already has an account")
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.SanitizedUser{}, err
	}

	user := models.User{
		ID:          len(dbs.Users) + 1,
		Email:       email,
		Password:    string(passwordHash),
		IsChirpyRed: false,
	}
	dbs.Users[user.ID] = user

	err = db.writeDB(dbs)
	if err != nil {
		log.Fatalf("error writing database: %s", err)
	}

	return models.SanitizedUser{
		Email:       user.Email,
		ID:          user.ID,
		IsChirpyRed: false,
	}, nil
}

func (db *DB) GetUsers() ([]models.User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	users := []models.User{}
	for _, u := range dbs.Users {
		users = append(users, u)
	}
	sort.Slice(users, func(i, j int) bool { return users[i].ID < users[j].ID })
	return users, nil
}

func (db *DB) GetUserByID(id int) (models.User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return models.User{}, err
	}
	if user, ok := dbs.Users[id]; !ok {
		return models.User{}, ErrUserNotFound
	} else {
		return user, nil
	}
}

func (db *DB) GetUserByEmail(email string) (models.User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return models.User{}, err
	}
	for _, user := range dbs.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return models.User{}, ErrUserNotFound
}

func (db *DB) UpdateUser(userID int, userEmail, passwordHash string, isChirpyRed bool) (models.SanitizedUser, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return models.SanitizedUser{}, err
	}

	// WARN: this is 2 read operations
	user, err := db.GetUserByID(userID)
	if err != nil {
		return models.SanitizedUser{}, err
	}

	user.Email = userEmail
	user.Password = passwordHash

	dbs.Users[userID] = user
	err = db.writeDB(dbs)
	if err != nil {
		return models.SanitizedUser{}, err
	}

	return models.SanitizedUser{
		Email: userEmail,
		ID:    userID,
	}, nil
}

func (db *DB) ActivateChirpyRed(userID int) error {
	dbs, err := db.loadDB()
	if err != nil {
		return err
	}
	user, ok := dbs.Users[userID]
	if !ok {
		return errors.New("couldn't find user")
	}

	user.IsChirpyRed = true
	dbs.Users[userID] = user

	err = db.writeDB(dbs)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) RevokeToken(token string) error {
	dbs, err := db.loadDB()
	if err != nil {
		return err
	}
	dbs.Revokations[token] = time.Now().UTC()
	err = db.writeDB(dbs)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) IsTokenRevoked(token string) bool {
	dbs, err := db.loadDB()
	if err != nil {
		log.Fatalf("could not load database: %s", err)
	}
	if _, ok := dbs.Revokations[token]; !ok {
		return false
	}
	return true
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	fi, err := os.ReadFile(db.path)
	db.mu.RUnlock()

	if err != nil {
		log.Println("loadDB errored reading the file")
		return DBStructure{}, err
	}

	reader := bytes.NewReader(fi)
	decoder := json.NewDecoder(reader)
	dbs := DBStructure{}
	err = decoder.Decode(&dbs)
	if err != nil {
		log.Println("loadDB errored decoding json")
		return DBStructure{}, err
	}

	return dbs, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	db.mu.Lock()
	os.WriteFile(db.path, dat, 0644)
	db.mu.Unlock()
	return nil
}
