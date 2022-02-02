package http_api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	url_handler "github.com/msm-k8s-svc-helper/pkg/url-handler"
)

// API -> api type with embedded cluster api-server client
type API struct {
	url_handler.UrlHandler // compose k8s api-server client
	endpoints              map[string]func(http.ResponseWriter, *http.Request)
}

// Given a service clusterIP, return the src endpoints behind it
func (api *API) svcEndpointHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case http.MethodGet:
		resource := r.URL.Query().Get("url")

		urls := api.GetInternalURLs(resource)
		log.Printf("request handled: %s -> : %v", resource, urls)

		err := json.NewEncoder(w).Encode(urls)

		if err != nil {
			log.Printf("Error encoding service endpoints: %s", err.Error())
			fmt.Fprintf(w, "Error encoding url endpoints: %s", err.Error())
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// linkEndpoints -> create mapping of url paths to function handlers
func (api *API) addEndpoints() {
	api.endpoints = map[string]func(http.ResponseWriter, *http.Request){
		"/apiv1/url-routing": api.svcEndpointHandler,
	}
}

func StartAPI(listenAddr string) {
	api := API{}
	api.addEndpoints()
	api.InitializeUrlHandler()

	// start http handlers
	for url, handler := range api.endpoints {
		http.HandleFunc(url, handler)
	}

	err := http.ListenAndServe(listenAddr, nil)

	if err != nil {
		log.Fatalf("error serving: %v", err)
	}
}
