package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fabianofranz/origin/pkg/cmd/server"
	"github.com/fabianofranz/origin/pkg/cmd/server/etcd"
	"github.com/golang/glog"
	"github.com/openshift/origin/pkg/cmd/flagtypes"
	"github.com/openshift/origin/pkg/cmd/util/docker"
)

func main() {
	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	etcdMasterPort := 16001
	etcdBindPort := 17001
	etcdPeerBindPort := 18001

	etcdConfig := &etcd.Config{
		MasterAddr:   fmt.Sprintf("%s:%v", host, etcdMasterPort),
		BindAddr:     fmt.Sprintf("%s:%v", host, etcdBindPort),
		PeerBindAddr: fmt.Sprintf("%s:%v", host, etcdPeerBindPort),
	}
	etcdConfig.Run()

	etcdAddr := flagtypes.Addr{DefaultScheme: "http", DefaultPort: etcdBindPort}
	etcdAddr.Set(fmt.Sprintf("%s:%v", host, etcdBindPort))

	bind := fmt.Sprintf("%s:%v", host, port)

	cfg := &server.Config{
		Docker: docker.NewHelper(),

		MasterAddr:     flagtypes.Addr{Value: bind, DefaultScheme: "http", DefaultPort: port, AllowPrefix: true}.Default(),
		BindAddr:       flagtypes.Addr{Value: bind, DefaultScheme: "http", DefaultPort: port, AllowPrefix: true}.Default(),
		EtcdAddr:       etcdAddr,
		KubernetesAddr: flagtypes.Addr{Value: bind, DefaultScheme: "http", DefaultPort: port}.Default(),
		AssetAddr:      flagtypes.Addr{Value: bind, DefaultScheme: "http", DefaultPort: port}.Default(),
		PortalNet:      flagtypes.DefaultIPNet("172.30.17.0/24"),

		RunAssetServer: false,
		Hostname:       host,
		NodeList:       flagtypes.StringList{},
	}

	if err := server.Start(cfg, []string{"master"}); err != nil {
		glog.Fatal(err)
	}
}

