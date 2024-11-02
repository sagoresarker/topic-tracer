package indexer

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/sagoresarker/topic-tracer/internal/store"

	"github.com/ledongthuc/pdf"
)

type Indexer struct {
	store *store.Store
	mu    sync.Mutex
}

func New() *Indexer {
	return &Indexer{
		store: store.New(),
	}
}

func (i *Indexer) IndexDirectory(dir string) error {
	var wg sync.WaitGroup
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".pdf" {
			continue
		}

		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			i.indexFile(filepath.Join(dir, filename))
		}(file.Name())
	}

	wg.Wait()
	return i.store.Save()
}

func (i *Indexer) indexFile(path string) error {
	f, r, err := pdf.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for pageNum := 1; pageNum <= r.NumPage(); pageNum++ {
		p := r.Page(pageNum)
		text, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}

		i.mu.Lock()
		i.store.AddDocument(path, pageNum, text)
		i.mu.Unlock()
	}

	return nil
}
