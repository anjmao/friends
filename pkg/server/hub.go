package server

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/anjmao/friends/pkg/types"
)

const (
	CheckUsersStateInterval = 300 * time.Millisecond
	pingWaitTime            = 1000 * time.Millisecond
)

type ping struct {
	userID int
	time   time.Time
}

type userLogin struct {
	req  *types.LoginRequest
	conn *ConnContext
}

type User struct {
	UserID int
	Online bool
	// Since we are sending users friends array from the client
	// for simplicity we store friends as int array, but in real world scenario
	// map[int]User, linkedList or bitmap would be a better data structure as
	// all friends could be already got from the server side.
	Friends      []int
	Conn         *ConnContext
	LastPingTime time.Time
}

// Hub holds all game state with online users
// and handles incoming TCP/UDP traffic.
type Hub struct {
	users map[int]*User

	login chan *userLogin
	ping  chan *ping
}

func NewHub() *Hub {
	return &Hub{
		users: make(map[int]*User),
		login: make(chan *userLogin, 10),
		ping:  make(chan *ping, 10),
	}
}

// Run stars main game loop and controls state changes
// via channels which allows to prevent use of mutexes.
// done channel could be used for both stopping the loop and
// making unit testing much easier.
func (h *Hub) Run(checkTick <-chan time.Time, done <-chan struct{}) {
	for {
		select {
		case u := <-h.login:
			if err := h.handleLogin(u); err != nil {
				logrus.Errorf("could not handle login: %v", err)
			}
		case p := <-h.ping:
			if err := h.handlePing(p); err != nil {
				logrus.Errorf("could not handle ping: %v", err)
			}
		case <-checkTick:
			h.checkUsersState()
		case <-done:
			return
		}
	}
}

// IncomingMessageHandler handles incoming message.
func (h *Hub) IncomingMessageHandler(ctx *ConnContext, msg *types.Msg) {
	switch msg.Cmd {
	case types.CmdLogin:
		req := new(types.LoginRequest)
		if err := json.Unmarshal(msg.Data, req); err != nil {
			logrus.Errorf("could not parse login message: %v", err)
			return
		}

		h.login <- &userLogin{
			req:  req,
			conn: ctx,
		}
	case types.CmdPing:
		req := new(types.PingRequest)
		if err := json.Unmarshal(msg.Data, req); err != nil {
			logrus.Errorf("could not parse ping message: %v", err)
			return
		}

		h.ping <- &ping{userID: req.UserID, time: time.Now()}
	default:
		logrus.Errorf("unknown command: %b", msg.Cmd)
	}
}

// Users returns readonly users map.
func (h *Hub) Users() map[int]*User {
	return h.users
}

func (h *Hub) handleLogin(login *userLogin) error {
	logrus.Infof("user=%d friends=%v connected", login.req.UserID, login.req.Friends)
	u := &User{
		UserID:       login.req.UserID,
		Friends:      login.req.Friends,
		Online:       true,
		Conn:         login.conn,
		LastPingTime: time.Now().Add(pingWaitTime),
	}
	h.users[u.UserID] = u
	return h.notifyFriends(u, true)
}

func (h *Hub) handlePing(p *ping) error {
	u, ok := h.users[p.userID]
	if !ok {
		return fmt.Errorf("user %d not found", p.userID)
	}
	u.LastPingTime = p.time
	return nil
}

func (h *Hub) checkUsersState() {
	// 1 Step. Loop through all users and check last ping time.
	// Mark user as offline if no ping was received after pingWaitTime interval.
	for _, u := range h.users {
		if u.Online && u.LastPingTime.Add(pingWaitTime).Before(time.Now()) {
			logrus.Infof("user=%d disconnected", u.UserID)
			u.Online = false
			if err := u.Conn.close(); err != nil {
				logrus.Errorf("could not close client Conn: %v", err)
			}
		}
	}

	// 2 Step. Loop through all offline users and notify such users
	// friends about status change.
	for _, u := range h.users {
		if u.Online {
			continue
		}
		if err := h.notifyFriends(u, false); err != nil {
			logrus.Errorf("could not notify User's %d Friends: %v", u.UserID, err)
		}

		delete(h.users, u.UserID)
	}
}

// notifyFriends notifies all user's online friends about his
// status change.
func (h *Hub) notifyFriends(u *User, online bool) error {
	status := &types.StatusChangeReply{UserID: u.UserID, Online: online}
	msg, err := types.EncodeMsg(types.CmdStatusChange, status)
	if err != nil {
		return err
	}

	var writeErr error
	for _, friendID := range u.Friends {
		if f, ok := h.users[friendID]; ok && f.Online {
			writeErr = f.Conn.write(msg)
		}
	}
	return writeErr
}
