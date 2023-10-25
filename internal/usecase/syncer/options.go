package syncer

import "time"

type Options func(os *syncer)

func RequestInterval(i time.Duration) Options {
	return func(os *syncer) {
		os.i = i
	}
}
