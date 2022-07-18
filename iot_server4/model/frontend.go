package model

import (
	"errors"
	"time"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
)

//代表每个长连接客户端的缓存参数
type User struct {
	ClientSession cellnet.Session
	LastPingTime  time.Time
	TimerHandler  *time.Timer //用于延迟发送心跳响应
}

var FrontendSessionManager peer.SessionManager

//从session中获取user
func SessionToUser(clientSes cellnet.Session) *User {
	if clientSes == nil {
		return nil
	}
	if raw, ok := clientSes.(cellnet.ContextSet).GetContext("user"); ok {
		return raw.(*User)
	}
	return nil
}

//遍历所有持有session中的user
func VisitUser(callback func(*User) bool) {
	FrontendSessionManager.VisitSession(func(clientSes cellnet.Session) bool {
		if u := SessionToUser(clientSes); u != nil {
			return callback(u)
		}
		return true
	})
}

func NewUser(clientSes cellnet.Session) *User {
	return &User{
		ClientSession: clientSes, LastPingTime: time.Now(),
	}
}

func CreateUser(ses cellnet.Session) (*User, error) {
	u := SessionToUser(ses)
	if u != nil {
		return nil, errors.New("user already bind")
	}
	u = NewUser(ses)
	ses.(cellnet.ContextSet).SetContext("user", u)
	return u, nil
}
