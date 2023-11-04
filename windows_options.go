package ttask

import (
	"time"
)

type WindowingTime string

const (
	EventTime      WindowingTime = "event-time"
	InjestionTime  WindowingTime = "injestion-time"
	ProcessingTime WindowingTime = "processing-time"
)

// Defaults:
//   - Size: 1 (min: 1)
//   - MaxIncativity: 0 (no inactivity check)
type CWOptions[T any] struct {
	Size          int
	MaxInactivity time.Duration
}

func parseCWOptions[T any](o *CWOptions[T]) {
	if o.Size == 0 {
		o.Size = 1
	}
}

// Defaults:
//   - Size: 1 second
//   - Hop: 2 seconds
//   - Watermark: 0
type HWOptions[T any] struct {
	Size          time.Duration
	Hop           time.Duration
	Watermark     time.Duration
	WindowingTime WindowingTime
}

func parseHWOptions[T any](o *HWOptions[T]) {
	if o.Size == 0 {
		o.Size = 1 * time.Second
	}

	if o.Hop == 0 {
		o.Hop = 2 * time.Second
	}

	if o.WindowingTime == "" {
		o.WindowingTime = ProcessingTime
	}
}

// Defaults:
//   - Storage: memory (no persistence)
//   - Id: random uuid
//   - MaxInactivity: 1 second
//   - MaxSize: MaxInactivity * 2
type SWOptions[T any] struct {
	MaxSize       time.Duration
	MaxInactivity time.Duration
	Watermark     time.Duration
	WindowingTime WindowingTime
}

func parseSWOptions[T any](o *SWOptions[T]) {
	if o.MaxInactivity == 0 {
		o.MaxInactivity = 1 * time.Second
	}

	if o.MaxSize == 0 {
		o.MaxSize = 2 * o.MaxInactivity
	}

	if o.WindowingTime == "" {
		o.WindowingTime = ProcessingTime
	}
}

// Defaults:
//   - Storage: memory (no persistence)
//   - Id: random uuid
//   - Size: 1 second
//   - Hop: 2 seconds
type TWOptions[T any] struct {
	Size          time.Duration
	Watermark     time.Duration
	WindowingTime WindowingTime
}

func parseTWOptions[T any](o *TWOptions[T]) {
	if o.Size == 0 {
		o.Size = 1 * time.Second
	}

	if o.WindowingTime == "" {
		o.WindowingTime = ProcessingTime
	}
}
