package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

func main() {
	templateString := fmt.Sprintf("package global\n\nfunc GetDefaultDirectory() string {\n\treturn \"{{ .DefaultDirFromEnv }}\"\n}\n")
	var fileTemplate = template.Must(template.New("").Parse(templateString))

	var defaultDirFromEnv = os.Getenv("GITCONVEX_DEFAULT_PATH")
	file, fileErr := os.Create("./global/default_dir_const.go")
	if fileErr != nil {
		log.Fatal(fileErr.Error())
	}

	err := fileTemplate.ExecuteTemplate(file, "", struct {
		DefaultDirFromEnv string
	}{
		DefaultDirFromEnv: defaultDirFromEnv,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}
