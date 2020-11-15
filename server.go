package main

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph"
	"github.com/neel1996/gitconvex-server/graph/generated"
	"github.com/neel1996/gitconvex-server/utils"
	"github.com/rs/cors"
	"log"
	"net/http"
)

const defaultPort = "9001"

var (
	Port string
)

func main() {
	logger := global.Logger{}
	logger.Log("Starting Gitconvex server modules", global.StatusInfo)

	// checks if the env_config file is accessible. If not then the EnvConfigFileGenerator will be invoked
	// to generate a default env_config file

	if envError := utils.EnvConfigValidator(); envError == nil {
		logger.Log("Using available env config file", global.StatusInfo)
		envConfig := *utils.EnvConfigFileReader()
		Port = envConfig.Port
	} else {
		logger.Log("No env config file is present. Falling back to default config data", global.StatusWarning)
		envGeneratorError := utils.EnvConfigFileGenerator()
		if envGeneratorError != nil {
			panic(envGeneratorError)
		}
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/gitconvexapi", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	router := mux.NewRouter()

	// http route handler for provisioning graphql playground UI when the API is directly opened from the browser

	router.Path("/gitconvexapi/graph").Handler(playground.Handler("GraphQL", "/query"))
	router.Handle("/query", srv)
	router.Handle("/gitconvexapi", srv)

	// Static file supplier for hosting the react application

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./build/")))

	if Port != "" && len(Port) > 0 {
		logger.Log(fmt.Sprintf("Gitconvex started on  http://localhost:%v", Port), global.StatusInfo)
		log.Fatal(http.ListenAndServe(":"+Port, cors.Default().Handler(router)))
	} else {
		logger.Log(fmt.Sprintf("Gitconvex started on  http://localhost:%v", defaultPort), global.StatusInfo)
		log.Fatal(http.ListenAndServe(":"+defaultPort, cors.Default().Handler(router)))
	}
}
