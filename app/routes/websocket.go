package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	error_model "github.com/MauricioMilano/stock_app/models/error"
	"github.com/MauricioMilano/stock_app/services/rabbitmq"
	"github.com/MauricioMilano/stock_app/services/websocket"
	error_utils "github.com/MauricioMilano/stock_app/utils/error"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Println("Erro ao atualizar a conexão:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Erro ao ler a mensagem:", err)
			return
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			fmt.Println("Erro ao escrever a mensagem:", err)
			return
		}
	}
}

var RegisterWebsocketRoute = func(router *mux.Router) {
	pool := websocket.NewPool()
	go pool.Start()
	sb := router.PathPrefix("/v1").Subrouter()
	sb.HandleFunc("/test", handleConnections)
	sb.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		jwtToken := r.URL.Query().Get("jwt")
		jwtSecret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			handleWebsocketAuthenticationErr(w, err)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			handleWebsocketAuthenticationErr(w, err)
			return
		}

		serveWS(pool, w, r, claims)

	})
}

func serveWS(pool *websocket.Pool, w http.ResponseWriter, r *http.Request, claims jwt.MapClaims) {
	conn, err := websocket.Upgrade(w, r)
	error_utils.ErrorCheck(err)
	br := rabbitmq.GetRabbitMQBroker()

	client := &websocket.Client{
		Connection: conn,
		Pool:       pool,
		Name:       claims["UserName"].(string),
		UserID:     uint(claims["UserID"].(float64)),
	}

	fmt.Println("Websocket ready to accept connections")
	pool.Register <- client
	requestBody := make(chan []byte) // websocket.Message byte array channel
	go client.Read(requestBody)
	go br.ReadMessages(pool)
	go br.PublishMessage(requestBody)
}

func handleWebsocketAuthenticationErr(w http.ResponseWriter, err error) {
	log.Println("websocket error: ", err)
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := error_model.ErrorResponse{Message: err.Error(), Status: false, Code: http.StatusUnauthorized}
	data, err := json.Marshal(res)
	error_utils.ErrorCheck(err)
	w.Write(data)
}
