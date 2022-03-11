package handlers

import (
	"fmt"
	"net/http"

	"github.com/smithyworks/FaDO/util"
)

func Health(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/healthz" {
        util.PrintErr(fmt.Errorf("Path not found. Expected %v, got %v.", "/healthz", r.URL.Path))
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }

    fmt.Fprintf(w, "OK")
}
