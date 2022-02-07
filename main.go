package main

import (
	"flag"
	"io/ioutil"

	zippo "github.com/jan0ski/zippo/pkg"
	log "github.com/sirupsen/logrus"
	"golang.org/x/mod/semver"
)

var (
	listenAddr   string
	listenPort   int
	sshUser      string
	sshPubkey    string
	k8sVersion   string
	cniVersion   string
	templatePath string
	configPath   string
)

func init() {
	flag.StringVar(&listenAddr, "address", "0.0.0.0", "listen address of the server")
	flag.IntVar(&listenPort, "port", 8080, "listen port of the server")
	flag.StringVar(&sshUser, "ssh-user", "core", "initial user available via ssh")
	flag.StringVar(&sshPubkey, "ssh-pubkey", "", "key to add to `authorized_hosts` for ssh access")
	flag.StringVar(&k8sVersion, "k8s-version", "v1.23.0", "Kubernetes version to install")
	flag.StringVar(&cniVersion, "cni-version", "v1.0.0", "CNI version to install")
	flag.StringVar(&templatePath, "template", "/etc/zippo/templates/config.tmpl", "path to Butane template to render")
	flag.StringVar(&configPath, "config", "/etc/zippo/config.yaml", "path to place rendered config")
	flag.Parse()
}

// Define your own args to pass into the template
type args struct {
	Hostname      string
	SSHUser       string
	SSHPubkey     string
	CNIVersion    string
	CRICtlVersion string
	K8sVersion    string
}

func main() {
	// Don't set `Hostname`, it's populated from the host header
	args := args{
		Hostname:   "{{ .Hostname }}",
		SSHUser:    sshUser,
		SSHPubkey:  sshPubkey,
		K8sVersion: k8sVersion,
		// CRICTL only has major/minor releases
		CRICtlVersion: semver.MajorMinor(k8sVersion) + ".0",
		CNIVersion:    cniVersion,
	}

	// populate config file with rendered template
	template, err := zippo.Render(templatePath, args)
	if err != nil {
		log.Fatalf(err.Error())
	}
	ioutil.WriteFile(configPath, template.Bytes(), 0600)

	zippo.NewServer(listenAddr, listenPort, configPath, args).Run()
}
