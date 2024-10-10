package storage

import (
	"golang-island/internal/data"
	"testing"
	"time"
)

var id = 0
var toPut = data.Data{IslandCount: 33, CalculationDate: time.Now()}

func TestMockRepository_Put(t *testing.T) {
	var mockRepo = NewRepo(nil)

	mockRepo.Put(id, toPut)
	success := mockRepo.PutIfEmpty(id, toPut)

	if success {
		t.Error("PutIfEmpty should have returned false")
	}
}

func TestMockRepository_ExistsById(t *testing.T) {
	var mockRepo = NewRepo(nil)
	mockRepo.Put(id, toPut)

	actual, err := mockRepo.GetById(id)

	if err != nil {
		t.Errorf("Unexpected error, %s", err.Error())
	} else if actual != toPut {
		t.Errorf("Unexpected result")
	}
}

func TestMockRepository_GetById(t *testing.T) {
	var mockRepo = NewRepo(nil)
	mockRepo.Put(id, toPut)

	actual, err := mockRepo.GetById(id)

	if err != nil {
		t.Errorf("Unexpected error, %s", err.Error())
	} else if actual != toPut {
		t.Errorf("Unexpected result")
	}
}

func TestMockRepository_(t *testing.T) {
	var mockRepo = NewRepo(nil)
	mockRepo.Put(id, toPut)

	actual, err := mockRepo.GetById(id)

	if err != nil {
		t.Errorf("Unexpected error, %s", err.Error())
	} else if actual != toPut {
		t.Errorf("Unexpected result")
	}
}

func TestMockRepository_GetById2(t *testing.T) {
	var mockRepo = NewRepo(nil)
	mockRepo.Put(id, toPut)

	actual, err := mockRepo.GetById(id + 1)
	expected := data.Data{}

	if err == nil {
		t.Errorf("No error found but expected")
	} else if actual != expected {
		t.Errorf("Unexpected result, expected empty data")
	}
}
