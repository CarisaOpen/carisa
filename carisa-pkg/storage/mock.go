package storage

import (
	"context"
	"errors"
)

// ErrMockCRUD allows test the errors.
// For testing other functions use storage.Integration
type ErrMockCRUD struct {
	create bool
	close  bool
	get    bool
}

func (e *ErrMockCRUD) Close() error {
	if e.close {
		return errors.New("close")
	}
	return nil
}

func (e *ErrMockCRUD) Put(entity Entity) (OpeWrap, error) {
	if e.create {
		return OpeWrap{}, errors.New("create")
	}
	return OpeWrap{}, nil
}

func (e *ErrMockCRUD) Get(ctx context.Context, key string, entity Entity) (bool, error) {
	if e.get {
		return false, errors.New("get")
	}
	return true, nil
}

// Activate activates the methods to throw a error
func (e *ErrMockCRUD) Activate(methods ...string) {
	e.Clear()

	for _, method := range methods {
		switch method {
		case "Put":
			e.create = true
		case "Close":
			e.close = true
		case "Get":
			e.get = true
		default:
			panic("method not found")
		}
	}
}

// Clear deactivates all methods
func (e *ErrMockCRUD) Clear() {
	e.create = false
	e.close = false
	e.get = false
}

// ErrMockTxn allows test the errors.
// For testing other functions use storage.Integration
type ErrMockTxn struct {
	commit bool
}

func (e *ErrMockTxn) Commit(ctx context.Context) (bool, error) {
	if e.commit {
		return false, errors.New("commit")
	}
	return true, nil
}

// Activate activates the methods to throw a error
func (e *ErrMockTxn) Activate(methods ...string) {
	e.Clear()

	for _, method := range methods {
		switch method {
		case "Commit":
			e.commit = true
		default:
			panic("method not found")
		}
	}
}

// Clear deactivates all methods
func (e *ErrMockTxn) Clear() {
	e.commit = false
}

func (e *ErrMockTxn) DoFound(ope OpeWrap) {
}

func (e *ErrMockTxn) DoNotFound(ope OpeWrap) {
}

func (e *ErrMockTxn) Find(keyValue string) {
}

// ErrMockCRUDOper allows test the errors.
type ErrMockCRUDOper struct {
	create bool
	put    bool
}

// Activate activates the methods to throw a error
func (e *ErrMockCRUDOper) Activate(methods ...string) {
	e.Clear()

	for _, method := range methods {
		switch method {
		case "Create":
			e.create = true
		case "Put":
			e.put = true
		default:
			panic("method not found")
		}
	}
}

// Clear deactivates all methods
func (e *ErrMockCRUDOper) Clear() {
	e.create = false
	e.put = false
}

func (e *ErrMockCRUDOper) Create(loc string, storeTimeout StoreWithTimeout, entity Entity) (bool, error) {
	if e.create {
		return false, errors.New("create")
	}
	return true, nil
}

func (e *ErrMockCRUDOper) Put(loc string, storeTimeout StoreWithTimeout, entity Entity) (bool, error) {
	if e.put {
		return false, errors.New("put")
	}
	return true, nil
}
