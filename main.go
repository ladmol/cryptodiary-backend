package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func main() {
	err := supertokens.Init(SuperTokensConfig)

	if err != nil {
		panic(err.Error())
	}

	// Initialize a new Chi router
	r := chi.NewRouter()

	// Add useful Chi middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(corsMiddleware)
	r.Use(supertokens.Middleware)

	// Define routes
	r.Get("/hello", hello)

	// Health check endpoint for the API and SuperTokens connection
	r.Get("/health", healthCheck)

	// Protected routes using session verification
	r.Get("/sessioninfo", func(w http.ResponseWriter, r *http.Request) {
		// We need to manually call VerifySession here since it returns http.HandlerFunc
		// and not a middleware that chi's router.Use expects
		session.VerifySession(nil, sessioninfo).ServeHTTP(w, r)
	})

	// Tenants endpoint
	r.Get("/tenants", tenants)
	// Not found handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	})

	log.Println("Starting server on port 3001...")
	log.Printf("Connected to SuperTokens core at: %s\n", getSuperTokensURI())
	http.ListenAndServe(":3001", r)
}

// corsMiddleware handles CORS (Cross-Origin Resource Sharing) headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			response.Header().Set("Access-Control-Allow-Headers", strings.Join(append([]string{"Content-Type"}, supertokens.GetAllCORSHeaders()...), ","))
			response.Header().Set("Access-Control-Allow-Methods", "*")
			response.Write([]byte(""))
			return
		}
		next.ServeHTTP(response, r)
	})
}

// JSON returns a formatted JSON response with a specific status code
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "error in converting to json"})
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

// healthCheck checks if the server is healthy and SuperTokens core is reachable
func healthCheck(w http.ResponseWriter, r *http.Request) {
	// Try to get tenant information as a way to verify SuperTokens core connection
	_, err := multitenancy.ListAllTenants()
	status := "healthy"
	statusCode := http.StatusOK

	if err != nil {
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
		JSON(w, statusCode, map[string]interface{}{
			"status":  status,
			"message": "SuperTokens core connection failed: " + err.Error(),
		})
		return
	}

	JSON(w, statusCode, map[string]interface{}{
		"status":  status,
		"message": "API is running and SuperTokens core is reachable",
	})
}

func sessioninfo(w http.ResponseWriter, r *http.Request) {
	sessionContainer := session.GetSessionFromRequestContext(r.Context())

	if sessionContainer == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("no session found"))
		return
	}

	sessionData, err := sessionContainer.GetSessionDataInDatabase()
	if err != nil {
		err = supertokens.ErrorHandler(err, r, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		return
	}

	responseData := map[string]interface{}{
		"sessionHandle":      sessionContainer.GetHandle(),
		"userId":             sessionContainer.GetUserID(),
		"accessTokenPayload": sessionContainer.GetAccessTokenPayload(),
		"sessionData":        sessionData,
	}

	JSON(w, http.StatusOK, responseData)
}

func tenants(w http.ResponseWriter, r *http.Request) {
	tenantsList, err := multitenancy.ListAllTenants()

	if err != nil {
		err = supertokens.ErrorHandler(err, r, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		return
	}

	responseData := map[string]interface{}{
		"status":  "OK",
		"tenants": tenantsList.OK.Tenants,
	}

	JSON(w, http.StatusOK, responseData)
}
