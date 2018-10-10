package goutil

import (
	"encoding/json"
	"sort"
	"testing"
)

type Document struct {
	ID    int64
	Score float64
}

func NewDocument(id int64, score float64) *Document {
	return &Document{
		ID:    id,
		Score: score,
	}
}

type Documents struct {
	Documents []*Document
	Index     int
}

func NewDocuments(docs ...*Document) *Documents {
	return &Documents{
		Index:     0,
		Documents: docs,
	}
}

func (docs *Documents) GetDocument() *Document {
	if docs.IsEmpty() {
		return nil
	}
	return docs.Documents[docs.Index]
}

func (docs *Documents) GetScore() float64 {
	if docs.IsEmpty() {
		return 0
	}
	return docs.Documents[docs.Index].Score
}

func (docs *Documents) IsEmpty() bool {
	return docs.Index == len(docs.Documents)
}

func (docs *Documents) Next() {
	docs.Index++
}

var (
	elastic = []map[int64]float64{
		map[int64]float64{
			101: 9.2,
			102: 6.6,
			103: 5.7,
			104: 4.3,
		},
		map[int64]float64{
			201: 9.1,
			202: 6.5,
			203: 5.6,
			204: 4.2,
		},
		map[int64]float64{
			301: 9.3,
			302: 6.7,
			303: 5.8,
			304: 4.4,
		},
	}
)

func TestElasticDocuments(t *testing.T) {
	batches := make([]*Documents, 0, len(elastic))
	for _, partition := range elastic {
		documents := make([]*Document, 0, len(partition))
		for id, score := range partition {
			documents = append(documents, NewDocument(id, score))
		}
		sort.Slice(
			documents,
			func(i, j int) bool {
				return documents[i].Score > documents[j].Score
			},
		)
		batches = append(batches, NewDocuments(documents...))
	}
	queue := NewPriorityQueue(
		batches, len(batches),
		func(i, j int) bool {
			return batches[i].GetScore() > batches[j].GetScore()
		},
	)
	var lastdoc *Document
	for {
		documents := batches[0]
		document := documents.GetDocument()
		if document == nil {
			t.Logf("No More Document")
			break
		}
		t.Logf("Document (%d) Score (%f)", document.ID, document.Score)
		if lastdoc != nil && lastdoc.Score < document.Score {
			t.Errorf("Document (%d) and (%d) are in Wrong Order",
				lastdoc.ID, document.ID)
			return
		}
		lastdoc = document
		documents.Next()
		queue.Fix(0)
	}
	data, err := json.Marshal(batches)
	if err != nil {
		t.Errorf("JSON Marshal Error (%s)", err)
		return
	}
	t.Logf("%s", data)
}

func TestPrioritySort(t *testing.T) {
	var documents []*Document
	for _, partition := range elastic {
		for id, score := range partition {
			documents = append(documents, NewDocument(id, score))
		}
	}
	queue := NewPriorityQueue(
		documents, len(documents),
		func(i, j int) bool {
			return documents[i].Score > documents[j].Score
		},
	)
	for queue.Pop() {
	}
	var lastdoc *Document
	for index, document := range documents {
		t.Logf("[%02d] Document (%d) Score (%f)",
			index+1, document.ID, document.Score)
		if lastdoc != nil && lastdoc.Score > document.Score {
			t.Errorf("Document (%d) and (%d) are in Wrong Order",
				lastdoc.ID, document.ID)
			return
		}
		lastdoc = document
	}
}

func TestPriorityPushSort(t *testing.T) {
	var documents []*Document
	for _, partition := range elastic {
		for id, score := range partition {
			documents = append(documents, NewDocument(id, score))
		}
	}
	queue := NewPriorityQueue(
		documents, len(documents)/2,
		func(i, j int) bool {
			return documents[i].Score > documents[j].Score
		},
	)
	for queue.Push() {
	}
	for queue.Pop() {
	}
	var lastdoc *Document
	for index, document := range documents {
		t.Logf("[%02d] Document (%d) Score (%f)",
			index+1, document.ID, document.Score)
		if lastdoc != nil && lastdoc.Score > document.Score {
			t.Errorf("Document (%d) and (%d) are in Wrong Order",
				lastdoc.ID, document.ID)
			return
		}
		lastdoc = document
	}
}
