package ktntwitchgo

import "fmt"

type TwitchApiRateLimit struct {
	Limit		int			`json:"limit"`
	Remaining	int			`json:"remaining"`
	Reset		int			`json:"reset"`
}

type TwitchApiRateLimitError struct {
	RateLimit	TwitchApiRateLimit
}

func (e *TwitchApiRateLimitError) Error() string {
	return fmt.Sprintf("twitch api is rate limited (limit %d,remaining %d,reset %d)", e.RateLimit.Limit, e.RateLimit.Remaining, e.RateLimit.Reset)
}
