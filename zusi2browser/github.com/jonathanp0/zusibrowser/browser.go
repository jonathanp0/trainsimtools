package main

import "fmt"
import "net/http"
import "log"
import "path/filepath"
import "os"
import "regexp"
import "html/template"
import "bufio"
import "strconv"

type ZugHalt struct {
	Name string
	Arrival string
	Departure string
}

type Zug struct {
	Number int
	Category string
	Lok string
	Wagons []string
	Stops []ZugHalt
	Path string
}


var zusiRefTime = "1999-12-30  16:20:59"
var timetables = make([]string, 0, 200);

func timetableList(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
	templates := template.Must(template.New("list").Parse(`{{template "header"}}<p>{{len .}} timetable(s) available</p> <ul>{{range .}}<li><a href="/fpl/?path={{.}}">{{.}}</li>{{end}}</ul>{{template "footer"}}`))
	addCommon(templates)
	
	err := templates.Execute(w, timetables)
	if(err != nil) {
		log.Fatal(err)
	}
}

func trainList(w http.ResponseWriter, r *http.Request){
	
	path:= r.URL.Query().Get("path")
	fullpath := filepath.Join(root, path)
	
	fpl, err := parseFpl(fullpath)

	funcMap := template.FuncMap{
    	"filename": filepath.Base,
	}	
	
	templates := template.Must(template.New("list").Funcs(funcMap).Parse(
		`{{template "header"}}
		<ul>{{range .}}<li><a href="/zug/?path={{.}}">{{filename .}}</a></li>{{end}}</ul>
		{{template "footer"}}` ))
	addCommon(templates)
	
	err = templates.Execute(w, fpl)
	if(err != nil) {
		log.Fatal(err)
	}
}	

func trainDetails(w http.ResponseWriter, r *http.Request){
	
	path:= r.URL.Query().Get("path")
	fullpath := filepath.Join(root, path)
	zug, err := parseZug(fullpath)

	if(err != nil) {
		log.Fatal(err)
	}
	
	zug.Path = path
	
	w.Header().Add("Content-type", "text/html; charset=ISO-8859-1")
	templates := template.Must(template.New("zug").Parse(
		`{{template "header"}}
		<h3>Zug  {{.Category}}{{.Number}}</h3>
		<a href="/log/add/?zug={{.Category}}{{.Number}}&lok={{.Lok}}&path={{.Path}}">Add to trip log</a><br />
		<p><b>Lok: </b>{{.Lok}}</p><p><b>{{len .Wagons}} Waggons:</b><ol>{{range .Wagons}}<li>{{.}}</li>{{end}}</ol></p>
		<b>Fahrplan</b>
		<table border="1"><tr><th>Haltestelle</th><th>Ankunft</th><th>Abfahrt</th></tr>
		{{range .Stops}} 
		<tr><td>{{.Name}}</td><td>{{.Arrival}}</td><td>{{.Departure}}</td></tr>
		{{end}}
		{{template "footer"}}` ))
	addCommon(templates)
	
	err = templates.Execute(w, zug)
	if(err != nil) {
		log.Fatal(err)
	}
}	

func scanFpl(path string, info os.FileInfo, err error) error {
	if (err != nil) {
		return err
	}
	matched, err := regexp.MatchString(".*fpl$", info.Name())
	if(matched) {
		relpath, _ := filepath.Rel(root, path)
		fmt.Printf("%s\n", relpath)
		timetables = append(timetables, relpath)
	}
	return err
}

func parseFpl(path string) ([]string, error) {

	relpath,_ := filepath.Rel(root, filepath.Dir(path))
	file, err := os.Open(path)

	if(err != nil){
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	zug := make([]string, 0, 20)
	line := 0

	for scanner.Scan() {
	    if(line > 1) {
	    	zug = append(zug, filepath.Join(relpath,scanner.Text()))
	    }
	    line = line + 1
	}

	if err := scanner.Err(); err != nil {
	    return nil, err
	}

	return zug, nil
}

func parseZug(path string) (*Zug, error) {

	file, err := os.Open(path)

	if(err != nil){
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	zug := Zug{}
	line := -1
	hashes := 0
	mode := "basic"

	for scanner.Scan() {
		line = line + 1

		//Read fixed position data
		if(line == 1) {
			zug.Number, _ = strconv.Atoi(scanner.Text())
		}
		if(line == 2) {
			zug.Category = scanner.Text()
		}
		if(line == 9) {
			zug.Lok = scanner.Text()
		}

		//Count hashes
	    if(scanner.Text() == "#") {
	    	hashes = hashes + 1
	    }
		
	    //Switch mode
		if(hashes == 6) {
	    	mode = "stop"
			hashes = 0
	    } else if(hashes == 1 && mode == "endstop"){
			mode = "tracklist"
			hashes = 0
		} else if(hashes == 1 && mode == "tracklist"){
			mode = "actions"
			hashes = 0
		} else if(hashes == 1 && mode == "actions"){
			mode = "stop"
			hashes = 0
		} else if(hashes == 1 && mode == "endwagon"){
			mode = "wagon"
			hashes = 0
		} else if(scanner.Text() == "#IF"){
	    	if (mode == "basic") {
				mode = "somethingelse"
			} else if(mode == "stop"){
				mode = "endwagon"
			}
			hashes = 0
	    } else if(mode == "stop") {
	    	halt := ZugHalt{scanner.Text(), "", ""}
			scanner.Scan()
			line = line + 1
	    	if(scanner.Text() != "1899-12-30  00:00:00" && scanner.Text() != "") {
	    		halt.Arrival = scanner.Text()[12:]
	    	}
	    	scanner.Scan()
			line = line + 1
	    	if(scanner.Text() != "1899-12-30  00:00:00") {
	    		halt.Departure = scanner.Text()[12:]
	    	}

	    	zug.Stops = append(zug.Stops, halt)
			
			mode = "endstop"
			hashes = 0
	    } else if(mode == "wagon") {
			zug.Wagons = append(zug.Wagons, scanner.Text())
			mode = "endwagon"
			hashes = 0
		}

	}

	if err := scanner.Err(); err != nil {
	    return nil, err
	}

	return &zug, nil
}