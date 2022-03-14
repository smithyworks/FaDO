package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/smithyworks/FaDO/database"
	"github.com/smithyworks/FaDO/mutations"
	"github.com/smithyworks/FaDO/util"
)

type BucketsInput struct {
	Bucket database.BucketRecord `json:"bucket"`
	TargetReplicaCount int `json:"target_replica_count"`
	Zones []string `json:"zones"`
	ReplicaStorageIDs []int64 `json:"replica_storage_ids"`
}

func (bi *BucketsInput) IsValid() bool {
	if bi.Zones == nil { bi.Zones = make([]string, 0) }
	return bi.Bucket.Name != "" && bi.Bucket.StorageID != 0
}

func Buckets(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        SendResources(w)
		return
    } else if r.Method == "POST" {
		var input BucketsInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			util.PrintErr(err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		} else if !input.IsValid() {
			util.PrintErr(fmt.Errorf("Invalid input. Got %+v of type %v.", input, reflect.TypeOf(input)))
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if tx, err := database.Begin(); err != nil {
			util.PrintErr(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			defer tx.Rollback(ctx)

			if err = mutations.AddMasterBucket(tx, input.Bucket, input.TargetReplicaCount, input.Zones, input.ReplicaStorageIDs); err != nil {
				util.PrintErr(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				tx.Commit(ctx)
			}
		}

		SendResources(w)
		return
	} else if r.Method == "PUT" {
		var input BucketsInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			util.PrintErr(err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		} else if !input.IsValid() {
			util.PrintErr(fmt.Errorf("Invalid input. Got %+v of type %v.", input, reflect.TypeOf(input)))
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if tx, err := database.Begin(); err != nil {
			util.PrintErr(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			defer tx.Rollback(ctx)
			if err = mutations.EditMasterBucket(tx, input.Bucket, input.TargetReplicaCount, input.Zones, input.ReplicaStorageIDs); err != nil {
				util.PrintErr(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				tx.Commit(ctx)
			}
		}

		SendResources(w)
		return
	} else {
		util.PrintErr(fmt.Errorf("Method not supported. Got %v.", r.Method))
		http.Error(w, "Method Not Supported", http.StatusNotFound)
		return
	}
}

func Bucket(w http.ResponseWriter, r *http.Request) {
    if r.Method == "DELETE" {
		bucket_id, err := strconv.Atoi(mux.Vars(r)["bucket_id"])
		if err != nil {
			util.PrintErr(fmt.Errorf("Path not found. Expected %v, got %v.", "/api/buckets?bucket_id=<int>", r.URL.RequestURI()))
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		if tx, err := database.Begin(); err != nil {
			util.PrintErr(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			defer tx.Rollback(ctx)

			if bucket, err := database.QueryBucketRow(tx, "SELECT * FROM buckets WHERE bucket_id = $1", bucket_id); err != nil {
				util.PrintErr(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				if err = mutations.DeleteMasterBucket(tx, bucket); err != nil {
					util.PrintErr(err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				} else {
					tx.Commit(ctx)
				}
			}
		}

		SendResources(w)
		return
	} else {
		util.PrintErr(fmt.Errorf("Method not supported. Got %v.", r.Method))
		http.Error(w, "Method Not Supported", http.StatusNotFound)
		return
	}
}
