package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/clemilsonazevedo/9-apis/internal/dto"
	"github.com/clemilsonazevedo/9-apis/internal/entity"
	"github.com/clemilsonazevedo/9-apis/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	userDB       database.UserInterface
	jwt          *jwtauth.JWTAuth
	jwtExpiresIn int
}

func NewUserHandler(db database.UserInterface, jwt *jwtauth.JWTAuth, jwtExpiresIn int) *UserHandler {
	return &UserHandler{
		userDB:       db,
		jwt:          jwt,
		jwtExpiresIn: jwtExpiresIn,
	}
}

func (u *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var userInput dto.GetJWTInput
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	usr, err := u.userDB.FindUserByEmail(userInput.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !usr.ValidatePassword(userInput.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, tokenStr, _ := u.jwt.Encode(
		map[string]any{
			"sub": usr.ID,
			"exp": time.Now().Add(time.Second * time.Duration(u.jwtExpiresIn)).Unix(),
		},
	)

	accessToken := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: tokenStr,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	foundUser, _ := u.userDB.FindUserByEmail(user.Email)
	if foundUser != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	usr, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := u.userDB.Create(usr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
