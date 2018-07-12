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

func RegisterBlobstore(key string, c BlobstoreConstructor) {
	if key == "" {
		panic(fmt.Sprintf("key cannot be empty"))
	}

	blobstoreRegistryMu.Lock()
	defer blobstoreRegistryMu.Unlock()

	if _, ok := blobstoreRegistry[key]; ok {
		panic(fmt.Sprintf("already registered blobstore: %s", key))
	}

	blobstoreRegistry[key] = c
}

func RegisterIndex(key string, c IndexConstructor) {
	if key == "" {
		panic(fmt.Sprintf("key cannot be empty"))
	}

	indexRegistryMu.Lock()
	defer indexRegistryMu.Unlock()

	if _, ok := indexRegistry[key]; ok {
		panic(fmt.Sprintf("already registered index: %s", key))
	}

	indexRegistry[key] = c
}

func RegisterStore(key string, c StoreConstructor) {
	if key == "" {
		panic(fmt.Sprintf("key cannot be empty"))
	}

	storeRegistryMu.Lock()
	defer storeRegistryMu.Unlock()

	if _, ok := storeRegistry[key]; ok {
		panic(fmt.Sprintf("already registered store: %s", key))
	}

	storeRegistry[key] = c
}

func (f BlobstoreConstructorFunc) New(name string, c config.Config) (Blobstore, error) {
	return f(name, c)
}

func (f IndexConstructorFunc) New(name string, c config.Config) (Index, error) {
	return f(name, c)
}

func (f StoreConstructorFunc) New(name string, c config.Config) (Store, error) {
	return f(name, c)
}
