package playermanager

import (
	"github.com/DisgoOrg/disgolink/lavalink"
	"sync"
)

type PlayerManager struct {
	lavalink.PlayerEventAdapter
	Player  lavalink.Player
	Queue   []lavalink.AudioTrack
	QueueMu sync.Mutex
}

func (m *PlayerManager) AddQueue(tracks ...lavalink.AudioTrack) {
	m.QueueMu.Lock()
	defer m.QueueMu.Unlock()
	m.Queue = append(m.Queue, tracks...)
}

func (m *PlayerManager) PopQueue() lavalink.AudioTrack {
	m.QueueMu.Lock()
	defer m.QueueMu.Unlock()
	if len(m.Queue) == 0 {
		return nil
	}
	var track lavalink.AudioTrack
	track, m.Queue = m.Queue[0], m.Queue[1:]
	return track
}

func (m *PlayerManager) OnTrackEnd(player lavalink.Player, _ lavalink.AudioTrack, endReason lavalink.AudioTrackEndReason) {
	if !endReason.MayStartNext() {
		return
	}

	if track := m.PopQueue(); track != nil {
		_ = player.Play(track)
	}
}
