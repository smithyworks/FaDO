package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"
	"github.com/smithyworks/FaDO/cli"
	"github.com/smithyworks/FaDO/config"
	"github.com/smithyworks/FaDO/database"
	"github.com/smithyworks/FaDO/handlers"
	"github.com/smithyworks/FaDO/mutations"
	"github.com/smithyworks/FaDO/util"
)

func initAfterReady(configFilePath, databaseConnectionString, serverURL, caddyAdminURL string) {
	time.Sleep(5 * time.Second)
	handlers.SetReady(true)

	tx, err := database.Begin()
	if err != nil { util.PrintErr(err); return }
	defer tx.Rollback(context.Background())

	rows, err := database.Query(tx, "SELECT COUNT(*) FROM clusters")
	if err != nil { util.PrintErr(err); return }
	defer rows.Close()

	var clusterCount int
	if rows.Next() { err = rows.Scan(&clusterCount); rows.Close() }
	if err != nil { util.PrintErr(err); return }

	if clusterCount > 0 {
		log.Printf("INFO: Found database populated, skipping configuration from file.")
		log.Printf("INFO: Syncing database state with MinIO deployments.")
		if storageDeployments, err := database.QueryStorageDeployments(tx, "SELECT * FROM storage_deployments"); err != nil {
			util.PrintErr(err)
			return
		} else {
			for _, sd := range storageDeployments {
				if err = mutations.AddStorageDeployment(tx, sd); err != nil {
					util.PrintWarning(err)
					err = nil
				}
			}
		}
		if buckets, err := database.QueryBuckets(tx, "SELECT * FROM buckets"); err != nil {
			util.PrintErr(err)
			return
		} else {
			for _, b := range buckets {
				if err = mutations.TrackBucketObjects(tx, b); err != nil {
					util.PrintWarning(err)
					err = nil
				}
				if err = mutations.ResolveBucketReplicas(tx, b); err != nil {
					util.PrintWarning(err)
					err = nil
				}
			}
		}
		mutations.ConfigureLoadBalancer(tx)
		tx.Commit(context.Background())
	} else {
		log.Printf("INFO: Loading initial configuration file at '%v'.", configFilePath)
		if err = config.LoadConfigFromFile(configFilePath); err != nil {
			util.PrintErr(err)
			return
		}
	}
	log.Printf("INFO: Finished loading.")
}

//go:embed build/*
var staticFiles embed.FS

// fsFunc is short-hand for constructing a http.FileSystem
// implementation
type fsFunc func(name string) (fs.File, error)

func (f fsFunc) Open(name string) (fs.File, error) {
	return f(name)
}

func OpenFile(name string) (fs.File, error) {
	filePath := path.Join("build", name)

	f, err := staticFiles.Open(filePath)
	if os.IsNotExist(err) {
		return staticFiles.Open("build/index.html")
	}

	return f, err
} 

// Handle serving the single-page frontend application.
func SpaHandler() http.Handler {
    return http.FileServer(http.FS(fsFunc(OpenFile)))
}

func main() {
	// Server Init

	input := cli.ReadInput()

	log.Printf("INFO: Connecting to database.")
	if err := database.Connect(cli.Input.DatabaseConnectionString); err != nil {
		util.PrintErr(err)
		log.Fatal("Exiting.")
		return
	}
	defer database.Close()

	go initAfterReady(input.ConfigFilePath, input.DatabaseConnectionString, input.ServerURL, input.CaddyAdminURL)
	
	// HTTP Server

    r := mux.NewRouter()

	r.HandleFunc("/api/notify", handlers.Notify)
	r.HandleFunc("/api/resources", handlers.Resources)
	r.HandleFunc("/api/clusters", handlers.Clusters)
	r.HandleFunc("/api/clusters/{cluster_id:[0-9]+}", handlers.Cluster)
	r.HandleFunc("/api/faas-deployments", handlers.FaaSDeployments)
	r.HandleFunc("/api/faas-deployments/{faas_id:[0-9]+}", handlers.FaaSDeployment)
	r.HandleFunc("/api/storage-deployments", handlers.StorageDeployments)
	r.HandleFunc("/api/storage-deployments/{storage_id:[0-9]+}", handlers.StorageDeployment)
	r.HandleFunc("/api/buckets", handlers.Buckets)
	r.HandleFunc("/api/buckets/{bucket_id:[0-9]+}", handlers.Bucket)
	r.HandleFunc("/api/objects", handlers.Objects)
	r.HandleFunc("/api/objects/{object_id:[0-9]+}", handlers.Object)
	r.HandleFunc("/api/load-balancer", handlers.LoadBalancer)
	r.HandleFunc("/api/load-balancer/settings", handlers.LoadBalancer)
	r.HandleFunc("/api/load-balancer/route-overrides", handlers.LoadBalancer)

	r.HandleFunc("/healthz", handlers.Health)

	spaHandler := SpaHandler()
	r.PathPrefix("/").Handler(spaHandler)

	http.Handle("/", r)

	log.Printf("INFO: Starting server at port 9090.\n")
    if err := http.ListenAndServe(":9090", nil); err != nil {
		util.PrintErr(err)
        log.Fatal("Exiting.")
    }
}
