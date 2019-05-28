package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

var PORT string

func init() {
	PORT = os.Getenv("PORT")
	if len(PORT) == 0 {
		PORT = "8080"
	}
}

func main() {
	// Create Server and Route Handlers
	r := mux.NewRouter()
	r.HandleFunc("/deploy", handleDeploy)
	srv := &http.Server{
		Handler: r,
		Addr:    ":" + PORT,
		// ReadTimeout:  20 * time.Second,
		// WriteTimeout: 20 * time.Second,
	}

	log.Println("Starting Server on port - " + PORT)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func handleDeploy(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	repo := r.FormValue("repo")
	port := r.FormValue("port")
	if len(repo) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("no 'repo' provided")
		return
	}

	addr, err := generateContainer(repo, port)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	fmt.Println("running on address " + addr)
	fmt.Println("complete...")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"port": addr})
}

func generateContainer(repo, port string) (string, error) {
	// get container name
	ss := strings.Split(repo, "/")
	name := strings.Split(ss[len(ss)-1], ".")[0]

	// init docker client
	fmt.Println("initializing docker client...")
	cli, err := client.NewClientWithOpts(client.WithVersion("1.37"))
	if err != nil {
		return "", err
	}

	// remove container if already exists
	fmt.Println("removing container if already exists...")
	if container, err := cli.ContainerInspect(context.Background(), name); err == nil {
		cmd := exec.Command("docker", "rm", "-f", "/"+container.Name)
		_, err := cmd.CombinedOutput()
		if err != nil {
			return "", err
		}
	}

	// generate project container
	fail := make(chan error)
	go func() {
		fmt.Println("generating container...")
		cmd := exec.Command("bash", "build.sh", repo, port)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fail <- fmt.Errorf(string(out))
		}
	}()

	// get project address
	fmt.Println("getting container address...")
	for {
		select {
		case f := <-fail:
			return "", f
		default:
			container, err := cli.ContainerInspect(context.Background(), name)
			if err != nil {
				continue
			}

			for _, ports := range container.ContainerJSONBase.HostConfig.PortBindings {
				fmt.Println("accessing address...")
				return ports[0].HostPort, nil
			}
		}
		time.Sleep(2 * time.Second)
	}
}
