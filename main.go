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

	auth "github.com/LiamPimlott/spaces/middleware"
	"github.com/LiamPimlott/spaces/spaces"
	"github.com/LiamPimlott/spaces/users"
)

var (
	dbHost string
	dbPort string
	dbUser string
	dbPass string
	dbName string
	secret string
	port   string
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
	port = viper.GetString(`server.port`)
}

func main() {
	////////
	// DB //
	////////

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db connected")

	////////////
	// ROUTER //
	////////////

	r := mux.NewRouter()

	///////////
	// USERS //
	///////////

	// repo & service
	usersRepository := users.NewMysqlUsersRepository(db)
	usersService := users.NewUsersService(usersRepository, secret)

	// handlers
	createUserHandler := users.NewCreateUserHandler(usersService)
	loginUserHandler := users.NewLoginHandler(usersService)
	getUserByIDHandler := users.NewGetUserByIDHandler(usersService)

	// routes
	r.Handle("/users", createUserHandler).Methods("POST")
	r.Handle("/users/{id}", auth.Required(getUserByIDHandler, secret)).Methods("GET")
	r.Handle("/users/login", loginUserHandler).Methods("POST")

	////////////
	// SPACES //
	////////////

	// repo & service
	spacesRepository := spaces.NewMysqlSpacesRepository(db)
	spacesService := spaces.NewSpacesService(spacesRepository)

	// handlers
	createSpaceHandler := spaces.NewCreateSpaceHandler(spacesService)
	getSpaceByIDHandler := spaces.NewGetSpaceByIDHandler(spacesService)

	// routes
	r.Handle("/spaces", auth.Required(createSpaceHandler, secret)).Methods("POST")
	r.Handle("/spaces/{id}", auth.Optional(getSpaceByIDHandler, secret)).Methods("GET")

	////////////
	// STATIC //
	////////////

	// serve static assest like images, css from the /static/{file} route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/build/static"))))

	// root route will serve the built react app
	r.Handle("/", http.FileServer(http.Dir("./frontend/build")))

	// start server
	log.Printf("listening on port %s\n", port)
	http.ListenAndServe(port, handlers.LoggingHandler(os.Stdout, r))
}
