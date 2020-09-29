package ip

import (
	"encoding/json"
	"github.com/Tomer-Alony/IPLocator/src/services"
	"github.com/go-chi/chi"
	"net/http"
)

var ipService *services.IPService

func Routes(repo *services.IPService) *chi.Mux {
	ipService = repo

	router := chi.NewRouter()
	router.HandleFunc("/find-country", findCountry)

	return router
}

func findCountry(w http.ResponseWriter, r *http.Request) {

	ip := r.URL.Query().Get("ip")
	w.Header().Set("Content-Type", "application/json")

	if ip != "" {
		var ipObject, err = ipService.FindCountry(ip)

		if err == nil {
			json.NewEncoder(w).Encode(ipObject)
		} else {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(err)
		}
	}
}