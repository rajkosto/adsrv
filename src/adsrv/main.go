// adsrv project main.go
package main

import (
	"adsrv/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/adsrv/locateService", handlers.Make(handlers.LocateServiceHandler))
	http.HandleFunc("/adsrv/openSession", handlers.Make(handlers.OpenSessionHandler))
	http.ListenAndServe(":8123", nil)
}
