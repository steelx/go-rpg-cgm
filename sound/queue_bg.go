package sound

import (
	"github.com/faiface/beep"
)

type QueueBackground struct {
	streamers []beep.StreamSeeker
}

func (q *QueueBackground) Add(streamers ...beep.StreamSeeker) {
	q.streamers = append(q.streamers, streamers...)
}

func (q *QueueBackground) Stream(samples [][2]float64) (n int, ok bool) {
	if len(q.streamers) == 0 {
		for i := range samples[0:] {
			samples[i][0] = 0
			samples[i][1] = 0
		}
		return 0, true
	}

	// We stream from the first streamer in the queue.
	n, ok = q.streamers[0].Stream(samples[0:])
	// If it's drained, we pop it from the queue, thus continuing with
	// the next streamer.
	if !ok {
		q.Pop()
	}
	return
}

func (q *QueueBackground) Err() error {
	return nil
}

//Pop out first in QueueBackground
func (q *QueueBackground) Pop() {
	if len(q.streamers) == 0 {
		return
	}
	q.streamers = q.streamers[1:]
}
