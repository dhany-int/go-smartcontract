package redenvelope

import (
	"testing"
)

func TestGenerateRoomIdHash(t *testing.T) {
	roomId := "test-room-123"
	hash := GenerateRoomIdHash(roomId)

	if hash == [32]byte{} {
		t.Fatal("Generated hash is empty")
	}

	// Generate again, should be same
	hash2 := GenerateRoomIdHash(roomId)
	if hash != hash2 {
		t.Error("Same room ID should generate same hash")
	}

	// Different room ID should generate different hash
	hash3 := GenerateRoomIdHash("different-room")
	if hash == hash3 {
		t.Error("Different room ID should generate different hash")
	}

	t.Logf("Room ID '%s' hash: %x", roomId, hash)
}
