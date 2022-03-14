package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/smithyworks/FaDO/util"
)

var isReady = false

var readyMutex sync.Mutex

func SetReady(ready bool) {
	readyMutex.Lock()
	defer readyMutex.Unlock()

	isReady = ready
}

func IsReady() bool {
	return isReady
}

var ctx = context.Background()

type validationInput interface {
    IsValid() bool
}

func ValidateRequest(w http.ResponseWriter, r *http.Request, path, method string, inputVar validationInput) bool {
	if r.URL.Path != path {
		util.PrintErr(fmt.Errorf("Path not found. Expected %v, got %v.", path, r.URL.Path))
        http.Error(w, "Not Found", http.StatusNotFound)
        return false
    }

    if r.Method != method {
		util.PrintErr(fmt.Errorf("Method not supported. Expected %v, got %v.", method, r.Method))
        http.Error(w, "Method Not Supported", http.StatusNotFound)
        return false
    }

	if inputVar != nil {
		err := json.NewDecoder(r.Body).Decode(inputVar)
		if err != nil || !inputVar.IsValid() {
			util.PrintErr(err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return false
		}
	}

	return true
}