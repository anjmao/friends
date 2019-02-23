package server

import (
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/anjmao/friends/pkg/types"
)

func TestHubAcceptLogin(t *testing.T) {
	h := NewHub()
	mockTicker := make(chan time.Time)
	done := make(chan struct{})
	go h.Run(mockTicker, done)
	login := &types.LoginRequest{UserID: 1, Friends: []int{2, 3, 4}}
	b, _ := types.EncodeMsg(types.CmdLogin, login)
	msg := types.DecodeMsg(b)
	connCtx := createConnContext()

	h.IncomingMessageHandler(connCtx, msg)
	done <- struct{}{}

	if len(h.users) == 0 {
		t.Fatal("expected to add new user")
	}
	usr := h.users[1]
	if usr.UserID != login.UserID {
		t.Errorf("expected userID %d, got %d", login.UserID, usr.UserID)
	}
	if !usr.Online {
		t.Error("expected user to be online")
	}
	if usr.Conn.tcpConn != connCtx.tcpConn {
		t.Errorf("expected user to get tcp conn")
	}
	if !reflect.DeepEqual(usr.Friends, login.Friends) {
		t.Errorf("expected friends %v, got %v", login.Friends, usr.Friends)
	}
}

func TestHubAcceptPing(t *testing.T) {
	h := NewHub()
	mockTicker := make(chan time.Time)
	done := make(chan struct{})
	go h.Run(mockTicker, done)
	user := createOnlineUser(1, []int{})
	h.users[user.UserID] = user

	ping := &types.PingRequest{UserID: user.UserID}
	b, _ := types.EncodeMsg(types.CmdPing, ping)
	msg := types.DecodeMsg(b)
	connCtx := createConnContext()

	h.IncomingMessageHandler(connCtx, msg)
	done <- struct{}{}

	if user.LastPingTime.IsZero() {
		t.Error("expected to set last ping time")
	}
}

func TestHubRemoveOfflineUsers(t *testing.T) {
	h := NewHub()
	mockTicker := make(chan time.Time)
	done := make(chan struct{})
	go h.Run(mockTicker, done)
	user1 := createOnlineUser(1, []int{2, 3})
	user2 := createOnlineUser(2, []int{1})
	user3 := createOnlineUser(3, []int{1})
	h.users = map[int]*User{
		1: user1,
		2: user2,
		3: user3,
	}

	mockTicker <- time.Now()
	done <- struct{}{}

	if len(h.users) != 0 {
		t.Fatalf("expected to remove offline Users, got %v", h.users)
	}
}

// BenchmarkCheckUsersState shows that current user data structure
// there each user holds friends as int array can be improved and instead
// friends could be linked as linkedList, map or bitmap for better performance.
func BenchmarkCheckUsersState(b *testing.B) {
	h := NewHub()
	totalUsers := 20000
	offlineUsers := 10
	// Each user is a friend of each other user.
	for i := 0; i < totalUsers; i++ {
		var friends []int
		for j := 0; j < totalUsers; j++ {
			if j != i {
				friends = append(friends, j)
			}
		}
		user := createOnlineUser(i, friends)
		user.LastPingTime = time.Now().Add(10 * time.Second)
		h.users[i] = user
	}

	// Make some users offline.
	for i := 0; i < offlineUsers; i++ {
		h.users[i].LastPingTime = time.Now().Add(-10 * time.Second)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h.checkUsersState()
	}
}

func createConnContext() *ConnContext {
	return &ConnContext{tcpConn: &mockTCPConn{}}
}

func createOnlineUser(userID int, friends []int) *User {
	return &User{
		UserID:  userID,
		Online:  true,
		Friends: friends,
		Conn:    createConnContext(),
	}
}

type mockTCPConn struct {
	// Embed net.Conn so we don't need to implement all methods.
	net.Conn
}

func (mockTCPConn) Close() error {
	return nil
}

func (mockTCPConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (mockTCPConn) Write(b []byte) (int, error) {
	return 0, nil
}
