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
	"os"
	"path/filepath"
	"strconv"
)

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

	execName, _ := os.Executable()
	var (
		staticPath string
		buildPath  string
	)
	if execName != "" {
		currentDir := filepath.Dir(execName)
		buildPath = fmt.Sprintf("%s/gitconvex-ui", currentDir)
		staticPath = fmt.Sprintf("%s/static", buildPath)
	} else {
		logger.Log("Unable to serve UI bundle", global.StatusError)
	}

	// Setting the server to use the UI bundle from the current directory if the bundle is unavailable in the executable directory
	// Resolved UI file server issue when running docker containers
	_, uiOpenErr := os.Open(buildPath)
	if uiOpenErr != nil {
		logger.Log(uiOpenErr.Error(), global.StatusError)
		cwd, _ := os.Getwd()
		if cwd != "" {
			logger.Log("Using UI bundle from the current directory -> "+cwd, global.StatusInfo)
			buildPath = fmt.Sprintf("%s/gitconvex-ui", cwd)
		}
	}

	// Static file supplier for hosting the react static assets and scripts
	logger.Log(fmt.Sprintf("Serving static files from -> %s", staticPath), global.StatusInfo)
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir(staticPath))))

	// Route for serving the webpage logo from the reach build bundle
	router.PathPrefix("/gitconvex.png").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log(fmt.Sprintf("Serving logo from directory -> %s", buildPath), global.StatusInfo)
		http.ServeFile(w, r, buildPath+"/gitconvex.png")
	})

	// A default fallback route for handling all routes with '/' prefix.
	// For making it compatible with react router
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log(fmt.Sprintf("Serving UI from directory -> %s", buildPath), global.StatusInfo)
		http.ServeFile(w, r, buildPath+"/index.html")
	})

	// Checking and Assigning port received from the command line ( --port args )
	argFlag := flag.String("port", "", "To define the port dynamically while starting gitconvex")
	flag.Parse()

	if *argFlag != "" {
		logger.Log(fmt.Sprintf("Setting port received from the command line -> %s", *argFlag), global.StatusInfo)
		Port, portErr = strconv.Atoi(*argFlag)
		if portErr != nil {
			Port, _ = strconv.Atoi(global.DefaultPort)
		}
	}

	if Port > 0 {
		logger.Log(fmt.Sprintf("Gitconvex started on  http://localhost:%v", Port), global.StatusInfo)
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(Port), cors.Default().Handler(router)))
	} else {
		logger.Log(fmt.Sprintf("Gitconvex started on  http://localhost:%v", global.DefaultPort), global.StatusInfo)
		log.Fatal(http.ListenAndServe(":"+global.DefaultPort, cors.Default().Handler(router)))
	}
}
