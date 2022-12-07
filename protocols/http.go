package protocols

import (
	"io"
	"net/http"
	"strings"

	"github.com/darkcat013/pr-datastore/constants"
	"github.com/darkcat013/pr-datastore/datastore"
	"github.com/darkcat013/pr-datastore/utils"
	"github.com/gorilla/mux"
)

var NewLeaderChan = make(chan bool)

func StartHttp() {

	if !datastore.IsLeader {
		<-NewLeaderChan
	}

	router := mux.NewRouter()

	router.HandleFunc("/get/{id}", func(w http.ResponseWriter, r *http.Request) {

		p := mux.Vars(r)
		id := p["id"]

		value, err := datastore.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(value)

	}).Methods("GET")

	router.HandleFunc("/getkeys", func(w http.ResponseWriter, r *http.Request) {

		value := datastore.GetAllKeys()

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(strings.Join(value, "\n")))

	}).Methods("GET")

	router.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {

		value, err := io.ReadAll(r.Body)
		if err != nil {
			utils.Log.Info("HTTP /post | Error reading body")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("could not read request body"))
		}

		newId := TcpInsert(value)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(newId))
	}).Methods("POST")

	router.HandleFunc("/update/{id}", func(w http.ResponseWriter, r *http.Request) {

		p := mux.Vars(r)
		id := p["id"]

		value, err := io.ReadAll(r.Body)
		if err != nil {
			utils.Log.Info("HTTP /update | Error reading body")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("could not read request body"))
		}

		err = datastore.Update(id, value)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		TcpUpdate(id, value)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully updated"))
	}).Methods("PUT")

	router.HandleFunc("/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		p := mux.Vars(r)
		id := p["id"]

		err := datastore.Delete(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		TcpDelete(id)

		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("Successfully deleted"))
	}).Methods("DELETE")

	//log and handle all incoming requests
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.Log.Infow("HTTP | Requested",
			"method", r.Method,
			"path", r.RequestURI,
		)
		router.ServeHTTP(w, r)
	})

	utils.Log.Infow("HTTP | Starting server on port " + constants.HTTP_PORT)
	if err := http.ListenAndServe(constants.HTTP_PORT, handler); err != nil {
		utils.Log.Fatalw("HTTP | Could not start server",
			"error", err.Error(),
		)
	}
	utils.Log.Infow("HTTP | Started server on port " + constants.HTTP_PORT)
}
