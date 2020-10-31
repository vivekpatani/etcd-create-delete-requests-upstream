module example.com/m/v2

go 1.15

require (
	go.etcd.io/etcd v3.3.20+incompatible
	go.uber.org/zap v1.14.1
)

replace go.etcd.io/etcd => go.etcd.io/etcd v0.0.0-20191023171146-3cf2f69b5738 // 3cf2f69b5738 is the SHA for git tag v3.4.3
