package auth

import "time"

type timeStamp int64

func newTimestamp() timeStamp {
	return timeStamp(time.Now().Unix())
}

func (ts timeStamp) Add(seconds int64) timeStamp {
	return ts + timeStamp(seconds)
}

func (ts timeStamp) AddDuration(interval time.Duration) timeStamp {
	return ts + timeStamp(interval/time.Second)
}

func (ts timeStamp) GetTime() time.Time {
	return time.Unix(int64(ts), 0).Local()
}
