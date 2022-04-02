package main

import (
	"net/http"
	"os"
	"time"

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
	initBallVelocity float32       = 0.03
	leftPaddle       entity.Paddle = entity.Paddle{
		Position: entity.Position{
			X: 1,
			Y: 50,
		},
		Score:  0,
		Speed:  2.5,
		Width:  1,
		Height: 12,
	}
	rightPaddle entity.Paddle = entity.Paddle{
		Position: entity.Position{
			X: 99,
			Y: 50,
		},
		Score:  0,
		Speed:  2.5,
		Width:  1,
		Height: 12,
	}
	ball entity.Ball = entity.Ball{
		Position: entity.Position{
			X: 50,
			Y: 50,
		},
		Radius:    1.25,
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

func reset(b *entity.Ball, lPaddle *entity.Paddle, rPaddle *entity.Paddle) {
	b.X = 50
	b.Y = 50
	lPaddle.Y = 50
	rPaddle.Y = 50
}

func wsKeys(w http.ResponseWriter, r *http.Request) {
	ws, err = wsUpgrader.Upgrade(w, r, nil)
	isOpen := true
	go func(conn *websocket.Conn) {
		for isOpen {
			err = conn.ReadJSON(&keysPressed)
			if err != nil {
				isOpen = false
				reset(&ball, &leftPaddle, &rightPaddle)
				leftPaddle.Score = 0
				rightPaddle.Score = 0
				message.Update(&leftPaddle, &rightPaddle, &ball)
				break
			}
			if keysPressed.S || keysPressed.W {
				leftPaddle.Update(&keysPressed)
				ball.Update(&leftPaddle, &rightPaddle)
				message.Update(&leftPaddle, &rightPaddle, &ball)
			} else if keysPressed.Down || keysPressed.Up {
				rightPaddle.Update(&keysPressed)
				ball.Update(&leftPaddle, &rightPaddle)
				message.Update(&leftPaddle, &rightPaddle, &ball)
			}
		}
	}(ws)

	go func(conn *websocket.Conn) {
		for isOpen {
			conn.WriteJSON(message)
			ball.Update(&leftPaddle, &rightPaddle)
			if ball.X < 0 {
				rightPaddle.Score++
				reset(&ball, &leftPaddle, &rightPaddle)
			} else if ball.X > 100 {
				leftPaddle.Score++
				reset(&ball, &leftPaddle, &rightPaddle)
			}
			message.Update(&leftPaddle, &rightPaddle, &ball)
			time.Sleep(1 * time.Millisecond)
		}
	}(ws)
}

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://go-pong-client.herokuapp.com/"}
	config.AllowAllOrigins = true

	router.Use(cors.New(config))
	gin.SetMode(gin.ReleaseMode)
	router.GET("/ws", func(ctx *gin.Context) {
		wsKeys(ctx.Writer, ctx.Request)
	})

	router.Run(":" + os.Getenv("PORT"))
}
