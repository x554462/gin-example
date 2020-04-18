package mango

import (
	"fmt"
	"github.com/x554462/gin-example/middleware/mango/library/cache"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type session struct {
	sid         string
	request     *http.Request
	cacheData   map[string]interface{}
	redisClient *cache.RedisClient
}

const RedisDb = 0
const SidLength = 32
const sessionName = "sid"
const LifeTime = 864000 * time.Second

var sidReg = regexp.MustCompile(fmt.Sprintf("[a-z0-9]{%d}", SidLength))

func newSession(request *http.Request, writer http.ResponseWriter) *session {
	var sid string
	cookie, _ := request.Cookie(sessionName)
	if cookie == nil || !sidReg.MatchString(cookie.Value) {
		sid = genSid()
	} else {
		sid = cookie.Value
	}
	http.SetCookie(writer, &http.Cookie{Name: sessionName, Value: sid, Path: "/", HttpOnly: true, Secure: false, Expires: time.Now().Add(LifeTime)})
	return &session{
		request:     request,
		redisClient: cache.NewRedis(RedisDb),
		sid:         sid,
		cacheData:   make(map[string]interface{}),
	}
}
func genSid() string {
	strBuilder := strings.Builder{}
	strBuilder.Grow(SidLength)
	var str = "0123456789abcdefghijklmnopqrstuvwxyz"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < SidLength; i++ {
		strBuilder.WriteByte(str[r.Intn(len(str))])
	}
	return strBuilder.String()
}

func (s *session) expiry() {
	s.redisClient.Expire(s.sid, LifeTime)
}

func (s *session) GetSid() string {
	return s.sid
}

func (s *session) GetUserAgent() string {
	return s.request.Header.Get("HTTP_USER_AGENT")
}

func (s *session) Get(key string) interface{} {
	return s.redisClient.HGet(s.sid, key)
}

func (s *session) GetString(key string, load ...bool) string {
	val, ok := s.cacheData[key]
	if !ok || (len(load) > 0 && load[0]) {
		val = s.redisClient.HGetString(s.sid, key)
		s.cacheData[key] = val
	}
	return val.(string)
}

func (s *session) Set(key string, val interface{}) error {
	err := s.redisClient.HSet(s.sid, key, val)
	if err == nil {
		s.cacheData[key] = val
	}
	return err
}

func (s *session) Has(key string) bool {
	return s.redisClient.HExists(s.sid, key)
}

func (s *session) Del(key string) error {
	err := s.redisClient.HDel(s.sid, key)
	if err == nil {
		delete(s.cacheData, key)
	}
	return err
}

func (s *session) Destroy() {
	s.cacheData = make(map[string]interface{}, 0)
	s.redisClient.Del(s.sid)
}
