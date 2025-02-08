package rate_limit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cache "github.com/harryosmar/cache-go"
	coreError "github.com/harryosmar/go-echo-core/error"
)

type Visitor struct {
	Counter         int
	UnixTimeVisited int64
}

type RateLimiterSlidingWindowCacheStore struct {
	TimeWindow int64         // window of time until visitors' counters reset
	BurstLimit int           // Limit of how many requests are allowed per time window
	ExpiresIn  time.Duration // ExpiresIn is the duration after that a rate limiter is cleaned up
	cache      cache.CacheRepo
}

func NewRateLimiterSlidingWindowCacheStore(timeWindow int64, burst int, expiresIn time.Duration, cache cache.CacheRepo) (store *RateLimiterSlidingWindowCacheStore) {
	store = &RateLimiterSlidingWindowCacheStore{TimeWindow: timeWindow, BurstLimit: burst, ExpiresIn: expiresIn, cache: cache}
	return
}

func (r *RateLimiterSlidingWindowCacheStore) Allow(identifier string) (bool, error) {
	tempRedisKeyPatternString := "rate_limiter"
	identifierKey := fmt.Sprintf("%s%s", tempRedisKeyPatternString, identifier)
	tempContext := context.TODO()

	visitorData, found, err := r.cache.Get(tempContext, identifierKey)
	if err != nil {
		return false, coreError.ErrGeneral.WithError(err)
	}

	var visitor Visitor
	if !found {
		visitor = Visitor{
			Counter:         1,
			UnixTimeVisited: time.Now().Unix(),
		}
	} else {
		if err := json.Unmarshal(visitorData, &visitor); err != nil {
			return false, coreError.ErrGeneral.WithError(err)
		}
	}

	if !r.isAllowed(&visitor) {
		return false, nil
	}

	visitorJson, err := json.Marshal(visitor)
	if err != nil {
		return false, coreError.ErrGeneral.WithError(err)
	}
	if err := r.cache.Store(tempContext, identifierKey, visitorJson, r.ExpiresIn); err != nil {
		return false, coreError.ErrGeneral.WithError(err)
	}
	return true, nil
}

func (r *RateLimiterSlidingWindowCacheStore) isAllowed(visitor *Visitor) bool {
	unixTimeNow := time.Now().Unix()

	if unixTimeNow-visitor.UnixTimeVisited > r.TimeWindow {
		visitor.Counter = 1
		visitor.UnixTimeVisited = unixTimeNow
		return true
	}

	if visitor.Counter >= r.BurstLimit {
		return false
	}

	visitor.Counter++
	return true
}
