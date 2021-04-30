package srvg

import "time"

type Option func(group *serverGroup)

func Wait(wait time.Duration) Option {
	return func(group *serverGroup) {
		group.wait = wait
	}
}
