package store

import (
	"encoding/gob"
	"os"
	"strings"
)

type Document struct {
	Filename string
	Page     int
	Content  string
}

type Store struct {
	index     map[string][]Document
	indexPath string
}

func New() *Store {
	return &Store{
		index:     make(map[string][]Document),
		indexPath: "pdf_index.gob",
	}
}

func (s *Store) AddDocument(filename string, page int, content string) {
	words := strings.Fields(strings.ToLower(content))
	for _, word := range words {
		s.index[word] = append(s.index[word], Document{
			Filename: filename,
			Page:     page,
			Content:  content,
		})
	}
}

func (s *Store) Search(terms []string) []Document {
	if len(terms) == 0 {
		return nil
	}

	// Get documents containing first term
	docs := s.index[terms[0]]

	// Filter for documents containing all terms
	results := make([]Document, 0)
	for _, doc := range docs {
		matches := true
		for _, term := range terms[1:] {
			if !strings.Contains(strings.ToLower(doc.Content), term) {
				matches = false
				break
			}
		}
		if matches {
			results = append(results, doc)
		}
	}

	return results
}

func (s *Store) Save() error {
	f, err := os.Create(s.indexPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return gob.NewEncoder(f).Encode(s.index)
}

func (s *Store) Load() error {
	f, err := os.Open(s.indexPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return gob.NewDecoder(f).Decode(&s.index)
}
