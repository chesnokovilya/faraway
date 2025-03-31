package pow

import "time"

const (
	ChallengePrefix = "CHALLENGE:"
	FailPrefix      = "FAIL:"
	ErrorPrefix     = "ERROR:"
	WowPrefix       = "WOW:"
	Delimiter       = ":"

	Difficulty = 24

	Timeout = 10 * time.Second
)

type Challenge struct {
	Value      string
	Difficulty int
}

type Solution struct {
	Nonce string
	Hash  string
}
