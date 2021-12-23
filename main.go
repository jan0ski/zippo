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
	templatePath string
)

func init() {
	flag.StringVar(&listenAddr, "address", "0.0.0.0", "listen address of the server")
	flag.IntVar(&listenPort, "port", 8080, "key to add to `authorized_hosts` for ssh access")
	flag.StringVar(&sshUser, "ssh-user", "core", "initial user available via ssh")
	flag.StringVar(&sshPubkey, "ssh-pubkey", "", "key to add to `authorized_hosts` for ssh access")
	flag.StringVar(&k8sVersion, "k8s", "v1.20.14", "Kubernetes version to install")
	flag.StringVar(&templatePath, "template", "/etc/zippo/template.yaml", "path to Butane template")
	flag.Parse()
}

// Define your own args to pass into the template
// Don't include `Hostname`, its populated from the host header
type args struct {
	SSHUser       string
	SSHPubkey     string
	CNIVersion    string
	CRICtlVersion string
	K8sVersion    string
}

func main() {
	args := args{
		SSHUser:    sshUser,
		SSHPubkey:  sshPubkey,
		K8sVersion: k8sVersion,
	}

	// overwrite config file with rendered variables
	template, err := zippo.Render(templatePath, args)
	if err != nil {
		log.Fatalf(err.Error())
	}
	ioutil.WriteFile(templatePath, template.Bytes(), 0600)

	zippo.NewServer(listenAddr, listenPort, templatePath, args).Run()
}
