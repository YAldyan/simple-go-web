package main

import (
	"example-go-web/controller"
	"example-go-web/controller/handle"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type NotFound struct{}

func (n *NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

// Creates a new router
func NewRouter() *httprouter.Router {
	router := httprouter.New()
	notFound := new(NotFound)
	router.NotFound = notFound
	return router
}

// Creates a new router
// func NewRouter() *httprouter.Router {
// 	router := httprouter.New()
// 	router.NotFound = func(http.ResponseWriter, *http.Request) {}
// 	return router
// }
func init() {
	// Assign a user store
	store, err := controller.NewFileUserStore("./data/users.json")
	if err != nil {
		panic(fmt.Errorf("Error creating user store: %s", err))
	}
	controller.GlobalUserStore = store

	// Assign a session store
	sessionStore, err := controller.NewFileSessionStore("./data/sessions.json")
	if err != nil {
		panic(fmt.Errorf("Error creating session store: %s", err))
	}
	controller.SetglobalSessionStore(sessionStore)

	// Assign a sql database
	db, err := controller.NewMySQLDB("root:@tcp(127.0.0.1:3306)/Goes")
	if err != nil {
		panic(err)
	}
	controller.GlobalMySQLDB = db

	// Assign an image store
	controller.GlobalImageStore = controller.NewDBImageStore()
}

func main() {

	router := NewRouter()

	router.Handle("GET", "/", handle.HandleHome)
	router.Handle("GET", "/register", controller.HandleUserNew)
	router.Handle("POST", "/register", controller.HandleUserCreate)
	router.Handle("GET", "/login", controller.HandleSessionNew)
	router.Handle("POST", "/login", controller.HandleSessionCreate)
	router.Handle("GET", "/image/:imageID", controller.HandleImageShow)
	router.Handle("GET", "/user/:userID", controller.HandleUserShow)

	router.ServeFiles(
		"/assets/*filepath",
		http.Dir("assets/"),
	)
	// router.ServeFiles(
	// 	"/assets/*filepath",
	// 	http.Dir(".assets/images/"),
	// )

	secureRouter := NewRouter()
	secureRouter.Handle("GET", "/sign-out", controller.HandleSessionDestroy)
	secureRouter.Handle("GET", "/account", controller.HandleUserEdit)
	secureRouter.Handle("POST", "/account", controller.HandleUserUpdate)
	secureRouter.Handle("GET", "/images/new", controller.HandleImageNew)
	secureRouter.Handle("POST", "/images/new", controller.HandleImageCreate)

	middleware := controller.Middleware{}
	middleware.Add(router)
	middleware.Add(http.HandlerFunc(controller.RequireLogin))
	middleware.Add(secureRouter)

	fmt.Println("starting web server at http://localhost:3000/")
	log.Fatal(http.ListenAndServe(":3000", middleware))
}

/*
func renderingPages() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template.RenderTemplate(w, r, "index/home", nil)
	})

	mux.Handle(
		"/assets/",
		http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))),
	)

	fmt.Println("starting web server at http://localhost:3000/")
	http.ListenAndServe(":3000", mux)
}

func routing() {
*/
/*
	Router adalah library bawaan dari golang untuk
	mengetahui skenario dari tiap-tiap URL yang akan
	dieksekusi, apa response-nya.
*/
/*
	unauthenticatedRouter := NewRouter()
	unauthenticatedRouter.GET("/", controller.HandleHome)
	unauthenticatedRouter.GET("/register", users.HandleUserNew)

	authenticatedRouter := NewRouter()
	authenticatedRouter.GET("/images/new", controller.HandleImageNew)

	middleware := controller.Middleware{}
	middleware.Add(unauthenticatedRouter)
	middleware.Add(http.HandlerFunc(controller.AuthenticateRequest))
	middleware.Add(authenticatedRouter)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", middleware)
}

func registrationUser() {
	router := NewRouter()

	router.Handle("GET", "/", controller.HandleHome)
	router.Handle("GET", "/register", users.HandleUserNew)
	router.Handle("POST", "/register", users.HandleUserCreate)

	router.ServeFiles(
		"/assets/*filepath",
		http.Dir("assets/"),
	)

	middleware := controller.Middleware{}
	middleware.Add(router)
	fmt.Println("Listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", middleware))
}
*/
