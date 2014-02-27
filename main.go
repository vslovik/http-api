package main

import (
	"github.com/motain/mux"
	"flag"
	"fmt"
	_ "github.com/motain/mysql"
	"net/http"
	_ "net/http/pprof"
	"github.com/motain/sCoreAdmin/cfgloader"
	"log/syslog"
	"github.com/motain/sCoreAdmin/controller"
)

var (
	logger  *syslog.Writer
	env     = "dev"
	verbose = false
	port uint = 8484
)

func init() {
	flag.StringVar(&env, "env", env, "Application Environment")
	flag.BoolVar(&verbose, "v", verbose, "Provider verbose output")
}

func checkParsedFlags() {
	validEnvs := map[string]bool{
		"dev":  true,
		"prod": true,
	}

	if _, ok := validEnvs[env]; !ok {
		panic(fmt.Sprintf("Invalid env given: %s, supported environments: %+v", env, validEnvs))
	}
}

func main() {
	flag.Parse()
	config := cfgloader.GetConfig(env)
	router := GetRouter(env, config)
	http.Handle("/", router)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// Handles application routing
func GetRouter(env string, config *cfgloader.Config) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/{controller}/{id:[0-9]*}", func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			switch {
			case vars["controller"] == "country":
				controller.GetCountryController(env, config).GetCoreRouting(w, r)
			case vars["controller"] == "section":
				controller.GetSectionController(env, config).GetCoreRouting(w, r)
			case vars["controller"] == "competition":
				controller.GetCompetitionController(env, config).GetCoreRouting(w, r)
			case vars["controller"] == "top_competition":
				controller.GetTopCompetitionController(env, config).GetCoreRouting(w, r)
			default:
				controller.GetBaseController(env, config).WriteBadRequestResponse(w)
			}
		})

	// Parent - child relations
	router.HandleFunc("/{controller}/{entity_id:[0-9]+}/{child}/{child_id:[0-9a-z]*}", func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			switch {
			case vars["controller"] == "section":
				controller.GetSectionController(env, config).ManagementRouting(w, r)
			case vars["controller"] == "competition":
				controller.GetCompetitionController(env, config).ManagementRouting(w, r)
			case vars["controller"] == "top_competition":
				controller.GetTopCompetitionController(env, config).ManagementRouting(w, r)
			default:
				controller.GetBaseController(env, config).WriteBadRequestResponse(w)
			}
		})

	return router
}
