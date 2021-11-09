package zippo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/coreos/butane/config"
	"github.com/coreos/butane/config/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type args struct {
	Hostname string
}

type Server struct {
	Config *ServerConfig
}

type ServerConfig struct {
	Address        string
	Port           int
	ButaneTemplate string
}

func (s *Server) Run() {
	router := gin.Default()
	router.GET("/ignition", s.serveButaneTranslator)
	router.Run()
}

func (s *Server) serveButaneTranslator(c *gin.Context) {
	// Set hostname from host header
	hostname := c.Request.Host
	if hostname == "" {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	log.Print(hostname)

	// Translate human-readable Butane config to Ignition
	ignitionConfig, err := createIgnitionConfig(hostname, s.Config.ButaneTemplate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, string(ignitionConfig))
}

func createIgnitionConfig(hostname, butaneFile string) ([]byte, error) {
	// Parse config template
	file, err := ioutil.ReadFile(butaneFile)
	if err != nil {
		return nil, err
	}

	butaneTemplate, err := template.New("butaneConfig").Parse(string(file))
	if err != nil {
		return nil, err
	}

	// Render butane config template with given hostname
	butaneConfig := &bytes.Buffer{}
	err = butaneTemplate.Execute(butaneConfig, args{Hostname: hostname})
	if err != nil {
		return nil, err
	}
	fmt.Println(butaneConfig.String())

	ignitionConfig, r, err := config.TranslateBytes(butaneConfig.Bytes(), common.TranslateBytesOptions{Pretty: true})
	if err != nil {
		return nil, errors.Wrapf(err, "Error translating config: %s", r.String())
	}

	return ignitionConfig, nil
}
