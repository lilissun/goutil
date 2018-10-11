package goutil

import (
	"math"
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
			return documents[i].Score < documents[j].Score
		},
	)
	for queue.Pop() {
	}
	var lastdoc *Document
	for index, document := range documents {
		t.Logf("[%02d] Document (%d) Score (%f)",
			index+1, document.ID, document.Score)
		if lastdoc != nil && lastdoc.Score < document.Score {
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
			return documents[i].Score < documents[j].Score
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
		if lastdoc != nil && lastdoc.Score < document.Score {
			t.Errorf("Document (%d) and (%d) are in Wrong Order",
				lastdoc.ID, document.ID)
			return
		}
		lastdoc = document
	}
}

func TestPriorityTopKth(t *testing.T) {
	var orders []*Document
	for _, partition := range elastic {
		for id, score := range partition {
			orders = append(orders, NewDocument(id, score))
		}
	}
	sort.Slice(orders, func(i, j int) bool { return orders[i].Score > orders[j].Score })
	for _, k := range []int{1, 2, 4, 5, 7, 8} {
		var documents []*Document
		for _, partition := range elastic {
			for id, score := range partition {
				documents = append(documents, NewDocument(id, score))
			}
		}
		queue := NewPriorityQueue(
			documents[:k], k,
			func(i, j int) bool {
				return documents[i].Score < documents[j].Score
			},
		)
		for _, document := range documents[k:] {
			if document.Score > documents[0].Score {
				documents[0] = document
				queue.Fix(0)
			}
		}
		if documents[0].ID != orders[k-1].ID {
			t.Errorf("The Top (%d) Document does not have ID (%d)", k, orders[k].ID)
			return
		}
	}
}

type Documents struct {
	Documents []*Document
	Queue     *PriorityQueue
}

func NewDocuments(documents ...*Document) *Documents {
	docs := &Documents{Documents: documents}
	docs.Queue = NewPriorityQueue(
		docs.Documents, len(docs.Documents),
		func(i, j int) bool { return docs.Documents[i].Score > docs.Documents[j].Score },
	)
	return docs
}

func (docs *Documents) Push(doc *Document) {
	if docs.Queue.Length() < docs.Queue.Capacity() {
		docs.Documents[docs.Queue.Length()] = doc
		docs.Queue.Push()
		return
	}
	documents := docs.Documents
	length := docs.Queue.Length()
	docs.Documents = make([]*Document, 2*length+1)
	copy(docs.Documents, documents)
	docs.Documents[length] = doc
	docs.Queue = NewPriorityQueue(
		docs.Documents, length+1,
		func(i, j int) bool { return docs.Documents[i].Score > docs.Documents[j].Score },
	)
}

func (docs *Documents) Pop() *Document {
	if docs.Queue.Pop() {
		return docs.Documents[docs.Queue.Length()]
	}
	return nil
}

func (docs *Documents) Top() *Document {
	if len(docs.Documents) != 0 {
		return docs.Documents[0]
	}
	return nil
}

func (docs *Documents) Update(doc *Document, index int) bool {
	if index >= 0 && index < docs.Queue.Length() {
		docs.Documents[index] = doc
		docs.Queue.Fix(index)
		return true
	}
	return false
}

func TestPriorityRealPushSort(t *testing.T) {
	var docs *Documents
	for index, partition := range elastic {
		if index == 0 {
			var documents []*Document
			for id, score := range partition {
				documents = append(documents, NewDocument(id, score))
			}
			docs = NewDocuments(documents...)
		} else {
			for id, score := range partition {
				docs.Push(NewDocument(id, score))
			}
		}
	}
	var lastdoc *Document
	for {
		document := docs.Pop()
		if document == nil {
			break
		}
		t.Logf("Document (%d) Score (%f)", document.ID, document.Score)
		if lastdoc != nil && lastdoc.Score < document.Score {
			t.Errorf("Document (%d) and (%d) are in Wrong Order",
				lastdoc.ID, document.ID)
			return
		}
		lastdoc = document
	}
}

type Batch struct {
	Documents []*Document
	Index     int
}

func NewBatch(docs ...*Document) *Batch {
	return &Batch{
		Index:     0,
		Documents: docs,
	}
}

func (batch *Batch) GetDocument() *Document {
	if batch.IsEmpty() {
		return nil
	}
	return batch.Documents[batch.Index]
}

func (batch *Batch) GetScore() float64 {
	if batch.IsEmpty() {
		return math.Inf(-1)
	}
	return batch.Documents[batch.Index].Score
}

func (batch *Batch) IsEmpty() bool {
	return batch.Index == len(batch.Documents)
}

func (batch *Batch) MoveToNext() {
	batch.Index++
}

func TestPriorityKWayMerge(t *testing.T) {
	batches := make([]*Batch, 0, len(elastic))
	var count int
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
		count += len(documents)
		batches = append(batches, NewBatch(documents...))
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
		count--
		t.Logf("Document (%d) Score (%f)", document.ID, document.Score)
		if lastdoc != nil && lastdoc.Score < document.Score {
			t.Errorf("Document (%d) and (%d) are in Wrong Order",
				lastdoc.ID, document.ID)
			return
		}
		lastdoc = document
		documents.MoveToNext()
		queue.Fix(0)
	}
	if count != 0 {
		t.Errorf("Documents are not Fully Listed, left: %d", count)
		return
	}
}
