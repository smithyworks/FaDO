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

type FaaSInput struct {
	FaaSDeployment database.FaaSDeploymentRecord `json:"faas_deployment"`
}

func (fi *FaaSInput) IsValid() bool {
	return fi.FaaSDeployment.URL != "" && fi.FaaSDeployment.ClusterID != 0
}

type FaaSDeleteInput struct {
	FaaSID int64 `json:"faas_id"`
}

func FaaSDeployments(w http.ResponseWriter, r *http.Request) {
    if r.Method  == "GET" {
		SendResources(w)
		return
    } else if r.Method == "POST" {
		var input FaaSInput
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

			if err = mutations.AddFaaSDeployment(tx, input.FaaSDeployment); err != nil {
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
		var input FaaSInput
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

			if err = mutations.EditFaaSDeployment(tx, input.FaaSDeployment); err != nil {
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

func FaaSDeployment(w http.ResponseWriter, r *http.Request) {
    if r.Method == "DELETE" {
		faas_id, err := strconv.Atoi(mux.Vars(r)["faas_id"])
		if err != nil {
			util.PrintErr(fmt.Errorf("Path not found. Expected %v, got %v.", "/api/faas-deployments/<int>", r.URL.RequestURI()))
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		if tx, err := database.Begin(); err != nil {
			util.PrintErr(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			defer tx.Rollback(ctx)

			if faas, err := database.QueryFaaSDeploymentRow(tx, "SELECT * FROM faas_deployments WHERE faas_id = $1", faas_id); err != nil {
				util.PrintErr(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				if err = mutations.DeleteFaaSDeployment(tx, faas); err != nil {
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
