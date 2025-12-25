package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type Session struct {
	Cache   map[string]string
	History []ResponseData
}

var sessions = make(map[string]*Session)

func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func getSession(r *http.Request) (*Session, string) {
	// Try to get existing session cookie
	cookie, err := r.Cookie("session_id")
	var sessionID string

	if err != nil || cookie.Value == "" {
		// No session exists, create new one
		sessionID = generateSessionID()
	} else {
		sessionID = cookie.Value
	}

	mu.Lock()
	defer mu.Unlock()

	// Get or create session data
	session, exists := sessions[sessionID]
	if !exists {
		session = &Session{
			Cache:   make(map[string]string),
			History: make([]ResponseData, 0),
		}
		sessions[sessionID] = session
	}

	return session, sessionID
}
