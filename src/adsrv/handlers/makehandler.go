package handlers

import (
	"adsrv/msg"
	"adsrv/util"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

func Make(actualHandler func(util.Config, *sql.DB, *msg.MessageWriter, *msg.MessageReader, string) (int, *string, error), conf util.Config, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			fmt.Printf("non POST (%s) request denied from %s\n", r.Method, r.RemoteAddr)
			http.Error(w, "Only POST allowed to this URI", http.StatusMethodNotAllowed)
			return
		}

		const MASSIVE_CONTENT_TYPE = "application/massive"
		if r.Header.Get("Content-Type") != MASSIVE_CONTENT_TYPE {
			fmt.Printf("Non-Massive (%s) request denied from %s\n", r.Header.Get("Content-Type"), r.RemoteAddr)
			http.Error(w, "Content-type must be "+MASSIVE_CONTENT_TYPE, http.StatusUnsupportedMediaType)
			return
		}

		wr := msg.NewWriter()
		rdr := msg.NewReader(r.Body)

		statusCode, tokenPtr, err := actualHandler(conf, db, wr, rdr, r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		var bytes []byte
		bytes, err = wr.Finalize(tokenPtr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", MASSIVE_CONTENT_TYPE)
		w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
		w.WriteHeader(statusCode)
		w.Write(bytes)
	}
}
