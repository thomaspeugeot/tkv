/*
Package main of gao_tkv contains code for the runtime server that is deployed
on the google application engine (gae) as well as the html / javascript / css files of the web site.
*/
package main

import (
	"fmt"
	"net/http"

	"github.com/thomaspeugeot/tkv/handler"
	"google.golang.org/appengine"
)

// attach all handlers
func main() {

	http.HandleFunc("/translateLatLngInSourceCountryToLatLngInTargetCountry",
		handler.GetTranslationResult)
	http.HandleFunc("/checkEnv", checkEnv)

	// that is all that is needed to serve the file at the root level
	// (check app.yaml for definition of files that are uploaded by the server)
	appengine.Main()
}

func checkEnv(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "IsDevAppServer: %v\n", appengine.IsDevAppServer())
}
