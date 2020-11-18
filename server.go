package main

import (
	"flag"
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
	"strconv"
)

const defaultPort = "9001"

var (
	Port int
)

func main() {
	var portErr error
	Port = 0

	logger := global.Logger{}
	logger.Log("Starting Gitconvex server modules", global.StatusInfo)

	// checks if the env_config file is accessible. If not then the EnvConfigFileGenerator will be invoked
	// to generate a default env_config file

	if envError := utils.EnvConfigValidator(); envError == nil {
		logger.Log("Using available env config file", global.StatusInfo)
		envConfig := *utils.EnvConfigFileReader()
		Port, portErr = strconv.Atoi(envConfig.Port)
		if portErr != nil {
			Port = 9001
			portErr = nil
		}
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

	// Static file supplier for hosting the react static assets and scripts
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./build/static/"))))

	// Route for serving the webpage logo from the reach build bundle
	router.PathPrefix("/gitconvex.png").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./build/gitconvex.png")
	})

	// A default fallback route for handling all routes with '/' prefix.
	// For making it compatible with react router
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./build/index.html")
	})

	// Checking and Assigning port received from the command line ( --port args )
	argFlag := flag.String("port", "", "To define the port dynamically while starting gitconvex")
	flag.Parse()

	if *argFlag != "" {
		logger.Log(fmt.Sprintf("Setting port received from the command line -> %s", *argFlag), global.StatusInfo)
		Port, portErr = strconv.Atoi(*argFlag)
		if portErr != nil {
			Port = 9001
		}
	}

	if Port > 0 {
		logger.Log(fmt.Sprintf("Gitconvex started on  http://localhost:%v", Port), global.StatusInfo)
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(Port), cors.Default().Handler(router)))
	} else {
		logger.Log(fmt.Sprintf("Gitconvex started on  http://localhost:%v", defaultPort), global.StatusInfo)
		log.Fatal(http.ListenAndServe(":"+defaultPort, cors.Default().Handler(router)))
	}
}
