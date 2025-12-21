package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/JuD4Mo/go_api_web_user/internal/user"
	"github.com/JuD4Mo/go_api_web_user/pkg/bootstrap"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	//Instancia de un router de Gorilla Mux
	router := mux.NewRouter()

	//Cargamos las variables de entorno que están en el archivo .env por medio del package godotenv
	_ = godotenv.Load()

	//Instanciamos un logger propio
	l := bootstrap.InitLogger()

	db, err := bootstrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		l.Fatal("paginator limit default is required")
	}

	//Instancias de las capas: repositorio, servicio y controlador
	userRepo := user.NewRepo(l, db)
	userService := user.NewService(l, userRepo)
	userEnd := user.MakeEndpoints(userService, user.Config{LimitPage: pagLimDef})

	//Por medio del router de Gorilla Mux servimos los endpoints
	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userEnd.Get).Methods("GET")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEnd.Delete).Methods("DELETE")

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	//Se crea una instancia de un servidor
	srv := &http.Server{
		Handler:      router,
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	//Se sirve la aplicación y se le vanta el servidor
	err = srv.ListenAndServe()
	if err != nil {
		l.Fatal(err)
	}

	// port := ":3000"
	// http.HandleFunc("/users", getUsers)
	// http.HandleFunc("/courses", getCourses)

	// //Servir la app y levantar el servidor
	// err := http.ListenAndServe(port, nil)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

// func getUsers(w http.ResponseWriter, r *http.Request) {
// 	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
// 	fmt.Println("got /users")
// 	io.WriteString(w, "user endpoint\n")
// }

// func getCourses(w http.ResponseWriter, r *http.Request) {
// 	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
// 	fmt.Println("got /courses")
// 	io.WriteString(w, "course endpoint\n")
// }
