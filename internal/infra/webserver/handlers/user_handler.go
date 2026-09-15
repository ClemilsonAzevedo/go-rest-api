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
	userDB database.UserInterface
}

func NewUserHandler(db database.UserInterface) *UserHandler {
	return &UserHandler{
		userDB: db,
	}
}

// GetJWT GoDoc
// @Summary     Get a user JWT
// @Description Get a user JWT
// @Tags        users
// @Accept      JSON
// @Produce     JSON
// @Param       request     body        dto.GetJWTInput     true    "user credentials"
// @Success     200         {object}    dto.GetJWTOutput
// @Failure     404         Error
// @Failure     500         {object}    Error
// @Router      /users/generate_token   [post]
func (u *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("jwtExpiresIn").(int)

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

	_, tokenStr, _ := jwt.Encode(
		map[string]any{
			"sub": usr.ID,
			"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
		},
	)

	accessToken := dto.GetJWTOutput{AccessToken: tokenStr}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

// CreateUser create user GoDoc
// @Summary     Create user
// @Description Create user
// @Tags        users
// @Accept      JSON
// @Produce     JSON
// @Param       request     body        dto.CreateUserInput     true    "user request"
// @Success     201
// @Failure     500         {object}    Error
// @Router      /users      [post]
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
