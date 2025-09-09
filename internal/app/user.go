package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/notHim0/kanban/internal/models"
	"github.com/notHim0/kanban/internal/types"
	"github.com/notHim0/kanban/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

//func to handle registration
func (app *App) Register(w http.ResponseWriter, r *http.Request){
	var cred types.Credentials

	var err error = json.NewDecoder(r.Body).Decode(&cred)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return 
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cred.Password), bcrypt.DefaultCost)

	var user models.User
	user.Name = cred.Username
	user.Password = string(hashedPassword)
	user.Id = insertUser(*app, &user)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server error")
	}
	tokenString, err := app.generateToken(cred.Username, user.Id)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Token not generated")
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(types.UserResponse{Username: cred.Username, Token: tokenString, Id: user.Id})
}

func (app *App) Login(w http.ResponseWriter, r *http.Request){
	var cred types.Credentials

	var err error = json.NewDecoder(r.Body).Decode(&cred)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var savedCredentials types.Credentials
	var userId string

	err = app.DB.QueryRow(`SELECT id, name, password FROM "user" WHERE id=$1`, cred.Id).Scan(&userId, &savedCredentials.Username, &savedCredentials.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid username or password")
			return
		}
		log.Fatal(err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return 
	}

	err = bcrypt.CompareHashAndPassword([]byte(savedCredentials.Password), []byte(cred.Password))

	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}
	tokenString, err := app.generateToken(cred.Username, userId)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(types.UserResponse{Username: cred.Username, Token: tokenString, Id: userId})

}

func (app *App) generateToken (username, id string) (string, error){
	expirationTime := time.Now().Add(24*time.Hour)

	claims := &types.Claims{
		Username: username,
		Id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(app.JWTKEY))

	if err!=nil {
		return "", err
	}

	return tokenString, err
}

func insertUser(app App, user *models.User) string {
	var query string =
	`INSERT INTO "user" (name, password)
	VALUES($1, $2) RETURNING id`

	var uk string

	err := app.DB.QueryRow(query, user.Name, user.Password).Scan(&uk)

	if err != nil {
		log.Fatal(err)
	}
	return uk
}
