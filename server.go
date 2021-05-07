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

var logger global.Logger

func main() {
	versionFlag := flag.Bool("version", false, "To get the current version of gitconvex")
	portFlag := flag.String("port", "", "To define the port dynamically while starting gitconvex")
	defaultPath, dirCreateErr := utils.DefaultDirSetup()
	if dirCreateErr != nil {
		logger.Log(dirCreateErr.Error(), global.StatusError)
	}
	flag.Parse()

	if *versionFlag == true {
		fmt.Println(global.GetCurrentVersion())
		os.Exit(0)
	}

	var portErr error
	Port = 0

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

	execName, execErr := os.Executable()
	var (
		uiStaticPath string
		uiBuildPath  string
	)
	if execErr != nil {
		logger.Log(fmt.Sprintf("Unable to find exec path - %s", execErr.Error()), global.StatusError)
	} else {
		currentDir := filepath.Dir(execName)
		uiBuildPath = fmt.Sprintf("%s/gitconvex-ui", currentDir)
		uiStaticPath = fmt.Sprintf("%s/static", uiBuildPath)
	}

	// Setting the server to use the UI bundle from the current directory if the bundle is unavailable in the executable directory
	// Resolved UI file server issue when running docker containers
	_, uiOpenErr := os.Open(uiBuildPath)
	if uiOpenErr != nil {
		logger.Log(uiOpenErr.Error(), global.StatusError)
		cwd, _ := os.Getwd()
		if cwd != "" {
			logger.Log("Trying to use the UI bundle from the current directory -> "+cwd, global.StatusInfo)
			uiBuildPath = fmt.Sprintf("%s/gitconvex-ui", cwd)

			if _, localUIOpenErr := os.Open(uiBuildPath); localUIOpenErr != nil {
				logger.Log("Unable to find local ui directory. Falling back to the UI directory lookup in basedir -> "+defaultPath, global.StatusWarning)
				uiBuildPath = fmt.Sprintf("%s/gitconvex-ui", defaultPath)
			}
			uiStaticPath = fmt.Sprintf("%s/static", uiBuildPath)
		}
	}

	// Static file supplier for hosting the react static assets and scripts
	logger.Log(fmt.Sprintf("Serving static files from -> %s", uiStaticPath), global.StatusInfo)
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir(uiStaticPath))))

	// Route for serving the webpage logo from the reach build bundle
	router.PathPrefix("/gitconvex.png").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log(fmt.Sprintf("Serving logo from directory -> %s", uiBuildPath), global.StatusInfo)
		http.ServeFile(w, r, uiBuildPath+"/gitconvex.png")
	})

	// A default fallback route for handling all routes with '/' prefix.
	// For making it compatible with react router
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log(fmt.Sprintf("Serving UI from directory -> %s", uiBuildPath), global.StatusInfo)
		http.ServeFile(w, r, uiBuildPath+"/index.html")
	})

	// Checking and Assigning port received from the command line ( --port args )
	flag.Parse()

	if *portFlag != "" {
		logger.Log(fmt.Sprintf("Setting port received from the command line -> %s", *portFlag), global.StatusInfo)
		Port, portErr = strconv.Atoi(*portFlag)
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
