package apicache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

//Cache represents the OpenWeather reponse cache
type Cache interface {
	SetValue(id string, v []byte)
	GetValue(id string) []byte
}

type apiCache struct {
	*cache.Cache
}

//New returns new cache object
func New(d int) Cache {
	c := cache.New(time.Duration(d)*time.Minute, time.Duration(d+1)*time.Minute)
	return &apiCache{c}
}

func (ch *apiCache) SetValue(id string, v []byte) {
	ch.Set(id, v, cache.DefaultExpiration)
}

func (ch *apiCache) GetValue(id string) []byte {
	v, ok := ch.Get(id)
	if ok {
		return v.([]byte)
	}
	return nil
}
