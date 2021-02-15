package storage

import (
	"context"
	"errors"
)

// ErrMockCRUD allows test the errors.
// For testing other functions use storage.Integration
type ErrMockCRUD struct {
	put      bool
	putRaw   bool
	remove   bool
	get      bool
	getRaw   bool
	exists   bool
	startKey bool
	rang     bool
	rangRaw  bool
	close    bool
}

func (e *ErrMockCRUD) Close() error {
	if e.close {
		return errors.New("close")
	}
	return nil
}

func (e *ErrMockCRUD) Put(entity Entity) (OpeWrap, error) {
	if e.put {
		return OpeWrap{}, errors.New("put")
	}
	return OpeWrap{}, nil
}

func (e *ErrMockCRUD) PutRaw(key string, value string) OpeWrap {
	return OpeWrap{}
}

func (e *ErrMockCRUD) Get(ctx context.Context, key string, entity Entity) (bool, error) {
	if e.get {
		return false, errors.New("get")
	}
	return true, nil
}

func (e *ErrMockCRUD) GetRaw(ctx context.Context, key string) (bool, string, error) {
	if e.get {
		return false, "", errors.New("getRaw")
	}
	return true, "", nil
}

func (e *ErrMockCRUD) Remove(key string) OpeWrap {
	return OpeWrap{}
}

func (e *ErrMockCRUD) Exists(ctx context.Context, key string) (bool, error) {
	if e.exists {
		return false, errors.New("exists")
	}
	return true, nil
}

func (e *ErrMockCRUD) StartKey(ctx context.Context, key string, top int, empty func() Entity) ([]Entity, error) {
	if e.startKey {
		return nil, errors.New("startKey")
	}
	list := make([]Entity, top)
	return list, nil
}

func (e *ErrMockCRUD) Range(ctx context.Context, skey string, ekey string, top int, empty func() Entity) ([]Entity, error) {
	if e.startKey {
		return nil, errors.New("range")
	}
	list := make([]Entity, top)
	return list, nil
}

func (e *ErrMockCRUD) RangeRaw(ctx context.Context, skey string, ekey string, top int) (map[string]string, error) {
	if e.startKey {
		return nil, errors.New("rangeraw")
	}
	list := make(map[string]string, top)
	return list, nil
}

// Activate activates the methods to throw a error
func (e *ErrMockCRUD) Activate(methods ...string) {
	e.Clear()

	for _, method := range methods {
		switch method {
		case "Put":
			e.put = true
		case "PutRaw":
			e.putRaw = true
		case "Remove":
			e.remove = true
		case "Get":
			e.get = true
		case "GetRaw":
			e.get = true
		case "Exists":
			e.exists = true
		case "StartKey":
			e.startKey = true
		case "Range":
			e.rang = true
		case "RangeRaw":
			e.rangRaw = true
		case "Close":
			e.close = true
		default:
			panic("method not found")
		}
	}
}

// Clear deactivates all methods
func (e *ErrMockCRUD) Clear() {
	e.put = false
	e.putRaw = false
	e.remove = false
	e.get = false
	e.getRaw = false
	e.exists = false
	e.startKey = false
	e.rang = false
	e.rangRaw = false
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
	create        bool
	put           bool
	createWithRel bool
	putWithRel    bool
	update        bool
	connectTo     bool
	listDLR       bool
	store         CRUD
}

func NewErrMockCRUDOper() *ErrMockCRUDOper {
	return &ErrMockCRUDOper{
		store: &ErrMockCRUD{},
	}
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
		case "Update":
			e.update = true
		case "CreateWithRel":
			e.createWithRel = true
		case "PutWithRel":
			e.putWithRel = true
		case "LinkTo":
			e.connectTo = true
		case "ListDLR":
			e.listDLR = true
		default:
			panic("method not found")
		}
	}
}

// Clear deactivates all methods
func (e *ErrMockCRUDOper) Clear() {
	e.create = false
	e.put = false
	e.createWithRel = false
	e.putWithRel = false
	e.update = false
	e.connectTo = false
	e.listDLR = false
}

func (e *ErrMockCRUDOper) Store() CRUD {
	return e.store
}

func (e *ErrMockCRUDOper) Create(loc string, storeTimeout StoreWithTimeout, entity Entity) (bool, error) {
	if e.create {
		return false, errors.New("create")
	}
	return true, nil
}

func (e *ErrMockCRUDOper) CreateWithRel(loc string, storeTimeout StoreWithTimeout, entity EntityRelation) (bool, bool, error) {
	if e.createWithRel {
		return false, false, errors.New("createWithRel")
	}
	return true, true, nil
}

func (e *ErrMockCRUDOper) Put(loc string, storeTimeout StoreWithTimeout, entity Entity) (bool, error) {
	if e.put {
		return false, errors.New("put")
	}
	return true, nil
}

func (e *ErrMockCRUDOper) PutWithRel(loc string, storeTimeout StoreWithTimeout, entity EntityRelation) (bool, bool, error) {
	if e.putWithRel {
		return false, false, errors.New("putWithRel")
	}
	return true, true, nil
}

func (e *ErrMockCRUDOper) Update(
	loc string,
	storeTimeout StoreWithTimeout,
	entity Entity,
	upd func(entity Entity)) (bool, error) {
	//
	if e.update {
		return false, errors.New("update")
	}
	return true, nil
}

func (e *ErrMockCRUDOper) LinkTo(
	loc string,
	storeTimeout StoreWithTimeout,
	txn Txn,
	child EntityRelation,
	parentID string,
	fill func(child Entity)) (bool, bool, Entity, error) {
	//
	if e.connectTo {
		return false, false, nil, errors.New("linkTo")
	}
	return true, true, nil, nil
}

func (e *ErrMockCRUDOper) ListDLR(storeTimeout StoreWithTimeout, childID string) ([]Entity, error) {
	if e.listDLR {
		return nil, errors.New("listDLR")
	}
	return nil, nil
}
