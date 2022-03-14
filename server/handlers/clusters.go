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

type ClustersInput struct {
	Cluster database.ClusterRecord `json:"cluster"`
	Zones []string `json:"zones"`
}

func (ci *ClustersInput) IsValid() bool {
	if ci.Zones == nil { ci.Zones = make([]string, 0) }
	return ci.Cluster.Name != ""
}

func Clusters(w http.ResponseWriter, r *http.Request) {
    if r.Method  == "GET" {
        SendResources(w)
		return
    } else if r.Method == "POST" {
		var input ClustersInput
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

			if err = mutations.AddCluster(tx, input.Cluster, input.Zones); err != nil {
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
		var input ClustersInput
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

			if err = mutations.EditCluster(tx, input.Cluster, input.Zones); err != nil {
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

func Cluster(w http.ResponseWriter, r *http.Request) {
    if r.Method == "DELETE" {
		cluster_id, err := strconv.Atoi(mux.Vars(r)["cluster_id"])
		if err != nil {
			util.PrintErr(fmt.Errorf("Path not found. Expected %v, got %v.", "/api/clusters?cluster_id=<int>", r.URL.RequestURI()))
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

			if cluster, err := database.QueryClusterRow(tx, "SELECT * FROM clusters WHERE cluster_id = $1", cluster_id); err != nil {
				util.PrintErr(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				if err = mutations.DeleteCluster(tx, cluster, permanent); err != nil {
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
