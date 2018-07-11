package session

import (
	"net/http"
	"net/url"
	"sync"
	"time"

	"bitbucket.org/ataboo/lasecapgo/atautils"
)

var Sessions *SessionManager

func init() {
	Sessions, _ = NewSessionManager("gowssrv_session", 3600)

	go Sessions.GC()
}

type SessionError struct {
	message string
}

func (err SessionError) Error() string {
	return err.message
}

type SessionManager struct {
	cookieName  string
	lock        sync.Mutex
	storage     *StorageProvider
	maxLifeTime int64
}

func FlashMessage(w http.ResponseWriter, r *http.Request, msg string) {
	Sessions.SessionStart(w, r).Set("flash_message", msg)
}

func NewSessionManager(cookieName string, maxLifeTime int64) (*SessionManager, error) {
	return &SessionManager{
		cookieName,
		sync.Mutex{},
		NewStorageProvider(),
		maxLifeTime,
	}, nil
}

func (manager *SessionManager) SessionStart(w http.ResponseWriter, r *http.Request) *SessionStore {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		session := manager.createSession(w)
		return session
	}

	sid, _ := url.QueryUnescape(cookie.Value)
	session, _ := manager.storage.Get(sid)
	return session
}

func (manager *SessionManager) createSession(w http.ResponseWriter) *SessionStore {
	id := atautils.UniqueID()
	session, _ := manager.storage.Create(id)
	cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(id), Path: "/", HttpOnly: true, MaxAge: int(manager.maxLifeTime)}
	http.SetCookie(w, &cookie)

	return session
}

func (manager *SessionManager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}
	manager.storage.Destroy(cookie.Value)
	expiration := time.Now()
	clearCookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
	http.SetCookie(w, &clearCookie)
}

func (manager *SessionManager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.storage.SessionGC(manager.maxLifeTime)
	time.AfterFunc(time.Duration(manager.maxLifeTime), func() { manager.GC() })
}
