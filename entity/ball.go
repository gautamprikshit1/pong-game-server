package entity

type Position struct {
	X float32 `json:"X"`
	Y float32 `json:"Y"`
}

type Ball struct {
	Position  `json:"position"`
	Radius    float32 `json:"radius"`
	XVelocity float32 `json:"xvelocity"`
	YVelocity float32 `json:"yvelocity"`
}

type ScreenSize struct {
	Width  float32
	Height float32
}

func (b *Ball) Update(leftPaddle *Paddle, rightPaddle *Paddle) {
	var screenSize ScreenSize = ScreenSize{
		Width:  100,
		Height: 100,
	}

	b.X += b.XVelocity
	b.Y += b.YVelocity

	if b.Y-b.Radius > screenSize.Height {
		b.YVelocity = -b.YVelocity
		b.Y = screenSize.Height - b.Radius
	} else if b.Y+b.Radius < 0 {
		b.YVelocity = -b.YVelocity
		b.Y = b.Radius
	}
	if b.X-b.Radius < leftPaddle.X+float32(leftPaddle.Width/2) &&
		b.Y > leftPaddle.Y-float32(leftPaddle.Height/2) &&
		b.Y < leftPaddle.Y+float32(leftPaddle.Height/2) {
		b.XVelocity = -b.XVelocity
		b.X = leftPaddle.X + float32(leftPaddle.Width/2) + b.Radius
	} else if b.X+b.Radius > rightPaddle.X-float32(rightPaddle.Width/2) &&
		b.Y > rightPaddle.Y-float32(rightPaddle.Height/2) &&
		b.Y < rightPaddle.Y+float32(rightPaddle.Height/2) {
		b.XVelocity = -b.XVelocity
		b.X = rightPaddle.X - float32(rightPaddle.Width/2) - b.Radius
	}
}
