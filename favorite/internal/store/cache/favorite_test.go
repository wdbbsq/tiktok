package cache

import (
	"fmt"
	"github.com/JirafaYe/favorite/internal/store/local"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

func TestLike(t *testing.T) {
	m := &Manager{
		handler: redis.NewClient(
			&redis.Options{
				Addr:     C.Addr + ":" + C.Port,
				Password: C.Password,
			},
		),
	}
	nums := []string{"1", "2", "3", "786576567"}
	fmt.Println(m.AddSetWithExpire(60*time.Second, "xxxx", nums...))
	fmt.Println(m.CheckExpire("xxxx"))
}

func TestDelete(t *testing.T) {
	db, _ := local.New()
	fmt.Println(db.SelectLikesByVideo(1))
	fmt.Println(db.SelectLikesByVideo(2))
	//fmt.Println(db.InsertUserFavorite(123, 1))
	//fmt.Println(db.InsertUserFavorite(233, 1))
	//fmt.Println(db.InsertUserFavorite(321, 1))
	//fmt.Println(db.DeleteUserFavoriteInBatch(1, []int64{123, 233}))
}

func TestTTL(t *testing.T) {
	m := &Manager{
		handler: redis.NewClient(
			&redis.Options{
				Addr:     C.Addr + ":" + C.Port,
				Password: C.Password,
			},
		),
	}
	//fmt.Println(m.GetExpiredKeys(m.GetKeys("*")))

	keys := m.GetKeys("*like*")
	ans := make([]time.Duration, len(keys))
	for i := 0; i < len(keys); i++ {
		ans[i] = m.CheckExpire(keys[i])
	}
	fmt.Println(keys)
	fmt.Println(ans)
}
