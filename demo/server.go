package demo

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/duycs/demo-go/demo/api/config"
	"github.com/duycs/demo-go/demo/api/controllers"
	"github.com/duycs/demo-go/demo/application/middlewares"
	"github.com/duycs/demo-go/demo/application/services"
	"github.com/duycs/demo-go/demo/entities"
	"github.com/duycs/demo-go/demo/infrastructure/repository"
	"github.com/duycs/demo-go/demo/seed"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

var server = Server{}

func Run() {
	var err error

	// load values from .env into the system
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	// init data and migration
	server.InitializeDatabase(os.Getenv("DB_DRIVER"), config.DB_USER, config.DB_PASSWORD, config.DB_PORT, config.DB_HOST, config.DB_DATABASE)

	server.RegisterServiceAndRouter()

	server.ListenToPort(":8080")
}

func (server *Server) ListenToPort(addr string) {
	var err error

	fmt.Println("Listening to port ", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))

	// listen and serve
	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + strconv.Itoa(config.API_PORT),
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err.Error())
	}
}

func (server *Server) RegisterServiceAndRouter() {
	// inject repositories, services
	taskRepo := repository.NewTaskContext(server.DB)
	taskService := services.NewTaskService(taskRepo)
	userRepo := repository.NewUserContext(server.DB)
	userService := services.NewUserService(userRepo)
	assignmentUseCase := services.NewAssignmentService(userService, taskService)

	// register router
	r := mux.NewRouter()
	n := negroni.New(
		negroni.HandlerFunc(middlewares.Cors),
		negroni.NewLogger(),
	)
	controllers.RegisterTaskHandlers(r, *n, taskService)
	controllers.RegisterUserHandlers(r, *n, userService)
	controllers.RegisterAssignmentHandlers(r, *n, taskService, userService, assignmentUseCase)
	http.Handle("/", r)
	http.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func (server *Server) InitializeDatabase(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error

	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	}
	if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	}
	if Dbdriver == "sqlite3" {
		server.DB, err = gorm.Open(Dbdriver, DbName)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", Dbdriver)
		}
		server.DB.Exec("PRAGMA foreign_keys = ON")
	}

	//database migration
	server.DB.Debug().AutoMigrate(&entities.User{})
	server.DB.Debug().AutoMigrate(&entities.Task{})

	defer server.DB.Close()

	seed.Load(server.DB)
}
