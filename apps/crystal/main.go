package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	W = 44
	H = 44

	NUM_BOIDS = 140
	DECAY     = 0.96
)

type Vec2 struct {
	X, Y float64
}

func (v Vec2) Add(o Vec2) Vec2 { return Vec2{v.X + o.X, v.Y + o.Y} }
func (v Vec2) Sub(o Vec2) Vec2 { return Vec2{v.X - o.X, v.Y - o.Y} }
func (v Vec2) Mul(f float64) Vec2 {
	return Vec2{v.X * f, v.Y * f}
}
func (v Vec2) Len() float64 { return math.Sqrt(v.X*v.X + v.Y*v.Y) }

func (v Vec2) Normalize() Vec2 {
	l := v.Len()
	if l == 0 {
		return Vec2{}
	}
	return v.Mul(1 / l)
}

type Boid struct {
	Pos Vec2
	Vel Vec2

	Speed float64
	Aggro float64
	Eff   float64
}

var boids []Boid

var pheromone [H][W]float64
var power [H][W]float64
var city [H][W]float64

func clear() {
	fmt.Print("\x1b[2J\x1b[H")
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func dist(a, b Vec2) float64 {
	return a.Sub(b).Len()
}

func initWorld() {
	for i := 0; i < NUM_BOIDS; i++ {
		boids = append(boids, Boid{
			Pos:   Vec2{rand.Float64() * W, rand.Float64() * H},
			Vel:   Vec2{rand.Float64()*2 - 1, rand.Float64()*2 - 1},
			Speed: 0.5 + rand.Float64(),
			Aggro: rand.Float64(),
			Eff:   rand.Float64(),
		})
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			city[y][x] = rand.Float64() * 0.3
		}
	}
}

func boidForces(b *Boid) Vec2 {
	var align Vec2
	var coh Vec2
	var sep Vec2

	count := 0

	for i := range boids {
		o := &boids[i]
		d := dist(b.Pos, o.Pos)

		if d > 0 && d < 6 {
			align = align.Add(o.Vel)
			coh = coh.Add(o.Pos)

			if d < 2 {
				sep = sep.Add(b.Pos.Sub(o.Pos))
			}

			count++
		}
	}

	if count > 0 {
		align = align.Mul(1 / float64(count)).Normalize()
		coh = coh.Mul(1 / float64(count)).Sub(b.Pos).Normalize()
		sep = sep.Normalize()
	}

	return align.Add(coh).Add(sep.Mul(1.5))
}

func updateBoids() {
	for i := range boids {
		b := &boids[i]

		f := boidForces(b)

		px := int(b.Pos.X)
		py := int(b.Pos.Y)

		if px >= 0 && px < W && py >= 0 && py < H {
			ph := pheromone[py][px]
			if ph > 0.01 {
				f = f.Add(Vec2{rand.Float64()*2 - 1, rand.Float64()*2 - 1}.Mul(ph))
			}
		}

		b.Vel = b.Vel.Add(f).Normalize().Mul(b.Speed)

		b.Pos = b.Pos.Add(b.Vel)

		if b.Pos.X < 0 {
			b.Pos.X += W
		}
		if b.Pos.Y < 0 {
			b.Pos.Y += H
		}
		if b.Pos.X >= W {
			b.Pos.X -= W
		}
		if b.Pos.Y >= H {
			b.Pos.Y -= H
		}

		x := int(b.Pos.X)
		y := int(b.Pos.Y)

		if x >= 0 && x < W && y >= 0 && y < H {
			pheromone[y][x] += 0.4 * b.Eff
			power[y][x] += 0.2
			city[y][x] += 0.02
		}

		if rand.Float64() < 0.001 {
			b.Speed = clamp(b.Speed+rand.NormFloat64()*0.05, 0.2, 2)
			b.Aggro = clamp(b.Aggro+rand.NormFloat64()*0.05, 0, 1)
			b.Eff = clamp(b.Eff+rand.NormFloat64()*0.05, 0, 1)
		}
	}
}

func updateFields() {
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			pheromone[y][x] *= DECAY

			if power[y][x] > 0.1 {
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						ny := y + dy
						nx := x + dx
						if nx >= 0 && nx < W && ny >= 0 && ny < H {
							power[ny][nx] += power[y][x] * 0.15
						}
					}
				}
			}

			power[y][x] *= 0.88

			if city[y][x] > 0.3 {
				city[y][x] += 0.01
			} else {
				city[y][x] *= 0.995
			}

			city[y][x] = clamp(city[y][x], 0, 1)
		}
	}
}

func symbol(x, y int) string {
	p := pheromone[y][x]
	pw := power[y][x]
	c := city[y][x]

	switch {
	case c > 0.85:
		return "■"
	case c > 0.6:
		return "◆"
	case pw > 1.5:
		return "⚡"
	case p > 1.2:
		return "●"
	case p > 0.6:
		return "◉"
	case c > 0.2:
		return "❖"
	default:
		return "≈"
	}
}

func render() {
	grid := make([][]string, H)
	for i := range grid {
		grid[i] = make([]string, W)
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			grid[y][x] = symbol(x, y)
		}
	}

	for i := range boids {
		x := int(boids[i].Pos.X)
		y := int(boids[i].Pos.Y)
		if x >= 0 && x < W && y >= 0 && y < H {
			grid[y][x] = "✦"
		}
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			fmt.Print(grid[y][x])
		}
		fmt.Println()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	initWorld()

	for {
		clear()

		updateBoids()
		updateFields()

		render()

		time.Sleep(60 * time.Millisecond)
	}
}
