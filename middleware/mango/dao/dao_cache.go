package dao

import (
	"container/list"
	"fmt"
	"github.com/x554462/gin-example/middleware/mango/library/exception"
	"strings"
	"sync"
)

func (d *Dao) query(indexes ...interface{}) ModelInterface {
	key := d.buildKey(indexes...)
	if m, err := d.GetDaoSession().daoModelCache.Get(key); err == nil {
		return m
	}
	return nil
}

func (d *Dao) remove(indexes ...interface{}) {
	key := d.buildKey(indexes...)
	d.GetDaoSession().daoModelCache.Del(key)
}

func (d *Dao) save(model ModelInterface) {
	key := d.buildKey(model.GetIndexValues()...)
	d.GetDaoSession().daoModelCache.Put(key, model)
}

func (d *Dao) buildKey(indexes ...interface{}) string {
	var err error
	buildStr := strings.Builder{}
	buildStr.WriteString(d.tableName)
	for _, v := range indexes {
		switch m := v.(type) {
		case int, int64, uint, uint64, int32, int16, int8, uint32, uint16, uint8:
			if m == 0 {
				err = exception.New("the index key can not be zero", exception.ModelRuntimeError)
			}
			buildStr.WriteString(fmt.Sprintf("`%d", m))
		case string:
			if m == "" {
				err = exception.New("the index key can not be empty string", exception.ModelRuntimeError)
			}
			buildStr.WriteString("`")
			buildStr.WriteString(m)
		default:
			err = exception.New("not support index key", exception.ModelRuntimeError)
		}
		if err != nil {
			exception.Throw(err)
		}
	}
	return buildStr.String()
}

const maxLength = 200

type element struct {
	listElem *list.Element
	model    ModelInterface
}

type DaoLruCache struct {
	elements map[string]*element
	list     *list.List
	capacity int // 容量
	used     int // 使用量
	locker   sync.RWMutex
}

func newDaoLru(capacity int) *DaoLruCache {
	size := maxLength
	if maxLength > capacity {
		size = capacity
	}
	return &DaoLruCache{
		elements: make(map[string]*element, size),
		list:     list.New(),
		capacity: capacity,
		used:     0,
	}
}

func (lru *DaoLruCache) Clear() {
	lru.elements = make(map[string]*element, lru.used)
	lru.list.Init()
	lru.used = 0
}

func (lru *DaoLruCache) Get(key string) (ModelInterface, error) {
	if lru.used > 0 {
		lru.locker.RLock()
		defer lru.locker.RUnlock()
		if element, ok := lru.elements[key]; ok {
			lru.list.MoveToBack(element.listElem)
			return element.model, nil
		}
	}
	return nil, exception.ModelNotFoundError
}

func (lru *DaoLruCache) Put(key string, model ModelInterface) {

	lru.locker.Lock()
	defer lru.locker.Unlock()

	if elem, ok := lru.elements[key]; ok {
		lru.elements[key] = &element{listElem: elem.listElem, model: model}
		lru.list.MoveToBack(elem.listElem)
		return
	}
	lru.addElement(key, model)
	if lru.used > lru.capacity {
		lru.delListFrontElement()
	}
}

func (lru *DaoLruCache) Del(key string) {
	if element, ok := lru.elements[key]; ok {
		lru.list.Remove(element.listElem)
		delete(lru.elements, key)
		lru.used--
	}
}

func (lru *DaoLruCache) addElement(key string, model ModelInterface) {
	lru.used++
	listElem := lru.list.PushBack(key)
	lru.elements[key] = &element{listElem: listElem, model: model}
}

func (lru *DaoLruCache) delListFrontElement() {
	frontElem := lru.list.Front()
	key := frontElem.Value.(string)
	lru.Del(key)
}
