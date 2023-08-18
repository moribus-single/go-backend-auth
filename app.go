package main

import (
	"app/services"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var db services.DbService
var jwt services.JwtService

type RefreshInput struct {
	Refresh string `json:"refresh"`
	Guid    string `json:"guid"`
}

func main() {
	// load env variables into config
	config := services.LoadConfig()

	// get services
	dbService, err := services.GetDbService(config.DbName, config.DbTableName, config.DbURI)
	if err != nil {
		log.Fatal(err)
	}

	db = dbService
	jwt = services.GetJwtService(config.AccessLifetime, config.RefreshLifetime, config.Secret)

	// registering the routes
	http.HandleFunc("/token", tokenHandler)
	http.HandleFunc("/refresh", refreshHandler)

	// running the server
	port := fmt.Sprintf(":%s", config.Port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("Invalid port")
	}
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	// get query parameter GUID
	guid := r.URL.Query().Get("guid")

	// set content type
	w.Header().Set("Content-Type", "application/json")

	// get db instance for some data
	instance, err := db.Read(guid)
	if err != nil {
		http.Error(w, "unable to find by guid", http.StatusBadRequest)
		return
	}

	// generate jwt tokens
	access, refresh, err := updateTokenPair(w, guid, instance.TokenCounter)
	if err != nil {
		return
	}

	// prepare response data
	w.WriteHeader(http.StatusCreated)
	data := make(map[string]string)
	data["access"] = access
	data["refresh"] = refresh

	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "unable to stringify token pair struct", http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	// read the data from request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read the data from request body", http.StatusInternalServerError)
		return
	}

	// parse the data
	var input RefreshInput
	json.Unmarshal(body, &input)

	// get db instance
	instance, err := db.Read(input.Guid)
	if err != nil {
		http.Error(w, "unable to find by guid", http.StatusInternalServerError)
		return
	}

	// ensure refresh is valid
	err = bcrypt.CompareHashAndPassword([]byte(instance.Refresh), []byte(input.Refresh))
	if err != nil {
		http.Error(w, "invalid refresh", http.StatusInternalServerError)
		return
	}

	// generate jwt tokens
	access, refresh, err := updateTokenPair(w, input.Guid, instance.TokenCounter)
	if err != nil {
		return
	}

	// prepare response data
	w.WriteHeader(http.StatusCreated)
	data := make(map[string]string)
	data["access"] = access
	data["refresh"] = refresh

	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "unable to stringify token pair struct", http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func updateTokenPair(w http.ResponseWriter, guid string, tokenCounter int) (string, string, error) {
	// generate jwt tokens
	access, refresh, err := jwt.Generate(guid, tokenCounter+1)
	if err != nil {
		http.Error(w, "unable to generate token pair", http.StatusInternalServerError)
		return "", "", err
	}

	// hash generated refresh token
	hashed, err := bcrypt.GenerateFromPassword([]byte(refresh), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "unable to hash refresh token", http.StatusInternalServerError)
		return "", "", err
	}

	// update db with new hashed refresh token
	err = db.Update(guid, string(hashed))
	if err != nil {
		http.Error(w, "unable to update refresh token", http.StatusInternalServerError)
		return "", "", err
	}

	return access, refresh, nil
}
