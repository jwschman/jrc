/*
Jailhouse Roll Call
A simple Go application to retrieve and display information about jails from a TrueNAS server.
Makes an API call to fetch jail details, processes the returned data, and serves it in an HTML table on a local web server.
Created by John Schmanski
Date: October 25, 2024
*/

package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	APIURL string `json:"api_url"`
	APIKey string `json:"api_key"`
}

var config Config
var tmpl *template.Template

func init() {

	// Load configuration from config.json
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	// Pre-parse index.html template once at startup
	tmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("Failed to parse template file: %v", err)
	}
}

type JailAttributes struct {
	Id          int    `json:"jid"`
	Name        string `json:"id"`
	Status      string `json:"state"`
	Release     string `json:"release"`
	Address     string `json:"ip4_addr"`
	LastStarted string `json:"last_started"`
}

func main() {

	// Serve static files (css only at the moment)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route handler
	http.HandleFunc("/", showJails)

	// Start the http server on port 3000
	log.Print("Listening on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}

// getJails makes an HTTP GET request to the TrueNAS API to retrieve the list of jails.
// It unmarshals the JSON response into a slice of JailAttributes structs
// and removes any "vnet0|" prefix from the Address field.
func getJails() ([]JailAttributes, error) {

	// Create a new GET request to the TrueNAS API using the API URL in config.json
	request, err := http.NewRequest("GET", config.APIURL, nil)
	if err != nil {
		return nil, err
	}

	// Add the API key to the request header for authentication
	request.Header.Add("Authorization", "Bearer "+config.APIKey)

	// Send the request and capture the response
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body and store it as 'body'
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data from the response 'body' into a slice of JailAttributes structs 'jails'
	var jails []JailAttributes
	if err := json.Unmarshal(body, &jails); err != nil {
		return nil, err
	}

	// Remove the "vnet0|" prefix from the Address field if it exists
	for i, jail := range jails {
		jails[i].Address = strings.TrimPrefix(jail.Address, "vnet0|")
	}

	return jails, nil
}

// showJails is an HTTP handler that retrieves jail data, renders it in an HTML template, and sends it as a response
func showJails(wr http.ResponseWriter, r *http.Request) {

	// Call getJails to retrieve the list of jails from the TrueNAS API
	jails, err := getJails()
	if err != nil {
		http.Error(wr, "Failed to retrieve jails", http.StatusInternalServerError)
		return
	}

	// Render the pre-parsed template with the jail data and write the rendered HTML to the response
	if err := tmpl.Execute(wr, jails); err != nil {
		http.Error(wr, "Failed to render template", http.StatusInternalServerError)
	}
}
