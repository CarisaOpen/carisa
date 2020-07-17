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
}

func (e *ErrMockCRUD) Close() error {
	if e.close {
		return errors.New("close")
	}
	return nil
}

func (e *ErrMockCRUD) Create(entity Entity) (OpeWrap, error) {
	if e.create {
		return OpeWrap{}, errors.New("create")
	}
	return OpeWrap{}, nil
}

// Activate activates the methods to throw a error
func (e *ErrMockCRUD) Activate(methods ...string) {
	e.Clear()

	for _, method := range methods {
		switch method {
		case "Create":
			e.create = true
		case "Close":
			e.close = true
		default:
			panic("method not found")
		}
	}
}

// Clear deactivates all methods
func (e *ErrMockCRUD) Clear() {
	e.create = false
	e.close = false
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
}

// Activate activates the methods to throw a error
func (e *ErrMockCRUDOper) Activate(methods ...string) {
	e.Clear()

	for _, method := range methods {
		switch method {
		case "Create":
			e.create = true
		default:
			panic("method not found")
		}
	}
}

// Clear deactivates all methods
func (e *ErrMockCRUDOper) Clear() {
	e.create = false
}

func (e *ErrMockCRUDOper) Create(loc string, storeTimeout StoreWithTimeout, entity Entity) (bool, error) {
	if e.create {
		return false, errors.New("create")
	}
	return true, nil
}
