package sim

import (
	"testing"
	"time"

	"dragonsnshit/packages/common"
)

func clockSequence(times ...time.Time) func() time.Time {
	idx := 0
	return func() time.Time {
		if idx >= len(times) {
			return times[len(times)-1]
		}
		t := times[idx]
		idx++
		return t
	}
}

func TestServerConnectAssignsSlots(t *testing.T) {
	now := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	server := NewServerWithClock(2, clockSequence(now, now, now))

	id1, ok := server.Connect("alpha")
	if !ok || id1 != 0 {
		t.Fatalf("expected slot 0 to be assigned, got id=%d ok=%v", id1, ok)
	}
	id2, ok := server.Connect("beta")
	if !ok || id2 != 1 {
		t.Fatalf("expected slot 1 to be assigned, got id=%d ok=%v", id2, ok)
	}
	id3, ok := server.Connect("alpha")
	if !ok || id3 != 0 {
		t.Fatalf("expected reconnect to reuse slot 0, got id=%d ok=%v", id3, ok)
	}
	if server.slots[0].Key != "alpha" || server.slots[1].Key != "beta" {
		t.Fatalf("expected slots to hold alpha/beta keys")
	}
}

func TestServerConnectRespectsCapacity(t *testing.T) {
	server := NewServerWithClock(2, time.Now)

	if _, ok := server.Connect("alpha"); !ok {
		t.Fatalf("expected first client to connect")
	}
	if _, ok := server.Connect("beta"); !ok {
		t.Fatalf("expected second client to connect")
	}
	if _, ok := server.Connect("gamma"); ok {
		t.Fatalf("expected capacity limit to reject connection")
	}
}

func TestServerUserCmdRequiresConnection(t *testing.T) {
	first := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	second := first.Add(2 * time.Second)
	server := NewServerWithClock(1, clockSequence(first, second))

	cmd := common.UserCmd{Sequence: 7, Timestamp: 22}
	if _, ok := server.ApplyUserCmd("alpha", cmd); ok {
		t.Fatalf("expected usercmd to be rejected without connection")
	}
	if _, ok := server.Connect("alpha"); !ok {
		t.Fatalf("expected client to connect")
	}
	if _, ok := server.ApplyUserCmd("alpha", cmd); !ok {
		t.Fatalf("expected usercmd to be accepted after connection")
	}
	if server.slots[0].LastCmd.Sequence != cmd.Sequence {
		t.Fatalf("expected last cmd to be stored")
	}
	if !server.slots[0].LastSeen.Equal(second) {
		t.Fatalf("expected last seen to update to %v", second)
	}
}

func TestServerSnapshotIncludesActiveSlots(t *testing.T) {
	start := time.Date(2024, 3, 10, 10, 0, 0, 0, time.UTC)
	server := NewServerWithClock(3, clockSequence(start, start, start, start))

	if _, ok := server.Connect("alpha"); !ok {
		t.Fatalf("expected alpha to connect")
	}
	if _, ok := server.Connect("beta"); !ok {
		t.Fatalf("expected beta to connect")
	}
	cmd := common.UserCmd{Sequence: 12}
	if _, ok := server.ApplyUserCmd("alpha", cmd); !ok {
		t.Fatalf("expected usercmd to apply")
	}

	snap := server.Snapshot()
	if len(snap.Players) != 2 {
		t.Fatalf("expected 2 players in snapshot, got %d", len(snap.Players))
	}
	if snap.Players[0].ID != 0 || snap.Players[1].ID != 1 {
		t.Fatalf("expected snapshot players to be ordered by id")
	}
	if snap.Players[0].Health != defaultHealth || snap.Players[0].State != common.StateAlive {
		t.Fatalf("expected player 0 to be alive with default health")
	}
	if snap.Players[0].LastCmd.Sequence != cmd.Sequence {
		t.Fatalf("expected player 0 last cmd to be synced")
	}
}

func TestServerDamageTransitionsToDead(t *testing.T) {
	server := NewServerWithClock(1, time.Now)

	id, ok := server.Connect("alpha")
	if !ok {
		t.Fatalf("expected connect to succeed")
	}
	if !server.ApplyDamage(id, 25) {
		t.Fatalf("expected damage to apply")
	}
	if server.slots[id].Health != 75 {
		t.Fatalf("expected health to drop to 75, got %d", server.slots[id].Health)
	}
	if !server.ApplyDamage(id, 80) {
		t.Fatalf("expected fatal damage to apply")
	}
	if server.slots[id].Health != 0 || server.slots[id].State != common.StateDead {
		t.Fatalf("expected player to be dead with 0 health")
	}
}
