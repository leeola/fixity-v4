package fixity

import (
	"fmt"
	"sync"

	"github.com/leeola/fixity/config"
)

var (
	blobstoreRegistry   map[string]BlobstoreConstructor
	blobstoreRegistryMu sync.Mutex
	indexRegistry       map[string]IndexConstructor
	indexRegistryMu     sync.Mutex
	storeRegistry       map[string]StoreConstructor
	storeRegistryMu     sync.Mutex
)

func init() {
	blobstoreRegistry = map[string]BlobstoreConstructor{}
	indexRegistry = map[string]IndexConstructor{}
	storeRegistry = map[string]StoreConstructor{}
}

type BlobstoreConstructor interface {
	New(name string, c config.Config) (Blobstore, error)
}

type IndexConstructor interface {
	New(name string, c config.Config) (Index, error)
}

type StoreConstructor interface {
	New(name string, c config.Config) (Store, error)
}

type BlobstoreConstructorFunc func(string, config.Config) (Blobstore, error)

type IndexConstructorFunc func(string, config.Config) (Index, error)

type StoreConstructorFunc func(string, config.Config) (Store, error)

func RegisterBlobstore(blobstoreType string, c BlobstoreConstructor) {
	if blobstoreType == "" {
		panic(fmt.Sprintf("blobstoreType cannot be empty"))
	}

	blobstoreRegistryMu.Lock()
	defer blobstoreRegistryMu.Unlock()

	if _, ok := blobstoreRegistry[blobstoreType]; ok {
		panic(fmt.Sprintf("already registered blobstore: %s", blobstoreType))
	}

	blobstoreRegistry[blobstoreType] = c
}

func RegisterIndex(indexType string, c IndexConstructor) {
	if indexType == "" {
		panic(fmt.Sprintf("indexType cannot be empty"))
	}

	indexRegistryMu.Lock()
	defer indexRegistryMu.Unlock()

	if _, ok := indexRegistry[indexType]; ok {
		panic(fmt.Sprintf("already registered index: %s", indexType))
	}

	indexRegistry[indexType] = c
}

func RegisterStore(storeType string, c StoreConstructor) {
	if storeType == "" {
		panic(fmt.Sprintf("storeType cannot be empty"))
	}

	storeRegistryMu.Lock()
	defer storeRegistryMu.Unlock()

	if _, ok := storeRegistry[storeType]; ok {
		panic(fmt.Sprintf("already registered store: %s", storeType))
	}

	storeRegistry[storeType] = c
}

func (f BlobstoreConstructorFunc) New(n string, c config.Config) (Blobstore, error) {
	return f(n, c)
}

func (f IndexConstructorFunc) New(n string, c config.Config) (Index, error) {
	return f(n, c)
}

func (f StoreConstructorFunc) New(n string, c config.Config) (Store, error) {
	return f(n, c)
}
