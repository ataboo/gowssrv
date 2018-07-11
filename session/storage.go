package session

import (
	"container/list"
	"sync"
	"time"
)

var storage *StorageProvider

type SessionStore struct {
	id        string
	lastLogin time.Time
	value     map[interface{}]interface{}
}

func (store *SessionStore) Set(key interface{}, value interface{}) error {
	store.value[key] = value
	return nil
}

func (store *SessionStore) Get(key interface{}) interface{} {
	if val, ok := store.value[key]; ok {
		return val
	}

	return nil
}

func (store *SessionStore) Delete(key interface{}) error {
	delete(store.value, key)
	return nil
}

func (store *SessionStore) Id() string {
	return store.id
}

type StorageProvider struct {
	lock     sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

func NewStorageProvider() *StorageProvider {
	return &StorageProvider{
		sync.Mutex{},
		make(map[string]*list.Element, 0),
		list.New(),
	}
}

func (sp *StorageProvider) Create(sessionId string) (*SessionStore, error) {
	sp.lock.Lock()
	defer sp.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	s := &SessionStore{id: sessionId, lastLogin: time.Now(), value: v}
	element := sp.list.PushBack(s)
	sp.sessions[sessionId] = element
	return s, nil
}

func (sp *StorageProvider) Get(sessionId string) (*SessionStore, error) {
	if element, ok := sp.sessions[sessionId]; ok {
		return element.Value.(*SessionStore), nil
	} else {
		sess, err := sp.Create(sessionId)
		return sess, err
	}
}

func (sp *StorageProvider) Destroy(sessionId string) error {
	if element, ok := sp.sessions[sessionId]; ok {
		delete(sp.sessions, sessionId)
		sp.list.Remove(element)
		return nil
	}
	return nil
}

func (sp *StorageProvider) SessionGC(maxLifeTime int64) {
	sp.lock.Lock()
	defer sp.lock.Unlock()

	for {
		element := sp.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).lastLogin.Unix() + maxLifeTime) < time.Now().Unix() {
			sp.list.Remove(element)
			delete(sp.sessions, element.Value.(*SessionStore).id)
		} else {
			break
		}
	}
}

func (sp *StorageProvider) SessionUpdate(sessionId string) {
	sp.lock.Lock()
	defer sp.lock.Unlock()

	if element, ok := sp.sessions[sessionId]; ok {
		element.Value.(*SessionStore).lastLogin = time.Now()
		sp.list.MoveToFront(element)
	}
}
