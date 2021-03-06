package player

import (
	"github.com/DisgoOrg/disgolink/lavalink"
	"sync"
)

type Manager struct {
	lavalink.PlayerEventAdapter
	Player  lavalink.Player
	Queue   []lavalink.AudioTrack
	QueueMu sync.Mutex
	playing bool
	paused  bool
}

// AddQueue adds 1 or more lava link tracks to manager queue
func (m *Manager) AddQueue(tracks ...lavalink.AudioTrack) {
	m.QueueMu.Lock()
	defer m.QueueMu.Unlock()
	m.Queue = append(m.Queue, tracks...)
}

// PopQueue remove first lava link track from manager queue
func (m *Manager) PopQueue() lavalink.AudioTrack {
	m.QueueMu.Lock()
	defer m.QueueMu.Unlock()
	if len(m.Queue) == 0 {
		return nil
	}
	var track lavalink.AudioTrack
	track, m.Queue = m.Queue[0], m.Queue[1:]
	return track
}

// EmptyQueue remove all elements from manager queue
func (m *Manager) EmptyQueue() {
	m.QueueMu.Lock()
	defer m.QueueMu.Unlock()
	m.Queue = nil
}

// OnTrackEnd event when track ends
func (m *Manager) OnTrackEnd(player lavalink.Player, _ lavalink.AudioTrack, endReason lavalink.AudioTrackEndReason) {
	if !endReason.MayStartNext() {
		m.playing = false
		return
	}

	if track := m.PopQueue(); track != nil {
		_ = player.Play(track)
		return
	}
	m.playing = false
}
