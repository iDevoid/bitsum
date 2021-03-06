package routers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Router data will be registered to http listener
type Router struct {
	Method  string
	Path    string
	Handler fiber.Handler
}

type routing struct {
	host   string
	domain string
	routes []Router
}

// Routers contains the functions of http handler to clean payloads and pass it the service
type Routers interface {
	Serve()
}

// Initialize is for initialize the handler
func Initialize(host string, routes []Router, domain string) Routers {
	return &routing{
		host,
		domain,
		routes,
	}
}

// Serve is to start serving the HTTP Listener for every domain
func (r *routing) Serve() {
	server := fiber.New()

	group := server.Group(fmt.Sprintf("/%s", r.domain))

	for _, router := range r.routes {
		group.Add(router.Method, router.Path, router.Handler)
	}

	logrus.WithFields(logrus.Fields{
		"host":   r.host,
		"domain": r.domain,
	}).Info("Starts Serving on HTTP")
	err := server.Listen(r.host)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"host":   r.host,
			"domain": r.domain,
		}).Fatal(err)
		panic(err)
	}
}
