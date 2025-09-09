package main

import (
	"database/sql"

	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/justinas/alice"
	"github.com/notHim0/kanban/internal/app"
	"github.com/notHim0/kanban/internal/utils"
)



func main(){
	var err error = godotenv.Load()

	//checking for any env variable loading erors
	if err != nil {
		log.Fatalf("Error loading env file: %s",err.Error())
	}

	var loadErr error
	//loading json schemas for input validation
	userSchema, loadErr := utils.LoadSchema("schemas/user.json")

	//handling loading errors
	if loadErr != nil {
		log.Fatalf("Error loading user schema: %v", loadErr)
	}

	projectSchema, loadErr := utils.LoadSchema("schemas/project.json")

	if loadErr != nil {
		log.Fatalf("Error loadingn project schema: %v", loadErr)
	}

	//establishing connection to the db
	var connectionString string = os.Getenv("POSTGRESQL_URI")

	if len(connectionString) == 0 {
		log.Fatalf("Database uri is not set")
	}

	//loading jwt secret
	var JWTKEY []byte = []byte(os.Getenv("JWTKEY"))

	if len(JWTKEY) == 0 {
		log.Fatalf("JWT key is not set")
	}

	//connecting to db
	DB, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err.Error())
	}

	//checking if the connection is established
	if err := DB.Ping(); err !=nil {
		log.Fatal(err.Error())
	}
	defer DB.Close()

	//creating a app client for request routing
	var app *app.App = &app.App{DB: DB, JWTKEY: JWTKEY}
	
	//creating router
	var router mux.Router = *mux.NewRouter()

	//Middleware which checks schema validation for user route
	var userChain alice.Chain = alice.New(app.ValidateMiddleware((userSchema)))
	router.Handle("/register", userChain.ThenFunc(app.Register)).Methods("POST")
	router.Handle("/login", userChain.ThenFunc(app.Login)).Methods("POST")

	//Middleware chain for login requiring routes(eg: GET, DELETE)
	var projectChain alice.Chain = alice.New(app.JWTMiddleware)
	router.Handle("/projects", projectChain.ThenFunc(app.GetProject)).Methods("GET")
	router.Handle("/projects/{id}", projectChain.ThenFunc(app.DeleteProject)).Methods("DELETE")
	router.Handle("/projects/{id}", projectChain.ThenFunc(app.GetProject)).Methods("GET")

	//Middleware chain for login and schema validation requiring routes
	var projectChainWithValidation alice.Chain = projectChain.Append(app.ValidateMiddleware((projectSchema)))
	router.Handle("/projects", projectChainWithValidation.ThenFunc(app.CreateProject)).Methods("POST")
	router.Handle("/projects/{id}", projectChainWithValidation.ThenFunc(app.UpdateProject)).Methods("PUT")
	
	log.Fatal(http.ListenAndServe(":5000", &router))
}





// func createTable(app App){
// 	var query string = 
// 	`CREATE TABLE IF NOT EXISTS "user"(
// 	id SERIAL PRIMARY KEY,
// 	name VARCHAR(100) NOT NULL,
// 	password VARCHAR(200) NOT NULL,
// 	created timestamp DEFAULT NOW()
// 	)` 

// 	_, err := app.DB.Exec(query)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }


// func register(w http.ResponseWriter, r *http.Request) {
// 	var vars map[string]string = mux.Vars(r)
// 	var id string = vars["id"]
// 	log.Println(id)
// 	w.Header().Set("Content-type", "application/json")
// 	json.NewEncoder(w).Encode(RouteResponse{Message: "Hello from register", ID: id})
// }