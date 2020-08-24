package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"flag"

	"github.com/julienschmidt/httprouter"
)

func main() {

	fileLocation := flag.String("loc", "", "The location to serve static files from")

	portNumber := flag.String("port", ":8080", "The port to serve on, in form :#####")

	showFiles := flag.Bool("sf", false, "Set to true to show the file tree")

	

	flag.Parse()

	fmt.Println("Yo")
	fmt.Println(*showFiles)
	if *fileLocation == "" {
		panic("Please provide a file location")
	}
	staticRouter := httprouter.New()
	//Server all requests to static from direct URL, no static pre-fix
	//Since the static calls are done only after the activeRouter has not matched
	//Any special events with a page can be handled first
	if *showFiles {
		staticRouter.ServeFiles("/*filepath", http.Dir(*fileLocation))
	} else {
		staticRouter.ServeFiles("/*filepath", neuteredFileSystem{http.Dir(*fileLocation)})
	}

	server := &http.Server{
		Addr:    *portNumber,
		Handler: staticRouter,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	//Modified code, no need to have .html at the end of our html static filepaths
	//If our filepath is not a directory and is looking for a file
	if filepath.Base(path) != "\\" && filepath.Base(path) != "" {
		if filepath.Ext(path) == "" {
			//Append the file path with a .html
			path = fmt.Sprintf("%s.html", path)
		}
	}

	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := nfs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}
