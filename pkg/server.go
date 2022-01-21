package zippo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Config *serverConfig
}

type serverConfig struct {
	Address      string
	Port         int
	TemplatePath string
	Args         interface{}
}

func httpError(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": err})
}

func NewServer(address string, port int, templatePath string, args interface{}) *Server {
	return &Server{&serverConfig{
		Address:      address,
		Port:         port,
		TemplatePath: templatePath,
		Args:         args,
	}}
}

func (s *Server) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/ignition", s.serveButaneTranslator)

	address := fmt.Sprintf("%s:%d", s.Config.Address, s.Config.Port)
	srv := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	log.Infof("Zippo 🔥")
	log.Infof("Starting Butane translator web service: %s", address)
	log.Fatalln(srv.ListenAndServe())
}

func (s *Server) serveButaneTranslator(w http.ResponseWriter, r *http.Request) {
	// Set hostname from host header
	hostname := strings.Split(r.Host, ":")[0]
	if hostname == "" {
		httpError(w, "No hostname specified")
		return
	}
	log.Debugf("Served ignition config for %s at %s", hostname, r.RemoteAddr)

	ignitionConfig, err := CreateIgnitionConfig(s.Config.TemplatePath, hostname)
	if err != nil {
		log.Error(err.Error())
		httpError(w, "failed to render ignition config")
	} else {
		log.Infof("Successfully rendered ignition config")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ignitionConfig)
	}
}
