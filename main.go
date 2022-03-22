package main

import (
	"log"
	"net/http"

	"github.com/gautamprikshit1/pong-game-backend/entity"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	wsUpgrader *websocket.Upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	keysPressed      entity.KeysPressed
	initBallVelocity float32       = 0.05
	leftPaddle       entity.Paddle = entity.Paddle{
		Position: entity.Position{
			X: 5,
			Y: 50,
		},
		Score:  0,
		Speed:  0.1,
		Width:  1,
		Height: 12,
	}
	rightPaddle entity.Paddle = entity.Paddle{
		Position: entity.Position{
			X: 95,
			Y: 50,
		},
		Score:  0,
		Speed:  0.1,
		Width:  1,
		Height: 12,
	}
	ball entity.Ball = entity.Ball{
		Position: entity.Position{
			X: 50,
			Y: 50,
		},
		Radius:    2.5,
		XVelocity: initBallVelocity,
		YVelocity: initBallVelocity,
	}
	message entity.Message = entity.Message{
		LeftPaddle:  &leftPaddle,
		RightPaddle: &rightPaddle,
		PongBall:    &ball,
	}
	ws  *websocket.Conn
	err error
)

func wsKeys(w http.ResponseWriter, r *http.Request) {
	ws, err = wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error upgrading the connection: ", err)
	}
	defer ws.Close()
	for {
		err = ws.WriteJSON(message)
		if err != nil {
			log.Fatal("Error occured writing JSON: ", err)
		}
		ball.Update(&leftPaddle, &rightPaddle)
		message.Update(&leftPaddle, &rightPaddle, &ball)
		err = ws.ReadJSON(&keysPressed)
		if err != nil {
			log.Println("Error reading JSON: ", err)
			return
		}
		if keysPressed.S || keysPressed.W {
			leftPaddle.Update(&keysPressed)
			message.Update(&leftPaddle, &rightPaddle, &ball)
			ws.WriteJSON(message)
		} else if keysPressed.Down || keysPressed.Up {
			rightPaddle.Update(&keysPressed)
			message.Update(&leftPaddle, &rightPaddle, &ball)
			ws.WriteJSON(message)
		}
	}
}

/*func wsPaddleLeft(w http.ResponseWriter, r *http.Request) {
	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error upgrading the connection: ", err)
	}
	defer ws.Close()
	for {
		err = ws.WriteJSON(leftPaddle)
		if err != nil {
			log.Fatal("Error occured writing JSON: ", err)
		}
		leftPaddle.Update(&keysPressed)
		fmt.Println("Left Paddle: ", leftPaddle)
	}
}
*/

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/ws", func(ctx *gin.Context) {
		wsKeys(ctx.Writer, ctx.Request)
	})

	router.Run(":5000")
}
