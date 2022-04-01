package cache

import (
	"credit-line/internal/model"
)

// cache the store cache struct
type cache struct {
	CurrentCreditStatus string
	RequestFailed       map[string]uint
}

// cacheRequest variable with the cache
var cacheRequest *cache = &cache{
	CurrentCreditStatus: "",
	RequestFailed:       make(map[string]uint),
}

// RetrieveRequestCache retrieves the current cache state
func RetrieveRequestCache() *cache {
	return cacheRequest
}

// UpdateRequestCache set a new values in the cache store
func UpdateRequestCache(cs model.CreditStatus, ip string) {
	cacheRequest.CurrentCreditStatus = string(cs)

	if model.Approved == cs {
		cacheRequest.RequestFailed[ip] = 0
		return
	}
	cacheRequest.RequestFailed[ip] = cacheRequest.RequestFailed[ip] + 1
}
