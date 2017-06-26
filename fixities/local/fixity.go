package local

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"sync"

	"github.com/fatih/structs"
	"github.com/inconshreveable/log15"
	"github.com/leeola/errors"
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/chunkers/restic"
	"github.com/leeola/fixity/q"
)

var (
	blockMetaBucketKey = []byte("blockMeta")
	idsBucketKey       = []byte("ids")
	lastBlockKey       = []byte("lastBlock")
)

type Config struct {
	Index    fixity.Index `toml:"-"`
	Store    fixity.Store `toml:"-"`
	Log      log15.Logger `toml:"-"`
	RootPath string       `toml:"rootPath"`
}

type Fixity struct {
	config     Config
	blockchain *Blockchain
	db         Db
	idLock     *sync.Mutex
	index      fixity.Index
	store      fixity.Store
	log        log15.Logger
}

func New(c Config) (*Fixity, error) {
	if c.RootPath == "" {
		return nil, errors.New("missing required config: rootPath")
	}

	if c.Index == nil {
		return nil, errors.New("missing reqired config: Index")
	}

	if c.Store == nil {
		return nil, errors.New("missing reqired config: Store")
	}

	if c.Log == nil {
		c.Log = log15.New()
	}

	db, err := newBoltDb(c.RootPath)
	if err != nil {
		return nil, err
	}

	return &Fixity{
		config:     c,
		blockchain: NewBlockchain(c.Log, db, c.Store),
		db:         db,
		idLock:     &sync.Mutex{},
		index:      c.Index,
		store:      c.Store,
		log:        c.Log,
	}, nil
}

// loadPreviousInfo is a helper to load the hash and the chunksize of the
// previous content. Empty values are returned if no id is found.
func (l *Fixity) loadPreviousInfo(id string) (string, uint64, error) {
	c, err := l.Read(id)
	if err == fixity.ErrIdNotFound {
		return "", 0, nil
	}
	if err != nil {
		return "", 0, err
	}

	b, err := c.Blob()
	if err != nil {
		return "", 0, err
	}

	return c.Hash, b.AverageChunkSize, nil
}

func (l *Fixity) Blob(h string) (io.ReadCloser, error) {
	return l.store.Read(h)
}

func (l *Fixity) Blockchain() fixity.Blockchain {
	return l.blockchain
}

func (f *Fixity) Close() error {
	return f.db.Close()
}

func (l *Fixity) Delete(id string) error {
	c, err := l.Read(id)
	if err != nil {
		return err
	}

	cs := []fixity.Content{c}
	for c.PreviousContentHash != "" {
		c, err = c.Previous()
		if err != nil {
			return err
		}
		cs = append(cs, c)
	}

	_, err = l.Blockchain().DeleteContent(cs...)
	return err
}

func (l *Fixity) Search(q *q.Query) ([]string, error) {
	return l.index.Search(q)
}

func (l *Fixity) ReadHash(h string) (fixity.Content, error) {
	var c fixity.Content
	if err := ReadAndUnmarshal(l.store, h, &c); err != nil {
		return fixity.Content{}, err
	}

	if structs.IsZero(c) {
		return fixity.Content{}, fixity.ErrNotContent
	}

	c.Hash = h
	c.Store = l.store

	return c, nil
}

func (l *Fixity) Read(id string) (fixity.Content, error) {
	h, err := l.db.GetIdHash(id)
	if err != nil {
		return fixity.Content{}, err
	}

	return l.ReadHash(h)
}

func (l *Fixity) Write(id string, r io.Reader, f ...fixity.Field) (fixity.Content, error) {
	req := fixity.NewWrite(id, ioutil.NopCloser(r))
	req.Fields = f
	return l.WriteRequest(req)
}

func (l *Fixity) WriteRequest(req *fixity.WriteRequest) (fixity.Content, error) {
	if req.Blob == nil {
		return fixity.Content{}, errors.New("no data given to write")
	}
	defer req.Blob.Close()

	averageChunkSize := req.AverageChunkSize
	var previousContentHash string
	if req.Id != "" {
		l.idLock.Lock()
		defer l.idLock.Unlock()

		pch, acs, err := l.loadPreviousInfo(req.Id)
		if err != nil {
			return fixity.Content{}, err
		}
		previousContentHash = pch
		averageChunkSize = acs
	}

	if averageChunkSize == 0 {
		averageChunkSize = fixity.DefaultAverageChunkSize
	}

	chunker, err := restic.New(req.Blob, averageChunkSize)
	if err != nil {
		return fixity.Content{}, err
	}

	cHashes, totalSize, err := WriteChunker(l.store, chunker)
	if err != nil {
		return fixity.Content{}, err
	}

	blob := fixity.Blob{
		ChunkHashes:      cHashes,
		Size:             totalSize,
		AverageChunkSize: req.AverageChunkSize,
	}

	blobHash, err := MarshalAndWrite(l.store, blob)
	if err != nil {
		return fixity.Content{}, err
	}

	content := fixity.Content{
		Id:                  req.Id,
		PreviousContentHash: previousContentHash,
		BlobHash:            blobHash,
		IndexedFields:       req.Fields,
	}

	cHash, err := MarshalAndWrite(l.store, content)
	if err != nil {
		return fixity.Content{}, err
	}
	content.Store = l.store
	content.Hash = cHash

	// TODO(leeola): return the block instead of hashes directly.
	if _, err := l.Blockchain().AppendContent(content); err != nil {
		return fixity.Content{}, err
	}

	// if the id was supplied, update the new id
	if req.Id != "" {
		if err := l.db.SetIdHash(req.Id, cHash); err != nil {
			return fixity.Content{}, err
		}
	}

	// TODO(leeola): move this to a goroutine, no reason to
	// block writes while we index in the background.
	if err := l.index.Index(cHash, content.Id, req.Fields); err != nil {
		return fixity.Content{}, err
	}

	return content, nil
}

// WriteReader writes the given reader's content to the store.
func WriteReader(s fixity.Store, r io.Reader) (string, error) {
	if s == nil {
		return "", errors.New("Store is nil")
	}
	if r == nil {
		return "", errors.New("Reader is nil")
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", errors.Wrap(err, "failed to readall")
	}

	h, err := s.Write(b)
	return h, errors.Wrap(err, "store failed to write")
}

// MarshalAndWrite marshals the given interface to json and writes that to the store.
func MarshalAndWrite(s fixity.Store, v interface{}) (string, error) {
	if s == nil {
		return "", errors.New("Store is nil")
	}
	if v == nil {
		return "", errors.New("Interface is nil")
	}

	b, err := json.Marshal(v)
	if err != nil {
		return "", errors.Stack(err)
	}

	h, err := s.Write(b)
	if err != nil {
		return "", errors.Stack(err)
	}

	return h, nil
}

func ReadAll(s fixity.Store, h string) ([]byte, error) {
	rc, err := s.Read(h)
	if err != nil {
		return nil, errors.Stack(err)
	}
	defer rc.Close()

	return ioutil.ReadAll(rc)
}

func ReadAndUnmarshal(s fixity.Store, h string, v interface{}) error {
	_, err := ReadAndUnmarshalWithBytes(s, h, v)
	return err
}

func ReadAndUnmarshalWithBytes(s fixity.Store, h string, v interface{}) ([]byte, error) {
	b, err := ReadAll(s, h)
	if err != nil {
		return nil, errors.Stack(err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return nil, errors.Stack(err)
	}

	return b, nil
}

func WriteChunker(s fixity.Store, r fixity.Chunker) ([]string, int64, error) {
	var totalSize int64
	var hashes []string
	for {
		c, err := r.Chunk()
		if err != nil && err != io.EOF {
			return nil, 0, err
		}

		totalSize += c.Size

		if err == io.EOF {
			break
		}

		h, err := MarshalAndWrite(s, c)
		if err != nil {
			return nil, 0, err
		}
		hashes = append(hashes, h)
	}
	return hashes, totalSize, nil
}
