package entity

type Message struct {
	LeftPaddle  *Paddle `json:"leftPaddle"`
	RightPaddle *Paddle `json:"rightPaddle"`
	PongBall    *Ball   `json: "pongBall"`
}

func (*Message) Update(p1 *Paddle, p2 *Paddle, ball *Ball) Message {
	return Message{
		LeftPaddle:  p1,
		RightPaddle: p2,
		PongBall:    ball,
	}
}
