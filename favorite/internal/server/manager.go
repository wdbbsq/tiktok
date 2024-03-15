package server

import (
	"github.com/JirafaYe/favorite/internal/store/cache"
	"github.com/JirafaYe/favorite/internal/store/local"
	"log"
)

var m *Manager

type Manager struct {
	db    *local.Manager
	redis *cache.Manager
}

func init() {
	var err error
	db, err := local.New()
	if err != nil {
		log.Fatal(err)
	}
	redis, err := cache.New()
	if err != nil {
		log.Fatal(err)
	}

	m = &Manager{
		db:    db,
		redis: redis,
	}
}
