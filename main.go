package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dpc-sdp/bay-section-ip-controller/internal/handler"
	"github.com/dpc-sdp/bay-section-ip-controller/internal/middleware"
	"github.com/dpc-sdp/bay-section-ip-controller/internal/util"
)

var (
	version string
	commit  string

	// Variables for running the webserver.
	port             = flag.String("p", "80", "TCP listen port")
	blockedIps       = flag.String("b", "", "Comma separated list of IPs to always include in the blocklist")
	environments     = flag.String("e", "Develop", "Comma separated list of environments to update")
	applications     = flag.String("a", "", "Comma separate list of applications to update")
	sectionUsername  = flag.String("u", os.Getenv("SECTION_IO_USERNAME"), "User for Section API")
	sectionToken     = flag.String("t", os.Getenv("SECTION_IO_TOKEN"), "Token for Section API")
	sectionAccountId = flag.String("i", os.Getenv("SECTION_IO_ACCOUNT_ID"), "Account ID for Section API")
)

func main() {
	flag.Parse()

	s := util.SectionAPI{
		Username:               *sectionUsername,
		Token:                  *sectionToken,
		AccountId:              *sectionAccountId,
		ActionableEnvironments: strings.Split(*environments, ","),
		ActionableApplications: strings.Split(*applications, ","),
		Host:                   "https://aperture.section.io/api/v1",
		Client:                 &http.Client{},
		BlockedIps: util.SectionIpRestrictionSchema{
			IpBlacklist: strings.Split(*blockedIps, ","),
		},
	}

	s.Init()

	router := http.NewServeMux()
	router.HandleFunc("/_healthz", (&handler.HealthCheck{Section: s}).Serve)
	router.HandleFunc("/v1/ip/add", (&handler.ThreatIPSavedSearch{Section: s}).Serve)

	// Register the middleware.
	username := os.Getenv("BASIC_AUTH_USERNAME")
	password := os.Getenv("BASIC_AUTH_PASSWORD")

	b := middleware.BasicAuth{
		Username:  username,
		Password:  password,
		AppliesTo: []string{"/v1/ip/add"},
	}

	handler := applyMiddleware(router, b.Do)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), handler))
}

func applyMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
