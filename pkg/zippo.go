package zippo

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	bc "github.com/coreos/butane/config"
	"github.com/coreos/butane/config/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type args struct {
	Hostname      string
	SSHUser       string
	SSHPubkey     string
	CNIVersion    string
	CRICtlVersion string
	K8sVersion    string
}

type Server struct {
	Config *ServerConfig
}

type ServerConfig struct {
	Address        string
	Port           int
	SSHUser        string
	SSHPubkey      string
	CNIVersion     string
	CRICtlVersion  string
	K8sVersion     string
	ButaneTemplate string
}

func (s *Server) Run() {
	router := gin.Default()
	router.GET("/ignition", s.serveButaneTranslator)
	router.Run()
}

func (s *Server) serveButaneTranslator(c *gin.Context) {
	// Set hostname from host header
	hostname := strings.Split(c.Request.Host, ":")[0]
	if hostname == "" {
		c.JSON(http.StatusInternalServerError, "No hostname specified")
		return
	}
	log.Infof("Serving ignition config for %s at %s", hostname, c.Request.RemoteAddr)

	// Translate human-readable Butane config to Ignition
	ignitionConfig, err := createIgnitionConfig(hostname, s.Config)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, string(ignitionConfig))
}

func createIgnitionConfig(hostname string, config *ServerConfig) ([]byte, error) {
	log.Infof("Using config at %s", config.ButaneTemplate)

	// Parse config template
	file, err := ioutil.ReadFile(config.ButaneTemplate)
	if err != nil {
		return nil, err
	}

	butaneTemplate, err := template.New("butaneConfig").Parse(string(file))
	if err != nil {
		return nil, err
	}

	// Render butane config template with given hostname
	vars := args{
		Hostname:      hostname,
		SSHUser:       config.SSHUser,
		SSHPubkey:     config.SSHPubkey,
		CNIVersion:    config.CNIVersion,
		CRICtlVersion: config.CRICtlVersion,
		K8sVersion:    config.K8sVersion,
	}
	butaneConfig := &bytes.Buffer{}
	err = butaneTemplate.Execute(butaneConfig, vars)
	if err != nil {
		return nil, err
	}
	log.Info("Populated Butane template")

	ignitionConfig, r, err := bc.TranslateBytes(butaneConfig.Bytes(), common.TranslateBytesOptions{Pretty: true})
	if err != nil {
		return nil, errors.Wrapf(err, "Error translating config: %s", r.String())
	}
	log.Info("Generated Ignition config")

	return ignitionConfig, nil
}
