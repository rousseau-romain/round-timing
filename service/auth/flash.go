package auth

import (
	"encoding/json"
	"net/http"

	"github.com/markbates/goth/gothic"
)

const flashSessionName = "flash"

type flashData struct {
	Title    string   `json:"title"`
	Messages []string `json:"messages"`
	Type     string   `json:"type"`
}

// SetFlash stores a flash message in the session.
func SetFlash(w http.ResponseWriter, r *http.Request, title string, messages []string, msgType string) {
	session, _ := gothic.Store.Get(r, flashSessionName)
	data, _ := json.Marshal(flashData{
		Title:    title,
		Messages: messages,
		Type:     msgType,
	})
	session.AddFlash(string(data))
	session.Save(r, w)
}

// GetFlash retrieves and clears the flash message from the session.
func GetFlash(w http.ResponseWriter, r *http.Request) (title string, messages []string, msgType string) {
	session, err := gothic.Store.Get(r, flashSessionName)
	if err != nil {
		return "", nil, ""
	}
	flashes := session.Flashes()
	if len(flashes) == 0 {
		return "", nil, ""
	}

	// Save to clear the flash
	session.Save(r, w)

	raw, ok := flashes[0].(string)
	if !ok {
		return "", nil, ""
	}
	var fd flashData
	if err := json.Unmarshal([]byte(raw), &fd); err != nil {
		return "", nil, ""
	}
	return fd.Title, fd.Messages, fd.Type
}
