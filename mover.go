package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"

	"github.com/gorilla/mux"
)

var project string = os.Getenv("PROJECT")
var instance string = os.Getenv("INSTANCE")

func moveHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	currentZone := vars["current_zone"]
	targetZone := vars["target_zone"]

	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}

	rb := &compute.InstanceMoveRequest{
		TargetInstance:  fmt.Sprintf("zones/%s/instances/%s", currentZone, instance),
		DestinationZone: targetZone,
	}

	resp, err := computeService.Projects.MoveInstance(project, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	log.Print(resp, resp.Status)
}

func main() {
	log.Print("Moving instance: ", instance)

	r := mux.NewRouter()

	r.HandleFunc("/{current_zone}/{target_zone}", moveHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
