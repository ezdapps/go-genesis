package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeToGenerate(t *testing.T) {
	cases := []struct {
		firstBlockTime time.Time
		blockGenTime   time.Duration
		blocksGap      time.Duration
		nodesCount     int64
		clock          Clock
		nodePosition   int64

		result bool
		err    error
	}{
		{
			firstBlockTime: time.Unix(1, 0),

			clock: func() Clock {
				mc := &MockClock{}
				mc.On("Now").Return(time.Unix(0, 0))
				return mc
			}(),

			err: TimeError,
		},
		{
			firstBlockTime: time.Unix(1, 0),
			blockGenTime:   time.Second * 2,
			blocksGap:      time.Second * 3,
			nodesCount:     3,
			nodePosition:   2,

			clock: func() Clock {
				mc := &MockClock{}
				mc.On("Now").Return(time.Unix(16, 0))
				return mc
			}(),

			result: true,
		},
		{
			firstBlockTime: time.Unix(0, 0),
			blockGenTime:   time.Second * 2,
			blocksGap:      time.Second * 3,
			nodesCount:     3,
			nodePosition:   3,

			clock: func() Clock {
				mc := &MockClock{}
				mc.On("Now").Return(time.Unix(0, 0))
				return mc
			}(),

			result: false,
		},
	}

	for _, c := range cases {
		btc := NewBlockTimeCalculator(c.clock,
			c.firstBlockTime,
			c.blockGenTime,
			c.blocksGap,
			c.nodesCount,
		)

		execResult, execErr := btc.TimeToGenerate(c.nodePosition)
		require.Equal(t, c.err, execErr)
		assert.Equal(t, c.result, execResult)
	}
}

func TestCountBlockTime(t *testing.T) {
	cases := []struct {
		firstBlockTime time.Time
		blockGenTime   time.Duration
		blocksGap      time.Duration
		nodesCount     int64
		clock          Clock

		result blockGenerationState
		err    error
	}{
		// Current time before first block case
		{
			firstBlockTime: time.Unix(1, 0),

			clock: func() Clock {
				mc := &MockClock{}
				mc.On("Now").Return(time.Unix(0, 0))
				return mc
			}(),

			err: TimeError,
		},

		// Zero duration case
		{
			firstBlockTime: time.Unix(0, 0),
			blockGenTime:   time.Second * 0,
			blocksGap:      time.Second * 0,
			nodesCount:     5,

			clock: func() Clock {
				mc := &MockClock{}
				mc.On("Now").Return(time.Unix(0, 0))
				return mc
			}(),

			result: blockGenerationState{
				start:    time.Unix(0, 0),
				duration: time.Second * 0,

				nodePosition: 0,
			},
		},

		// Duration testing case
		{
			firstBlockTime: time.Unix(0, 0),
			blockGenTime:   time.Second * 1,
			blocksGap:      time.Second * 0,
			nodesCount:     5,

			clock: func() Clock {
				mc := &MockClock{}
				mc.On("Now").Return(time.Unix(0, 0))
				return mc
			}(),

			result: blockGenerationState{
				start:    time.Unix(0, 0),
				duration: time.Second * 1,

				nodePosition: 0,
			},
		},

		// Duration testing case
		{
			firstBlockTime: time.Unix(0, 0),
			blockGenTime:   time.Second * 0,
			blocksGap:      time.Second * 1,
			nodesCount:     5,

			clock: func() Clock {
				mc := &MockClock{}
				mc.On("Now").Return(time.Unix(0, 0))
				return mc
			}(),

			result: blockGenerationState{
				start:    time.Unix(0, 0),
				duration: time.Second * 1,

				nodePosition: 0,
			},
		},

		// Duration testing case
		{
			firstBlockTime: time.Unix(0, 0),
			blockGenTime:   time.Second * 4,
			blocksGap:      time.Second * 6,
			nodesCount:     5,

			clock: func() Clock {
				mc := &MockClock{}
				mc.On("Now").Return(time.Unix(0, 0))
				return mc
			}(),

			result: blockGenerationState{
				start:    time.Unix(0, 0),
				duration: time.Second * 10,

				nodePosition: 0,
			},
		},

		// Block lowest time boundary case
		{
			firstBlockTime: time.Unix(0, 0),
			blockGenTime:   time.Second * 1,
			blocksGap:      time.Second * 1,
			nodesCount:     10,

			clock: func() Clock {
				mc := &MockClock{}
				mc.On("Now").Return(time.Unix(0, 0))
				return mc
			}(),

			result: blockGenerationState{
				start:    time.Unix(0, 0),
				duration: time.Second * 2,

				nodePosition: 0,
			},
		},

		// Block highest time boundary case
		{
			firstBlockTime: time.Unix(0, 0),
			blockGenTime:   time.Second * 2,
			blocksGap:      time.Second * 3,
			nodesCount:     10,

			clock: func() Clock {
				mc := &MockClock{}
				mc.On("Now").Return(time.Unix(5, 999999999))
				return mc
			}(),

			result: blockGenerationState{
				start:    time.Unix(0, 0),
				duration: time.Second * 5,

				nodePosition: 0,
			},
		},

		// Last nodePosition case
		{
			firstBlockTime: time.Unix(0, 0),
			blockGenTime:   time.Second * 0,
			blocksGap:      time.Second * 1,
			nodesCount:     3,

			clock: func() Clock {
				mc := &MockClock{}
				mc.On("Now").Return(time.Unix(6, 0))
				return mc
			}(),

			result: blockGenerationState{
				start:    time.Unix(6, 0),
				duration: time.Second * 1,

				nodePosition: 0,
			},
		},

		// One node case
		{
			firstBlockTime: time.Unix(0, 0),
			blockGenTime:   time.Second * 2,
			blocksGap:      time.Second * 2,
			nodesCount:     1,

			clock: func() Clock {
				mc := &MockClock{}
				mc.On("Now").Return(time.Unix(6, 0))
				return mc
			}(),

			result: blockGenerationState{
				start:    time.Unix(5, 0),
				duration: time.Second * 4,

				nodePosition: 0,
			},
		},

		// Custom firstBlockTime case
		{
			firstBlockTime: time.Unix(1, 0),
			blockGenTime:   time.Second * 2,
			blocksGap:      time.Second * 3,
			nodesCount:     3,

			clock: func() Clock {
				mc := &MockClock{}
				mc.On("Now").Return(time.Unix(13, 0))
				return mc
			}(),

			result: blockGenerationState{
				start:    time.Unix(13, 0),
				duration: time.Second * 5,

				nodePosition: 2,
			},
		},

		// Current time is in middle of interval case
		{
			firstBlockTime: time.Unix(1, 0),
			blockGenTime:   time.Second * 2,
			blocksGap:      time.Second * 3,
			nodesCount:     3,

			clock: func() Clock {
				mc := &MockClock{}
				mc.On("Now").Return(time.Unix(16, 0))
				return mc
			}(),

			result: blockGenerationState{
				start:    time.Unix(13, 0),
				duration: time.Second * 5,

				nodePosition: 2,
			},
		},

		// Real life case
		{
			firstBlockTime: time.Unix(1519240000, 0),
			blockGenTime:   time.Second * 4,
			blocksGap:      time.Second * 5,
			nodesCount:     101,

			clock: func() Clock {
				mc := &MockClock{}
				mc.On("Now").Return(time.Unix(1519241010, 1234))
				return mc
			}(),

			result: blockGenerationState{
				start:    time.Unix(1519241010, 0),
				duration: time.Second * 9,

				nodePosition: 0,
			},
		},
	}

	for _, c := range cases {
		btc := NewBlockTimeCalculator(c.clock,
			c.firstBlockTime,
			c.blockGenTime,
			c.blocksGap,
			c.nodesCount,
		)

		execResult, execErr := btc.countBlockTime(false, time.Time{})
		require.Equal(t, c.err, execErr)
		assert.Equal(t, c.result, execResult)
	}
}
