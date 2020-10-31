package main

import (
	"go.etcd.io/etcd/clientv3"
	"math/rand"
	"time"
)

// endpoints connection string
const (
	defaultEndpoints = "localhost:2379"
	// dialTimeout sets an expiration
	defaultDialTimeout = 4 * time.Second
)

func generator(size int) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, size)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return "foo/" + string(b)
}

// createMCli creates a new maintenance Client
func createMCli() clientv3.Maintenance {
	if cli != nil {
		return clientv3.NewMaintenance(cli)
	} else {
		return clientv3.NewMaintenance(createCli(false))
	}
}

// createKVCli creates a new KV Client
func createKVCli() clientv3.KV {
	if cli != nil {
		return clientv3.NewKV(cli)
	} else {
		return clientv3.NewKV(createCli(false))
	}
}

// createCli creates a new client
func createCli(force bool) *clientv3.Client {
	// Create a cli and return
	if cli == nil || force {
		var err error
		cli, err = clientv3.New(clientv3.Config{
			DialTimeout: defaultDialTimeout,
			Endpoints:   []string{endpoints},
		})
		if cli == nil || err != nil {
			log.Error("error creating etcd client")
		}
	}

	// return
	return cli
}