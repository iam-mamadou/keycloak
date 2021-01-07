package main

import (
	"encoding/json"
	"log"
	"time"

	"context"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

var states = make([]state, 100)

type state struct {
	id          string
	createdDate int64
}

func stateFactory() state {
	return state {
		id:          generateUUID(),
		createdDate: nowEpoc(),
	}
}

func nowEpoc() int64 {
	return time.Now().Unix()
}

func generateUUID() string {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		log.Fatalln("Failed to generate UUID")
	}
	return newUUID.String()
}

func isStateValid(state string) bool {
	// TODO add stateTTL to 120 sec
	for _, s := range states {
		if s.id ==  state && s.createdDate <= nowEpoc() {
			return true
		}
	}
	return false
}

func handleRedirect(config *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := stateFactory()
		states = append(states, state)
		http.Redirect(w, r, config.AuthCodeURL(state.id), http.StatusFound)
	}
}

func handleCallback(ctx context.Context, provider *oidc.Provider, config *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isStateValid(r.URL.Query().Get("state")) {
			http.Error(w, "state did not match", http.StatusBadRequest)
			return
		}

		token, err := config.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(token))
		if err != nil {
			http.Error(w, "Failed to get userinfo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json, err := json.Marshal(userInfo)
		if err != nil {
			http.Error(w, "Failed to marshall user info: "+err.Error(), http.StatusInternalServerError)
		}

		w.Write(json)
	}
}

func main() {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "http://localhost:8080/auth/realms/demo")
	if err != nil {
		log.Fatal(err)
	}

	config := &oauth2.Config{
		ClientID:     "demo",
		ClientSecret: "842ef196-1cc7-4a04-a0f4-f95cc7d36f9b",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "location"},
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8000/auth/callback",
	}

	http.Handle("/", handleRedirect(config))

	http.HandleFunc("/auth/callback", handleCallback(ctx, provider, config))

	log.Println("Server listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
