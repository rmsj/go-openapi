package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/respond"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"

	"github.com/a-h/rest"
	"github.com/a-h/rest/chiadapter"
	"github.com/a-h/rest/examples/chiexample/models"
	"github.com/a-h/rest/swaggerui"
)

func main() {
	// Define routes in any router.
	router := chi.NewRouter()

	router.Get("/topic/{id}", func(w http.ResponseWriter, r *http.Request) {
		resp := models.Topic{
			Namespace: "example",
			Topic:     "topic",
			Private:   false,
			ViewCount: 412,
		}
		respond.WithJSON(w, resp, http.StatusOK)
	})

	router.Get("/topics", func(w http.ResponseWriter, r *http.Request) {
		resp := models.TopicsGetResponse{
			Topics: []models.TopicRecord{
				{
					ID: "testId",
					Topic: models.Topic{
						Namespace: "example",
						Topic:     "topic",
						Private:   false,
						ViewCount: 412,
					},
				},
			},
		}
		respond.WithJSON(w, resp, http.StatusOK)
	})

	router.Post("/topics", func(w http.ResponseWriter, r *http.Request) {
		resp := models.TopicsPostResponse{ID: "123"}
		respond.WithJSON(w, resp, http.StatusOK)
	})

	// Create the API definition.
	api := rest.NewAPI("Messaging API", "1.0.0")

	// Create the routes and parameters of the Router in the REST API definition with an
	// adapter, or do it manually.
	err := chiadapter.Merge(api, router)
	if err != nil {
		log.Fatalf("failed to create routes: %v", err)
	}

	// Because this example is all in the main package, we can strip the `main_` namespace from
	// the types.
	api.StripPkgPaths = []string{"main", "github.com/a-h"}

	// It's possible to customise the OpenAPI schema for each type.
	_, _, err = api.RegisterModel(*rest.ModelOf[respond.Error](), rest.WithDescription("Standard JSON error"), func(s *openapi3.Schema) {
		status := s.Properties["statusCode"]
		status.Value.WithMin(100).WithMax(600)
	})
	if err != nil {
		log.Fatalf("failed to register model: %v", err)
	}

	// Document the routes.
	api.Get("/topic/{id}").
		HasResponse(http.StatusOK, rest.ModelOf[models.TopicsGetResponse](), "topic response").
		HasResponse(http.StatusInternalServerError, rest.ModelOf[respond.Error](), "error response")

	api.Get("/topics").
		HasResponse(http.StatusOK, rest.ModelOf[models.TopicsGetResponse](), "topic response").
		HasResponse(http.StatusInternalServerError, rest.ModelOf[respond.Error](), "error response")

	api.Post("/topics").
		HasRequest(rest.ModelOf[models.TopicsPostRequest](), "topic request").
		HasResponse(http.StatusOK, rest.ModelOf[models.TopicsPostResponse](), "topic response").
		HasResponse(http.StatusInternalServerError, rest.ModelOf[respond.Error](), "error response")

	// Create the spec.
	spec, err := api.Spec()
	if err != nil {
		log.Fatalf("failed to create spec: %v", err)
	}

	// Customise it.
	spec.Info.Version = "v1.0.0"
	spec.Info.Description = "Messages API"

	// Attach the UI handler.
	ui, err := swaggerui.New(spec)
	if err != nil {
		log.Fatalf("failed to create swagger UI handler: %v", err)
	}
	router.Handle("/swagger-ui*", ui)
	// And start listening.
	fmt.Println("Listening on :8080...")
	fmt.Println("Visit http://localhost:8080/swagger-ui to see API definitions")
	http.ListenAndServe(":8080", router)
}
