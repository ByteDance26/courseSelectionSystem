package middle

import (
	_type "courseSelectionSystem/type"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type SimpleSession struct {
	Name      string
	SessionID string
	Value     map[string]interface{}
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func (tmp *SimpleSession) CreateSession(name string) {
	tmp.SessionID = GetMd5String(time.ANSIC + strconv.Itoa(rand.Int()))
	tmp.Value = make(map[string]interface{})
	tmp.Name = name
}

func (tmp *SimpleSession) Get(name string) interface{} {
	v, ok := tmp.Value[name]
	if !ok {
		return nil
	} else {
		return v
	}
}

func (tmp *SimpleSession) Set(name string, val interface{}) {
	tmp.Value[name] = val
}

var sessionPool map[string]*SimpleSession
var userIdMapSessionID map[string]string

func InitSimpleSessionPool() {
	sessionPool = make(map[string]*SimpleSession)
	userIdMapSessionID = make(map[string]string)
}

func getSimpleSessionBySessionID(key string) *SimpleSession {
	v, ok := sessionPool[key]
	if !ok {
		return nil
	} else {
		return v
	}
}

func InsertIntoSessionMap(s *SimpleSession) {
	sessionPool[s.SessionID] = s
}

func HandleSimpleSession(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key, err := c.Cookie(name)
		var s *SimpleSession
		if (err != nil) || (getSimpleSessionBySessionID(key) == nil) {
			s = new(SimpleSession)
			s.CreateSession(name)
			InsertIntoSessionMap(s)
			s.Value["userId"] = nil
		} else {
			s = getSimpleSessionBySessionID(key)
		}
		c.Set("session", s)
		//c.Next()
	}
}

func GetSimpleSession(c *gin.Context) (*SimpleSession, _type.ErrNo) {
	if val, has := c.Get("session"); !has {
		return nil, _type.SessionError
	} else {
		return val.(*SimpleSession), _type.OK
	}
}

/*
session has not login -> "", _type.LoginRequired
session login -> userId, _type.OK
*/
func GetUserId(c *gin.Context) (userId string, err _type.ErrNo) {
	var s *SimpleSession
	if val, has := c.Get("session"); !has {
		return "", _type.SessionError
	} else {
		s = val.(*SimpleSession)
	}
	if s.Value["userId"] != nil {
		return s.Value["userId"].(string), _type.OK
	} else {
		return "", _type.LoginRequired
	}
}

/*
set userId to session and log out previous session with same userId -> _type.OK
*/
func SetUserId(c *gin.Context, userId string) (err _type.ErrNo) {
	var s *SimpleSession
	_ = s
	if val, has := c.Get("session"); !has {
		return _type.SessionError
	} else {
		s = val.(*SimpleSession)
	}
	var mutex sync.Mutex
	mutex.Lock()
	if val, has := userIdMapSessionID[userId]; has {
		getSimpleSessionBySessionID(val).Value["userId"] = nil
	}
	userIdMapSessionID[userId] = s.SessionID
	s.Value["userId"] = userId
	mutex.Unlock()
	return _type.OK
}

/*
session not login -> _type.LoginRequired
session login, set session not login -> _type.OK
*/
func DelUserId(c *gin.Context) (err _type.ErrNo) {
	var s *SimpleSession
	_ = s
	if val, has := c.Get("session"); !has {
		return _type.SessionError
	} else {
		s = val.(*SimpleSession)
	}
	if s.Value["userId"] == nil {
		return _type.LoginRequired
	} else {
		delete(userIdMapSessionID, s.Value["userId"].(string))
		s.Value["userId"] = nil
		return _type.OK
	}
}
