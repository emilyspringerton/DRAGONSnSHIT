package sim

import (
	"sort"
	"time"

	"dragonsnshit/packages/common"
)

const defaultHealth = 100

type Slot struct {
	ID       int
	Key      string
	Active   bool
	LastCmd  common.UserCmd
	LastSeen time.Time
	Health   int
	State    int
}

type PlayerSnapshot struct {
	ID       int
	Health   int
	State    int
	LastCmd  common.UserCmd
	LastSeen time.Time
}

type Snapshot struct {
	Timestamp time.Time
	Players   []PlayerSnapshot
}

type Server struct {
	slots []Slot
	clock func() time.Time
}

func NewServer(maxClients int) *Server {
	return NewServerWithClock(maxClients, time.Now)
}

func NewServerWithClock(maxClients int, clock func() time.Time) *Server {
	if maxClients <= 0 {
		maxClients = common.MaxClients
	}
	return &Server{
		slots: make([]Slot, maxClients),
		clock: clock,
	}
}

func (s *Server) Connect(key string) (int, bool) {
	if key == "" {
		return -1, false
	}
	if idx, ok := s.findSlot(key); ok {
		return idx, true
	}
	for i := range s.slots {
		if !s.slots[i].Active {
			s.slots[i] = Slot{
				ID:       i,
				Key:      key,
				Active:   true,
				LastSeen: s.clock(),
				Health:   defaultHealth,
				State:    common.StateAlive,
			}
			return i, true
		}
	}
	return -1, false
}

func (s *Server) Disconnect(key string) {
	idx, ok := s.findSlot(key)
	if !ok {
		return
	}
	s.slots[idx] = Slot{ID: idx}
}

func (s *Server) ApplyUserCmd(key string, cmd common.UserCmd) (int, bool) {
	idx, ok := s.findSlot(key)
	if !ok {
		return -1, false
	}
	slot := s.slots[idx]
	slot.LastCmd = cmd
	slot.LastSeen = s.clock()
	s.slots[idx] = slot
	return idx, true
}

func (s *Server) ApplyDamage(id int, amount int) bool {
	if id < 0 || id >= len(s.slots) {
		return false
	}
	slot := s.slots[id]
	if !slot.Active {
		return false
	}
	if amount <= 0 {
		return true
	}
	slot.Health -= amount
	if slot.Health <= 0 {
		slot.Health = 0
		slot.State = common.StateDead
	}
	s.slots[id] = slot
	return true
}

func (s *Server) Snapshot() Snapshot {
	players := make([]PlayerSnapshot, 0, len(s.slots))
	for _, slot := range s.slots {
		if !slot.Active {
			continue
		}
		players = append(players, PlayerSnapshot{
			ID:       slot.ID,
			Health:   slot.Health,
			State:    slot.State,
			LastCmd:  slot.LastCmd,
			LastSeen: slot.LastSeen,
		})
	}
	sort.Slice(players, func(i, j int) bool { return players[i].ID < players[j].ID })
	return Snapshot{Timestamp: s.clock(), Players: players}
}

func (s *Server) findSlot(key string) (int, bool) {
	for i, slot := range s.slots {
		if slot.Active && slot.Key == key {
			return i, true
		}
	}
	return -1, false
}
