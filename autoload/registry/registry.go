package registry

import (
	"errors"

	"github.com/leeola/fixity"
	cu "github.com/leeola/fixity/util/configumarshaller"
)

var (
	// If the loaders are nil, index(es) have already been produced and the
	// loaders have had their memory freed.
	fixityLoaders []FixityLoader
	indexLoaders  []IndexLoader
	storeLoaders  []StoreLoader
)

func init() {
	// init the loaders so they're not nil. Nil loader slice represents a freed
	// slice of loaders.
	fixityLoaders := []FixityLoader{}
	indexLoaders := []IndexLoader{}
	storeLoaders := []StoreLoader{}
}

type FixityLoader func(cu.ConfigUnmarshaller) (fixity.Fixity, error)
type IndexLoader func(cu.ConfigUnmarshaller) (fixity.Index, error)
type StoreLoader func(cu.ConfigUnmarshaller) (fixity.Store, error)

func RegisterFixity(l FixityLoader) error {
	if loadedFixity != nil {
		return errors.New("fixity already loaded")
	}

	fixityLoaders = append(fixityLoaders, l)
	return nil
}

func RegisterIndex(l IndexLoader) error {
	if loadedFixity != nil {
		return errors.New("fixity already loaded")
	}

	indexLoaders = append(indexLoaders, l)
	return nil
}

func RegisterStore(l StoreLoader) error {
	if loadedFixity != nil {
		return errors.New("fixity already loaded")
	}

	storeLoaders = append(storeLoaders, l)
	return nil
}

// LoadFixity from the given configunmarshaller.
//
// Note that LoadFixity purges the registered fixities if successful.
func LoadFixity(cu cu.ConfigUnmarshaller) (fixity.Fixity, error) {
	if fixityLoaders == nil {
		return nil, errors.New("fixity already loaded")
	}

	if len(fixityLoaders) == 0 {
		return nil, errors.New("no fixities registered")
	}

	var fixities []fixity.Fixity
	for _, l := range fixityLoaders {
		i, err := l(cu)
		if err != nil {
			return nil, err
		}

		// the loader may not have actually loaded a value, so check the resulting fixity.
		if i != nil {
			fixities = append(fixities, i)
		}
	}

	if len(fixities) == 0 {
		return nil, errors.New("configuration does not define valid fixity")
	}

	// all fixity implementations only use a single fixity atm.
	if len(fixities) > 1 {
		return nil, errors.New("configuration defines multiple fixities")
	}

	// purge the slice to free memory. Only one autoload is allowed with
	// this usage.
	fixityLoaders = nil

	return fixities[0], nil
}

// LoadIndex from the given configunmarshaller.
//
// Note that LoadIndex purges the registered indexes if successful.
func LoadIndex(cu cu.ConfigUnmarshaller) (fixity.Index, error) {
	if indexLoaders == nil {
		return nil, errors.New("index already loaded")
	}

	if len(indexLoaders) == 0 {
		return nil, errors.New("no indexes registered")
	}

	var indexes []fixity.Index
	for _, l := range indexLoaders {
		i, err := l(cu)
		if err != nil {
			return nil, err
		}

		// the loader may not have actually loaded a value, so check the resulting index.
		if i != nil {
			indexes = append(indexes, i)
		}
	}

	if len(indexes) == 0 {
		return nil, errors.New("configuration does not define valid index")
	}

	// all fixity implementations only use a single index atm.
	if len(indexes) > 1 {
		return nil, errors.New("configuration defines multiple indexes")
	}

	// purge the slice to free memory. Only one autoload is allowed with
	// this usage.
	indexLoaders = nil

	return indexes[0], nil
}

// LoadStore from the given configunmarshaller.
//
// Note that LoadStore purges the registered stores if successful.
func LoadStore(cu cu.ConfigUnmarshaller) (fixity.Store, error) {
	if storeLoaders == nil {
		return nil, errors.New("store already loaded")
	}

	if len(storeLoaders) == 0 {
		return nil, errors.New("no stores registered")
	}

	var stores []fixity.Store
	for _, l := range storeLoaders {
		i, err := l(cu)
		if err != nil {
			return nil, err
		}

		// the loader may not have actually loaded a value, so check the resulting store.
		if i != nil {
			stores = append(stores, i)
		}
	}

	if len(stores) == 0 {
		return nil, errors.New("configuration does not define valid store")
	}

	// all fixity implementations only use a single store atm.
	if len(stores) > 1 {
		return nil, errors.New("configuration defines multiple stores")
	}

	// purge the slice to free memory. Only one autoload is allowed with
	// this usage.
	storeLoaders = nil

	return stores[0], nil
}
