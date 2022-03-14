package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/smithyworks/FaDO/database"
	"github.com/smithyworks/FaDO/mutations"
	"github.com/smithyworks/FaDO/util"
)

type LBSettingsInput struct {
	MatchHeader string `json:"match_header"`
	Policy string `json:"policy"`
}

func (i *LBSettingsInput) IsValid() bool {
	return i.Policy != "" && i.MatchHeader != ""
}

type LBOverridesInput struct {
	RouteOverrides map[string]database.LoadBalancerRouteSettings `json:"route_overrides"`
}

func (i *LBOverridesInput) IsValid() bool {
	if i.RouteOverrides == nil { i.RouteOverrides = make(map[string]database.LoadBalancerRouteSettings) }
	return true
}

func LoadBalancer(w http.ResponseWriter, r *http.Request) {
    if r.Method  == "GET" {
        SendResources(w)
		return
    } else if r.Method == "PUT" {
		if r.URL.Path == "/api/load-balancer/settings" {
			var input LBSettingsInput
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

				if err := database.SetGlobalPolicy(tx, "lb_match_header", input.MatchHeader); err != nil {
					util.PrintErr(err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				if err := database.SetGlobalPolicy(tx, "lb_policy", input.Policy); err != nil {
					util.PrintErr(err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				if err := mutations.ConfigureLoadBalancer(tx); err != nil {
					util.PrintErr(err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				if err := tx.Commit(ctx); err != nil {
					util.PrintErr(err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}
		} else if r.URL.Path == "/api/load-balancer/route-overrides" {
			var input LBOverridesInput
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
	
				if err := database.SetGlobalPolicy(tx, "lb_route_overrides", input.RouteOverrides); err != nil {
					util.PrintErr(err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				if err := mutations.ConfigureLoadBalancer(tx); err != nil {
					util.PrintErr(err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				if err := tx.Commit(ctx); err != nil {
					util.PrintErr(err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}
		} else {
			util.PrintErr(fmt.Errorf("Path not found. Expected %v or %v, got %v.", "/api/load-balancer/settings", "/api/load-balancer/route-overrides", r.URL.Path))
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		SendResources(w)
		return
	} else {
		util.PrintErr(fmt.Errorf("Method not supported. Got %v.", r.Method))
		http.Error(w, "Method Not Supported", http.StatusNotFound)
		return
	}
}