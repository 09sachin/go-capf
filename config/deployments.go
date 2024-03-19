package config

import (
	"net/http"
	"os/exec"
	"fmt"
)

var (
	// ClientID     = getEnv("ClientID")
	// ClientSecret = getEnv("ClientSecret")
	ScriptPath   = getEnv("ScriptPath")
)

// Define your expected client ID and client secret
const (
	expectedClientID     = "your_expected_client_id"
	expectedClientSecret = "your_expected_client_secret"
)

func Deployments(w http.ResponseWriter, r *http.Request) {
	// Extract the client ID and client secret from headers
	client_ID := r.Header.Get("Client-ID")
	client_Secret := r.Header.Get("Client-Secret")

	// Check if the provided client ID and client secret match the expected values
	if client_ID != expectedClientID || client_Secret != expectedClientSecret {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	cmd := exec.Command("bash", ScriptPath)

    // Run the command and check for errors
    err := cmd.Run()
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        fmt.Println("Error executing script:", err)
        return
    }

    // Respond with success message
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Script executed successfully")

}
