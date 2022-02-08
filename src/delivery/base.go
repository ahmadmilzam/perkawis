package delivery

import (
	"github.com/ReneKroon/ttlcache/v2"
)

var notFound = ttlcache.ErrNotFound
var cache ttlcache.SimpleCache = ttlcache.NewCache()
