package main

import "net/http"
import "html/template"
import "log"
import "time"
import "encoding/json"
import "io/ioutil"

var storage = "journeylog.json"

type Journey struct {
	CompleteTime time.Time
	Zug string
	Lok string
	Path string
	Comment string
}

type JourneyLog []Journey

func journeyList(w http.ResponseWriter, r *http.Request){

	journeys := readJourneys()
	
	templates := template.Must(template.New("list").Parse(
		`{{template "header"}}
		<table border="1"><tr><th>Zeit</th><th>Zug</th><th>Lok</th><th>Bemerkung</th><tr>
		{{range .}}
		<tr><td>{{.CompleteTime.Format "02.01.2006 15:04" }}</td>
		<td><a href="/zug/?path={{.Path}}">{{.Zug}}</a></td>
		<td>{{.Lok}}</td>
		<td>{{.Comment}}</td></tr>
		{{end}}
		{{template "footer"}}` ))
	addCommon(templates)
	
	err := templates.Execute(w, journeys)
	if(err != nil) {
		log.Fatal(err)
	}
}

func readJourneys() JourneyLog {
	contents, err := ioutil.ReadFile(storage)

	input := make(JourneyLog, 0, 0)

	if(err != nil) {
		log.Print("Error reading journey log:", err)
		return input
	}

	err = json.Unmarshal(contents, &input)

	if(err != nil) {
		log.Fatal(err)
		return input
	}

	return input
}

func writeJourneys(logs JourneyLog) {
	jsontest, _ := json.Marshal(logs)

	ioutil.WriteFile(storage, jsontest, 0)
}

func journeyFormAdd(w http.ResponseWriter, r *http.Request){
	
	if(r.Method == "POST"){
		r.ParseForm()

		journey := Journey{time.Now(), r.PostFormValue("zug"), r.PostFormValue("lok"), r.PostFormValue("path"), r.PostFormValue("comment")}

		stored := readJourneys()
		stored = append(stored, journey)
		writeJourneys(stored)

		http.Redirect(w, r, "/log/", 303)
		return
	}

	journey := Journey{time.Now(), r.URL.Query().Get("zug"), r.URL.Query().Get("lok"), r.URL.Query().Get("path"), ""}

	templates := template.Must(template.New("list").Parse(
		`{{template "header"}}
		<form action="/log/add/" method="post">
		<fieldset><h3>Add Trip</h3><table>
	    <tr><td>Zug Number</td><td><input type="text" name="zug" value="{{.Zug}}" size="8"></td></tr>
		<tr><td>Datei</td><td><input type="text" name="path" value="{{.Path}}" size="50"></td></tr>
	    <tr><td>Lok</td><td><input type="text" name="lok" value="{{.Lok}}" size="50"></td></tr>
	    <tr><td>Bemerkung</td><td><input type="text" name="comment" value="{{.Comment}}" size="50"></td></tr>
	    </table>
	    <input type="submit" name="submit" value="Add">
	    </fieldset>
	    </form>
		{{template "footer"}}` ))
	addCommon(templates)
	
	err := templates.Execute(w, journey)
	if(err != nil) {
		log.Fatal(err)
	}
}