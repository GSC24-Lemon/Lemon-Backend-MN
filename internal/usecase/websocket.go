package usecase

import (
	"context"
	"errors"
	websocket "github.com/gorilla/websocket"
	"net/http"
)

var (
	WebsocketConnectionError   = errors.New("websocketc connection error")
	WebsocketUnauthorizedError = errors.New("websocketc unauthorized error")
)

// WebsocketUseCase bussines logic websocketc
type WebsocketUseCase struct {
	userPg       UserRepo
	geoRedisRepo GeoRedisRepo
}

// NewWebsocket Create new websocketUseCase
func NewWebsocket(hub *Hub, userPg UserRepo, geoRedisRepo GeoRedisRepo) *WebsocketUseCase {
	return &WebsocketUseCase{userPg, geoRedisRepo}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (uc *WebsocketUseCase) WebsocketHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) error {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return WebsocketConnectionError
	}

	_ = uc.
}