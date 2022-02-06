package main

import (
	"flag"
	"io/ioutil"

	zippo "github.com/jan0ski/zippo/pkg"
	log "github.com/sirupsen/logrus"
)

var (
	listenAddr   string
	listenPort   int
	sshUser      string
	sshPubkey    string
	k8sVersion   string
	cniVersion   string
	templatePath string
)

func init() {
	flag.StringVar(&listenAddr, "address", "0.0.0.0", "listen address of the server")
	flag.IntVar(&listenPort, "port", 8080, "key to add to `authorized_hosts` for ssh access")
	flag.StringVar(&sshUser, "ssh-user", "core", "initial user available via ssh")
	flag.StringVar(&sshPubkey, "ssh-pubkey", "", "key to add to `authorized_hosts` for ssh access")
	flag.StringVar(&k8sVersion, "k8s", "v1.20.14", "Kubernetes version to install")
	flag.StringVar(&cniVersion, "cni-version", "v0.8.2", "CNI version to install")
	flag.StringVar(&templatePath, "template", "/etc/zippo/config.tmpl", "path to Butane template")
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
		Hostname:      "{{ .Hostname }}",
		SSHUser:       sshUser,
		SSHPubkey:     sshPubkey,
		K8sVersion:    k8sVersion,
		CRICtlVersion: k8sVersion,
		CNIVersion:    cniVersion,
	}

	// populate config file with rendered template
	configPath := "/etc/zippo/config.yaml"
	template, err := zippo.Render(templatePath, args)
	if err != nil {
		log.Fatalf(err.Error())
	}
	ioutil.WriteFile(configPath, template.Bytes(), 0600)

	zippo.NewServer(listenAddr, listenPort, configPath, args).Run()
}
