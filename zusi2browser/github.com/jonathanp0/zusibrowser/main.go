package main

import "fmt"
import "net/http"
import "log"
import "path/filepath"
import "html/template"
import "flag"

var root = "C:\\Program Files (x86)\\Zusi"
var showHelp bool

func init() {

	flag.BoolVar(&showHelp, "h", false, "Show the help message")
	flag.StringVar(&root, "zusi", "C:\\Program Files (x86)\\Zusi", "Zusi 2 Install Directory")

}

func addCommon(templates *template.Template) {
	templates.Parse(`{{define "header"}}
		<html><head><title>Zusi 2 Info Browser</title>
		<style>
		body {
			font-family: Arial, sans-serif;
			font-size: 14px;
			background-color: #fff;
			line-height: 1.3em;
		}
		</style>
		</head>
		<body><h1>Zusi 2 Info Browser</h1>
		<a href="/">Home</a> | <a href="/log/">Trip Log</a> | <a href="/log/add/">Add Trip</a><hr />{{end}}`)
	templates.Parse(`{{define "footer"}}</body></html>{{end}}`)
}


func main() {

	flag.Parse()

	if(showHelp == true) {
		flag.PrintDefaults();
		return
	}

	strecken :=  filepath.Join(root, "Strecken")

	fmt.Printf("* Scanning for .fpl files in %s...\n", strecken)
    err := filepath.Walk(strecken, scanFpl)
	if(err != nil) {
		log.Fatal("Error scanning for timetables: \n", err.Error())
	}
	
	fmt.Printf("* Scanning complete. Go to http://localhost:8888/\n")
	fmt.Printf("* Press Ctrl+C to exit\n")
	
	http.HandleFunc("/", timetableList)
	http.HandleFunc("/fpl/", trainList)
	http.HandleFunc("/zug/", trainDetails)

	http.HandleFunc("/log/", journeyList)
	http.HandleFunc("/log/add/", journeyFormAdd)

	log.Fatal(http.ListenAndServe(":8888", nil))
}
