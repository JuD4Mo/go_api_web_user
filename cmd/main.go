package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JuD4Mo/go_api_web_user/internal/user"
	"github.com/JuD4Mo/go_api_web_user/pkg/bootstrap"
	"github.com/JuD4Mo/go_api_web_user/pkg/handler"
	"github.com/joho/godotenv"
)

func main() {

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

	ctx := context.Background()

	//Instancias de las capas: repositorio, servicio y controlador
	userRepo := user.NewRepo(l, db)
	userService := user.NewService(l, userRepo)
	h := handler.NewUserHTTPServer(ctx, user.MakeEndpoints(userService, user.Config{LimitPage: pagLimDef}))

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	//Se crea una instancia de un servidor
	srv := &http.Server{
		Handler:      accessControl(h),
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	errCh := make(chan error)
	go func() {
		l.Println("listen in", address)
		errCh <- srv.ListenAndServe()
	}()

	err = <-errCh
	if err != nil {
		log.Fatal(err)
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

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST, PATCH, OPTIONS, DELETE, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
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
