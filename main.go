package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"github.com/LiamPimlott/spaces/users"
)

var (
	dbHost string
	dbPort string
	dbUser string
	dbPass string
	dbName string
	secret string
)

func init() {
	viper.SetConfigFile(`config.json`)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	dbHost = viper.GetString(`database.host`)
	dbPort = viper.GetString(`database.port`)
	dbUser = viper.GetString(`database.user`)
	dbPass = viper.GetString(`database.pass`)
	dbName = viper.GetString(`database.name`)

	secret = viper.GetString(`jwt.secret`)
}

func main() {
	// db
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db is up.")

	// repos
	usersRepository := users.NewMysqlUsersRepository(db)

	// services
	usersService := users.NewUsersService(usersRepository, secret)

	// handlers
	createUserHandler := users.NewCreateUserHandler(usersService)
	getUserHandler := users.NewGetUserByIdHandler(usersService)

	// routing
	r := mux.NewRouter()

	// users
	r.Handle("/users", createUserHandler).Methods("POST")
	r.Handle("/users/{id}", getUserHandler).Methods("GET")

	// serve static assest like images, css from the /static/{file} route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/build/static"))))

	// root route will serve the built react app.
	r.Handle("/", http.FileServer(http.Dir("./frontend/build")))

	// start server
	log.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}
