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

type StorageInput struct {
	StorageDeployment database.StorageDeploymentRecord `json:"storage_deployment"`
}

func (si *StorageInput) IsValid() bool {
	return si.StorageDeployment.ClusterID != 0 && si.StorageDeployment.Alias != "" && si.StorageDeployment.Endpoint != "" && si.StorageDeployment.AccessKey != "" && si.StorageDeployment.SecretKey != ""
}

func StorageDeployments(w http.ResponseWriter, r *http.Request) {
    if r.Method  == "GET" {
		SendResources(w)
		return
    } else if r.Method == "POST" {
		var input StorageInput
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

			if err = mutations.AddStorageDeployment(tx, input.StorageDeployment); err != nil {
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

func StorageDeployment(w http.ResponseWriter, r *http.Request) {
    if r.Method == "DELETE" {
		storage_id, err := strconv.Atoi(mux.Vars(r)["storage_id"])
		if err != nil {
			util.PrintErr(fmt.Errorf("Path not found. Expected %v, got %v.", "/api/storage-deployments/<int>", r.URL.RequestURI()))
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		permanent := r.URL.Query().Get("permanent") == "true"

		if tx, err := database.Begin(); err != nil {
			util.PrintErr(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			defer tx.Rollback(ctx)

			if storage, err := database.QueryStorageDeploymentRow(tx, "SELECT * FROM storage_deployments WHERE storage_id = $1", storage_id); err != nil {
				util.PrintErr(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				if err = mutations.DeleteStorageDeployment(tx, storage, permanent); err != nil {
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
