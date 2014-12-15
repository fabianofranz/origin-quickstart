package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/golang/glog"
	"github.com/openshift/origin/pkg/cmd/flagtypes"
	"github.com/openshift/origin/pkg/cmd/util/docker"
)

type config struct {
	Docker *docker.Helper

	MasterAddr     flagtypes.Addr
	BindAddr       flagtypes.Addr
	EtcdAddr       flagtypes.Addr
	KubernetesAddr flagtypes.Addr
	PortalNet      flagtypes.IPNet

	Hostname  string
	VolumeDir string

	EtcdDir string

	StorageVersion string

	NodeList flagtypes.StringList

	CORSAllowedOrigins    flagtypes.StringList
	RequireAuthentication bool
}

func main() {
	hostname, err := defaultHostname()
	if err != nil {
		hostname = "localhost"
		glog.Warningf("Unable to lookup hostname, using %q: %v", hostname, err)
	}

	cfg := &config{
		Docker: docker.NewHelper(),

		MasterAddr:     flagtypes.Addr{Value: "localhost:8080", DefaultScheme: "http", DefaultPort: 8080, AllowPrefix: true}.Default(),
		BindAddr:       flagtypes.Addr{Value: "0.0.0.0:8080", DefaultScheme: "http", DefaultPort: 8080, AllowPrefix: true}.Default(),
		EtcdAddr:       flagtypes.Addr{Value: "0.0.0.0:4001", DefaultScheme: "http", DefaultPort: 4001}.Default(),
		KubernetesAddr: flagtypes.Addr{DefaultScheme: "http", DefaultPort: 8080}.Default(),
		PortalNet:      flagtypes.DefaultIPNet("172.30.17.0/24"),

		Hostname: hostname,
		NodeList: flagtypes.StringList{"127.0.0.1"},
	}

	if err := Start(cfg, []string{}); err != nil {
		glog.Fatal(err)
	}
}

// defaultHostname returns the default hostname for this system.
func defaultHostname() (string, error) {
	// Note: We use exec here instead of os.Hostname() because we
	// want the FQDN, and this is the easiest way to get it.
	fqdn, err := exec.Command("hostname", "-f").Output()
	if err != nil {
		return "", fmt.Errorf("Couldn't determine hostname: %v", err)
	}
	return strings.TrimSpace(string(fqdn)), nil
}
