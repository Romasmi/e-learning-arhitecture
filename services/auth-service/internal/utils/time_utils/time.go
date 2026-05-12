package time_utils

import "time"

func MinutesToDuration(minutes uint) time.Duration {
	return time.Duration(minutes) * time.Minute
}
