package main

import (
	"context"
	"flag"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"time"
)

const (
	// defaultDisableCreate disables creates in test
	defaultDisableCreate = false
	// defaultDisableDelete disables deletes in test
	defaultDisableDelete = false
	// defaultDisableCompaction disables compaction in test
	defaultDisableCompaction = false
	// defaultDisableDefragmentation disables defragmentation in test
	defaultDisableDefragmentation = false
	// defaultDisableCleanup disables cleanup in test
	defaultDisableCleanup = false
	// defaultKeyCount default number of keys to be created, deleted, compacted
	defaultKeyCount = 200
	// defaultKeySize default key size
	defaultKeySize = 102400
	// defaultFrequency default frequency
	defaultFrequency = 1
)

var (
	// log logger
	log *zap.SugaredLogger
	// cli etcd client
	cli *clientv3.Client
	// endpoints endpoints to connect to etcd
	endpoints string
	// disableCreate disables creates in test
	disableCreate bool
	// disableDelete disables deletes in test
	disableDelete bool
	// disableCompaction disables compaction in test
	disableCompaction bool
	// disableDefragmentation disables defragmentation in test
	disableDefragmentation bool
	// disableCleanup disables cleanup in test
	disableCleanup bool
	// keyCount number of keys to created, deleted, compacted
	keyCount int
	// keySize size of each key created
	keySize int
	// frequency of each operation
	frequency int
)

// main does main things
func main() {
	// 1. define logger
	log = zap.NewExample().Sugar()
	defer log.Sync()

	// 2. init all passed flags and create KV/Maintenance Client
	flagInit()
	kv := createKVCli()
	//m := createMCli()
	ticker := time.NewTicker(2 * time.Second)
	var rev int64 = 0

	// 3. Generate a keyList to be used
	var keyList []string
	for i := 0; i < keyCount; i++ {
		keyList = append(keyList, generator(keySize))
	}

	// 3. Run once or forever
	// 3a. Create n keys if enabled, run creation for * frequency
	if !disableCreate {
		// This index is to keep track of index while running every second
		for i := 0; i < frequency; i++ {
			select {
			case <-ticker.C:
				for j := 0; j < keyCount; j++ {
					key := keyList[j]
					resp, err := kv.Put(context.TODO(), key, "")
					if err != nil {
						log.Errorw("failed to PUT key", "key", key, "error", err)
					}
					if resp != nil && resp.Header != nil {
						rev = resp.Header.Revision
					}
					log.Debug(resp)
				}
			}
		}
		log.Infow("keys created", "count", keyCount, "size of each key", keySize, "revision", rev)
	}

	time.Sleep(5 * time.Second)

	// 3b. Delete n keys if enabled
	if !disableDelete {
		deleted := 0
		var rev int64 = 0
		for i := 0; i < frequency; i++ {
			select {
			case <-ticker.C:
				for j := 0; j < keyCount; j++ {
					key := keyList[j]
					resp, err := kv.Delete(context.TODO(), key)
					if err != nil {
						log.Errorw("failed to DELETE key", "key", key, "error", err)
					}
					if resp != nil && resp.Header != nil {
						rev = resp.Header.Revision
					}
					log.Debug(resp)
					if resp.Deleted == 1 {
						deleted += 1
					} else {
						log.Warnw("deleted 0 or multiple", "deleted ", resp.Deleted)
					}
				}
			}
		}
		log.Infow("keys created", "count ", keyCount, "size of each key ", keySize, "revision ", rev, "deleted ", deleted)
	}
}

// flagInit inits all flags
func flagInit() {
	flag.StringVar(&endpoints, "endpoints", defaultEndpoints, "etcd endpoints")
	flag.BoolVar(&disableCreate, "disableCreate", defaultDisableCreate, "disable creates")
	flag.BoolVar(&disableDelete, "disableDelete", defaultDisableDelete, "disable deletes")
	flag.BoolVar(&disableCompaction, "disableCompaction", defaultDisableCompaction, "disable compaction")
	flag.BoolVar(&disableDefragmentation, "disableDefrag", defaultDisableDefragmentation, "disable defragmentation")
	flag.BoolVar(&disableCleanup, "disableCleanup", defaultDisableCleanup, "disable cleanup")
	flag.IntVar(&keyCount, "keyCount", defaultKeyCount, "number of keys to be created, deleted, compacted")
	flag.IntVar(&keySize, "keySize", defaultKeySize, "size of each key created")
	flag.IntVar(&frequency, "frequency", defaultFrequency, "frequency of each operation")
	flag.Parse()
	log.Debug("Flag status")
	log.Debug("endpoint: ", endpoints)
	log.Debug("creates disabled: ", disableCreate)
	log.Debug("deletes disabled: ", disableDelete)
	log.Debug("compaction disabled: ", disableCompaction)
	log.Debug("defragmentation disabled: ", disableDefragmentation)
	log.Debug("cleanup disabled: ", disableCleanup)
	log.Debug("number of keys to be created: ", keyCount*frequency)
}
