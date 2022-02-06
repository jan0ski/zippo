package zippo

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/coreos/butane/config"
	"github.com/coreos/butane/config/common"
	"github.com/pkg/errors"
)

// Render renders a Butane config template file with variables from `args`
func Render(templatePath string, args interface{}) (*bytes.Buffer, error) {
	// Parse config template
	file, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}

	template, err := template.New("butaneConfig").Parse(string(file))
	if err != nil {
		return nil, err
	}

	// Render config template with given arguments
	config := &bytes.Buffer{}
	err = template.Execute(config, args)
	if err != nil {
		return config, err
	}

	return config, nil
}

// CreateIgnitionConfig creates an ignition config from a rendered butane template with a given hostname
func CreateIgnitionConfig(butaneTemplate string, hostname interface{}) ([]byte, error) {
	butaneConfig, err := Render(butaneTemplate, hostname)
	if err != nil {
		return nil, err
	}

	ignitionConfig, r, err := config.TranslateBytes(butaneConfig.Bytes(), common.TranslateBytesOptions{Pretty: true})
	if err != nil {
		return nil, errors.Wrapf(err, "error translating config: %s", r.String())
	}

	return ignitionConfig, nil
}
