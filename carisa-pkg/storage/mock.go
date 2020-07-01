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

func (_m *ErrMockCRUD) Close() error {
	if _m.close {
		return errors.New("close")
	}
	return nil
}

func (_m *ErrMockCRUD) Create(entity Entity) (OpeWrap, error) {
	if _m.create {
		return OpeWrap{}, errors.New("create")
	}
	return OpeWrap{}, nil
}

// Activate activates the methods to throw a error
func (_m *ErrMockCRUD) Activate(methods ...string) {
	_m.Clear()

	for _, method := range methods {
		switch method {
		case "Create":
			_m.create = true
		case "Close":
			_m.close = true
		default:
			panic("method not found")
		}
	}
}

// Clear deactivates all methods
func (_m *ErrMockCRUD) Clear() {
	_m.create = false
	_m.close = false
}

// ErrMockTxn allows test the errors.
// For testing other functions use storage.Integration
type ErrMockTxn struct {
	commit bool
}

func (_m *ErrMockTxn) Commit(ctx context.Context) (bool, error) {
	if _m.commit {
		return false, errors.New("commit")
	}
	return true, nil
}

// Activate activates the methods to throw a error
func (_m *ErrMockTxn) Activate(methods ...string) {
	_m.Clear()

	for _, method := range methods {
		switch method {
		case "Commit":
			_m.commit = true
		default:
			panic("method not found")
		}
	}
}

// Clear deactivates all methods
func (_m *ErrMockTxn) Clear() {
	_m.commit = false
}

func (_m *ErrMockTxn) DoFound(ope OpeWrap) {
}

func (_m *ErrMockTxn) DoNotFound(ope OpeWrap) {
}

func (_m *ErrMockTxn) Find(keyValue string) {
}
