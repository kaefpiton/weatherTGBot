package ttl

import (
	"context"
	"sync"
	"time"
	"weatherTGBot/internal/usecase/repository"
)

type TTLRepository struct {
	mu         sync.RWMutex
	ctx        context.Context
	cancelFunc context.CancelFunc
	expireTime time.Duration
	dirty      map[any]repository.TTLItem
}

func NewTTLRepository(ExpireTimeMinute int, checkExpireEverySecond int) (repository.TLLRepository, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	repo := &TTLRepository{
		ctx:        ctx,
		cancelFunc: cancel,
		expireTime: time.Minute * time.Duration(ExpireTimeMinute),
		dirty:      map[any]repository.TTLItem{}}
	closeFunc := repo.Close

	go func(repo *TTLRepository, checkExpireEverySecond int) {
		ticker := time.NewTicker(time.Second * time.Duration(checkExpireEverySecond))
		for {
			select {
			case <-repo.ctx.Done():
				break
			case <-ticker.C:
				repo.checkExpireTimeItems()
			}
		}
	}(repo, checkExpireEverySecond)
	return repo, closeFunc
}

func (tr *TTLRepository) checkExpireTimeItems() {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	for key, value := range tr.dirty {
		if value.ExpireTime < time.Now().Unix() {
			delete(tr.dirty, key)
		}
	}
}

func (tr *TTLRepository) Close() {
	tr.cancelFunc()
}

func (tr *TTLRepository) Get(key any) (repository.TTLItem, bool) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	item, ok := tr.dirty[key]
	return item, ok
}

func (tr *TTLRepository) Put(key any, item repository.TTLItem) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	item.ExpireTime = tr.getExpireTimeForItem()
	tr.dirty[key] = item
}

func (tr *TTLRepository) getExpireTimeForItem() int64 {
	return time.Now().Add(tr.expireTime).Unix()
}
