package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/smithyworks/FaDO/database"
	"github.com/smithyworks/FaDO/util"
)

func SendResources(w http.ResponseWriter) {
    if resources, err := database.QueryResources(); err != nil {
        util.PrintErr(err)
        http.Error(w, "500 internal server error.", http.StatusInternalServerError)
        return
    } else {
        if resourcesJSON, err := json.Marshal(resources); err != nil {
            util.PrintErr(err)
            http.Error(w, "500 internal server error.", http.StatusInternalServerError)
            return
        } else {
            w.Header().Set("Content-Type", "application/json")
	        w.Write(resourcesJSON)
        }
    }
}

func Resources(w http.ResponseWriter, r *http.Request) {
	if !ValidateRequest(w, r, "/api/resources", "GET", nil) { return }

    resources, err := database.QueryResources()
    if err != nil {
        util.PrintErr(err)
        http.Error(w, "500 internal server error.", http.StatusInternalServerError)
        return
    }

	resourcesJSON, err := json.Marshal(resources)
	if err != nil {
        util.PrintErr(err)
        http.Error(w, "500 internal server error.", http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.Write(resourcesJSON)
}