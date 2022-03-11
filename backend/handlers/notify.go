package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/smithyworks/FaDO/database"
	"github.com/smithyworks/FaDO/mutations"
	"github.com/smithyworks/FaDO/util"
)

type NotifyRecordResponseElements struct {
    DeploymentID string `json:"x-minio-deployment-id"`
}

type NotifyRecord struct {
    ResponseElements NotifyRecordResponseElements `json:"responseElements"`
}

type NotifyInput struct {
    EventName string `json:"EventName"`
    Key string `json:"Key"`
    Records []NotifyRecord `json:"Records"`
}

func Notify(w http.ResponseWriter, r *http.Request) {
    var input NotifyInput
    if r.URL.Path != "/api/notify" {
		util.PrintErr(fmt.Errorf("Path not found. Expected %v, got %v.", "/api/notify", r.URL.Path))
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }

    fmt.Fprintf(w, "OK")

    err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil || len(input.Records) < 1 { log.Println("INFO: Notify: input not actionable.") }

    name := input.Key

    var minioDeploymentID string
    for _, r := range input.Records {
        if r.ResponseElements.DeploymentID != "" {
            minioDeploymentID = r.ResponseElements.DeploymentID
        }
    }
    if minioDeploymentID == "" { return }

    nameTokens := strings.Split(name, "/")
    bucketName := nameTokens[0]

    log.Printf("INFO: Notify on %v (%v).", bucketName, minioDeploymentID)

    tx, err := database.Begin()
    if err != nil { util.PrintErr(err); return }
    defer tx.Rollback(ctx)

    bucket, err := database.QueryBucketRow(tx, "SELECT * FROM buckets WHERE name = $1", bucketName)
    if err != nil { util.PrintErr(err); return }
    storageDeployment, err := database.QueryStorageDeploymentRow(tx, "SELECT * FROM storage_deployments WHERE minio_deployment_id = $1", minioDeploymentID)

    if bucket.StorageID != storageDeployment.StorageID { return }

    if err = mutations.TrackBucketObjects(tx, bucket); err != nil {
        util.PrintErr(err); return
    }

    if err = mutations.ReplicateBucket(tx, bucketName); err != nil {
        util.PrintErr(err); return
    }

    tx.Commit(ctx)
}
