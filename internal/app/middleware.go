package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/notHim0/kanban/internal/types"
	"github.com/notHim0/kanban/internal/utils"
	"github.com/xeipuuv/gojsonschema"
)

// func (app *App) LogginMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
// 		log.Printf("from middleware %s %s %s", r.RemoteAddr, r.Method, r.URL)

// 		next.ServeHTTP(w, r)
// 	})
// }

//Validates input
func (app *App) ValidateMiddleware(schema string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			//to store the json object
			var body map[string] interface{}
			bodyBytes, err := io.ReadAll(r.Body)

			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
				return
			}
			err = json.Unmarshal(bodyBytes, &body)

			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
				return
			}

			//load the required schema to validate
			schemaLoader := gojsonschema.NewStringLoader(schema)
			//load the json object 
			documentLoader := gojsonschema.NewGoLoader(body)

			//validate according to schema
			result, err := gojsonschema.Validate(schemaLoader, documentLoader)

			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Error validating JSON")
				return
			}

			if !result.Valid(){
				var errs []string
				for _, err := range result.Errors(){
					errs=append(errs, err.String())
				}

				utils.RespondWithError(w, http.StatusBadGateway, strings.Join(errs, ", "))
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			next.ServeHTTP(w, r)
		})
	}
}

//check auth for a user
func (app *App) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//extract the token from header
		var authHeader string = r.Header.Get("Authorization")
		//check if token was provided
		if len(authHeader) == 0 {
			utils.RespondWithError(w, http.StatusUnauthorized, "No token provided")
			return 
		}

		//store the token after sanitizing it
		var tokenString string = strings.TrimPrefix(authHeader, "Bearer ")

		//to extract the user info from token
		var claims *types.Claims = &types.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface {}, error){
			return []byte(app.JWTKEY), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid{
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token signature")
				return
			}
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if !token.Valid {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token that")
			return
		}

		//globally set the user info
		var ctx context.Context = context.WithValue(r.Context(), "claims", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}