package entity

type Paddle struct {
	Position `json:"position"`
	Score    int     `json:"score"`
	Speed    float32 `json:"speed"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`
}

type KeysPressed struct {
	Up   bool `json:"up"`
	Down bool `json:"down"`
	W    bool `json:"W"`
	S    bool `json:"S"`
}

func (p *Paddle) Update(keysPressed *KeysPressed) {
	var screenSize ScreenSize = ScreenSize{
		Height: 100,
		Width:  100,
	}
	if keysPressed.Up || keysPressed.W {
		p.Y -= p.Speed
	} else if keysPressed.Down || keysPressed.S {
		p.Y += p.Speed
	}
	if p.Y-float32(p.Height/2) < 0 {
		p.Y = float32(1 + p.Height/2)
	} else if p.Y+float32(p.Height/2) > screenSize.Height {
		p.Y = float32(int(screenSize.Height) - p.Height/2 - 1)
	}
}
