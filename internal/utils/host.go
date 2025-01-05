package utils

import (
	"fmt"
	"net/url"
	"os"
)

func ServerHost() *url.URL {
	if os.Getenv("APP_ENV") == "local" {
		var u url.URL
		u.Scheme = "http"
		u.Host = fmt.Sprintf("localhost:%s", os.Getenv("PORT"))
		return &u
	}

	u, err := url.Parse(os.Getenv("DEPLOYMENT_URL"))
	if err != nil {
		panic(err)
	}

	return u
}
