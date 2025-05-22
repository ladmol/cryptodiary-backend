package main

import (
	"os"

	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// getSuperTokensURI returns the URI for the SuperTokens core
// It will use the environment variable SUPERTOKENS_URI if available,
// otherwise it will default to the local instance
func getSuperTokensURI() string {
	uri := os.Getenv("SUPERTOKENS_URI")
	if uri == "" {
		return "http://localhost:3567" // Default for local development
	}
	return uri
}

var SuperTokensConfig = supertokens.TypeInput{
	Supertokens: &supertokens.ConnectionInfo{
		ConnectionURI: getSuperTokensURI(),
		APIKey:        "", // No API key for local development
	},
	AppInfo: supertokens.AppInfo{
		AppName:       "CryptoDiary App",
		APIDomain:     "http://localhost:3001",
		WebsiteDomain: "http://localhost:3000",
	},
	RecipeList: []supertokens.Recipe{
		emailpassword.Init(nil),
		session.Init(nil),
		dashboard.Init(nil),
	},
}
