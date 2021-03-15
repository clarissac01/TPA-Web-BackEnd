package main

import (
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/clarissac01/Staem/database"
	_ "github.com/clarissac01/Staem/database"
	"github.com/clarissac01/Staem/graph"
	"github.com/clarissac01/Staem/graph/generated"
	_ "github.com/clarissac01/Staem/graph/model"
	model2 "github.com/clarissac01/Staem/graph/model"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

func main() {
	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver()}))
	srv.AddTransport(transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.Use(extension.Introspection{})

	router.Handle("/", playground.Handler("Staem", "/query"))
	router.Handle("/query", srv)
	// /asset/:id? -> return download file
	router.Get("/assets/{id}", GetFile)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}

func GetFile(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	// query db by id
	// w.write Files.file
	// ntar kalo ga bisa, tambahin http header octet/stream

	if s, err := strconv.ParseInt(id, 10, 64); err == nil {
		file := model2.Files{Id: int(s)}
		database.GetDB().Find(&file)
		writer.Write(file.File)
	}
}
