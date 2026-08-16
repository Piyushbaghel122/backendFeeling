package src

import (
	"fmt"
	"net/http"
	"time"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Reading a cookie (Go's built-in "cookie-parser")
	cookieToRead, err := r.Cookie("my_cookie")
	if err != nil {
		fmt.Println("Cookie not found or not set yet.")
	} else {
		fmt.Printf("Successfully read cookie! Value: %s\n", cookieToRead.Value)
	}

	// 2. Setting a new cookie
	cookieToSet := &http.Cookie{
		Name:     "my_cookie",
		Value:    "hello_cookie_value",
		Expires:  time.Now().Add(24 * time.Hour), // Expires in 24 hours
		HttpOnly: true,                           // Prevents JavaScript access (good for security)
		Path:     "/",
	}

	// Set the cookie in the response
	http.SetCookie(w, cookieToSet)

	fmt.Fprintf(w, "Hello World! A cookie has been set (and checked in the backend).")
}
