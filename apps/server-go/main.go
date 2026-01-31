package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"time"

	"dragonsnshit/packages/common"
	"dragonsnshit/server/player"
	"dragonsnshit/server/system"
)

type world struct{}

type rayResult struct {
	pos system.Vec3
}

func (r rayResult) Position() system.Vec3 { return r.pos }

func (w world) RayTrace(start, end system.Vec3) (player.RaycastResult, bool) {
	return rayResult{}, false
}

type shankPlayer struct {
	pos       system.Vec3
	eyeHeight float64
	world     world
}

func (p *shankPlayer) Position() system.Vec3 { return p.pos }
func (p *shankPlayer) EyeHeight() float64    { return p.eyeHeight }
func (p *shankPlayer) World() player.RaycastWorld {
	return p.world
}
func (p *shankPlayer) SendSound(name string, pos system.Vec3) {
	fmt.Printf("[sound] %s at %.2f %.2f %.2f\n", name, pos.X, pos.Y, pos.Z)
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":6969")
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Go backend listening on :6969")
	buf := make([]byte, 2048)
	p := &shankPlayer{pos: system.Vec3{}, eyeHeight: 1.62, world: world{}}

	for {
		conn.SetReadDeadline(time.Now().Add(250 * time.Millisecond))
		n, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			fmt.Printf("read error: %v\n", err)
			continue
		}
		if n < 1 {
			continue
		}
		const netHeaderSize = 12
		const userCmdSize = 36
		switch buf[0] {
		case common.PacketUserCmd:
			if n < netHeaderSize+1+userCmdSize {
				continue
			}
			count := int(buf[netHeaderSize])
			if count < 1 {
				continue
			}
			cmd := parseUserCmd(buf, netHeaderSize+1)
			if cmd.Buttons&common.BtnAttack != 0 {
				player.HandleShankFire(p, float64(cmd.Yaw), float64(cmd.Pitch), int(cmd.WeaponIdx))
			}
			_ = remote
		}
	}
}

func parseUserCmd(data []byte, offset int) common.UserCmd {
	off := offset
	cmd := common.UserCmd{}
	cmd.Sequence = binary.LittleEndian.Uint32(data[off:])
	off += 4
	cmd.Timestamp = binary.LittleEndian.Uint32(data[off:])
	off += 4
	cmd.Msec = binary.LittleEndian.Uint16(data[off:])
	off += 4
	cmd.Fwd = math.Float32frombits(binary.LittleEndian.Uint32(data[off:]))
	off += 4
	cmd.Str = math.Float32frombits(binary.LittleEndian.Uint32(data[off:]))
	off += 4
	cmd.Yaw = math.Float32frombits(binary.LittleEndian.Uint32(data[off:]))
	off += 4
	cmd.Pitch = math.Float32frombits(binary.LittleEndian.Uint32(data[off:]))
	off += 4
	cmd.Buttons = binary.LittleEndian.Uint32(data[off:])
	off += 4
	cmd.WeaponIdx = int32(binary.LittleEndian.Uint32(data[off:]))
	return cmd
}
