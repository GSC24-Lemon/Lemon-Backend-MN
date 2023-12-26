package usecase

import (
	"github.com/gorilla/websocket"
	"lemon_be/pkg/redispkg"
	"log"
	"sort"
	"sync"
	"time"
)

type User struct {
	io   sync.Mutex
	Conn *websocket.Conn

	Id   uint
	Name string
	Hub  *Hub
}

type Hub struct {
	mu       sync.RWMutex
	seq      uint
	Rds      *redispkg.Redis
	geoRedis GeoRedisRepo
	us       []*User

	register chan *User

	unregister chan *User
}

func NewHub(rds *redispkg.Redis, geoRedis GeoRedisRepo) *Hub {
	return &Hub{
		Rds:      rds,
		geoRedis: geoRedis,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case user := <-h.register:
			user.Id = h.seq
			user.Name = user.Name

			h.us = append(h.us, user)
			h.seq++

		case user := <-h.unregister:
			h.mu.Lock()
			// binary search utk cari index user di array us
			i := sort.Search(len(h.us), func(i int) bool {
				return h.us[i].Id >= user.Id
			})

			// hapus client dari array chat.us
			without := make([]*User, len(h.us)-1) // us = nil
			copy(without[:i], h.us[:i])
			copy(without[i:], h.us[i+1:])
			h.us = without
			h.mu.Unlock()
		}
	}
}

var (
	// pongWait : berapa lama server menunggu message pong dari client (30 detik)
	pongWait = 30 * time.Second
	// pingInterval : setiap 5 detik server mengirim ping message ke client.
	// pingInterval haruslah lebih kecil dari pongWait
	pingInterval = (pongWait * 5) / 10
)

// Receive membaca next message websocket dari client
// It blocks until full message received.
func (u *User) Recieve() error {
	defer func() {
		u.Hub.unregister <- u
		u.Conn.Close()
	}()
	// Set Max Size of Messages in Bytes
	u.Conn.SetReadLimit(1024)

	if err := u.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		u.Hub.unregister <- u
		u.Conn.Close()
		return err
	}

	u.Conn.SetPongHandler(u.pongHandler)

	for {
		// ReadMessage dari client websocket
		_, _, err := u.Conn.ReadMessage()

		if err != nil {
			break
		}

	}

}

// pongHandler handle message pong yang dikirim oleh client
// -> mereset durasi readDeadline (tambah 30 detik lagi) & set user online di dalam redis
// dan juga mengirim status online user ke semua kontaknya
func (u *User) pongHandler(pongMsg string) error {

	return u.Conn.SetReadDeadline(time.Now().Add(pongWait))
}