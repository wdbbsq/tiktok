package server

import (
	"context"
	"github.com/JirafaYe/favorite/internal/service"
	"github.com/JirafaYe/favorite/internal/store/local"
	"github.com/JirafaYe/favorite/pkg/token"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	hostname, _ = os.Hostname()
	str2Int64   = func(s string) int64 {
		ans, _ := strconv.ParseInt(s, 10, 64)
		return ans
	}
)

const (
	RedisLockSuffix = "_lock"
	LikeSuffix      = "_like__"
	UnlikeSuffix    = "_unlike"
	// ExpireTime && DDL must be bigger than FlushDuration
	ExpireTime = 30 * time.Second
	// min surviving time of keys
	DDL           = 7 * time.Second
	FlushDuration = 3 * time.Second
	MaxLockDuration
)

func init() {
	go cronJob()
}

func cronJob() {
	ticker := time.NewTicker(FlushDuration)
	ttlChecker := func(x, ttl time.Duration) bool {
		return x > 0 && x <= ttl
	}
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// check & get expired keys
			expiredKeys := m.redis.GetExpiredKeys(m.redis.GetKeys("*like*"), ttlChecker, DDL)
			// do delete cache
			for _, key := range expiredKeys {
				// double check
				if x := m.redis.CheckExpire(key); ttlChecker(x, DDL) {
					lockKey := key + RedisLockSuffix
					if m.redis.GetLock(lockKey, hostname, MaxLockDuration) {
						users := m.redis.GetSetMember(key)
						m.redis.DelKey(key)
						log.Printf("Key %s with %d values has been deleted (%d.s)\n", key, len(users), x/1000000000)
						// release lock
						m.redis.ReleaseLock(lockKey, hostname)
						go flush2Db(key, users)
					} else {
						// just fail fast and do nothing
						log.Printf("Fail to get redis lock")
					}
				}
			}
		}
	}
}

// consider mq
func flush2Db(key string, users []string) {
	if len(key) <= 7 || len(users) == 0 {
		return
	}
	videoId := str2Int64(key[:len(key)-7])
	if strings.Compare(key[len(key)-7:], LikeSuffix) == 0 {
		favorites := make([]local.UserFavorite, 0, len(users))
		for _, user := range users {
			favorites = append(favorites, local.UserFavorite{
				UserId:  str2Int64(user),
				VideoId: videoId,
			})
		}
		_ = m.db.InsertUserFavoriteInBatch(favorites)
		_ = m.db.UpdateVideoLike(videoId, len(users))
	} else {
		userIds := make([]int64, 0, len(users))
		for _, user := range users {
			userIds = append(userIds, str2Int64(user))
		}
		_ = m.db.DeleteUserFavoriteInBatch(videoId, userIds)
		_ = m.db.UpdateVideoLike(videoId, -len(users))
	}
}

type FavoriteServer struct {
	service.UnimplementedFavoriteServer
}

func (s *FavoriteServer) FavoriteAction(ctx context.Context, request *service.FavoriteActionRequest) (*service.FavoriteActionResponse, error) {
	// verify token
	claim, err := token.ParseToken(request.Token)
	if err != nil || m.redis.IsTokenExist(request.Token) {
		log.Println(err.Error())
		return ConvertActionResponse(http.StatusForbidden, "token已过期", err)
	}
	videoId := strconv.FormatInt(request.VideoId, 10)
	userId := strconv.FormatInt(claim.Id, 10)
	likeListKey := videoId + LikeSuffix
	if request.ActionType == 1 {
		// like
		if m.redis.IsTokenExist(likeListKey) {
			m.redis.AddSetWithExpire(ExpireTime, likeListKey, userId)
		} else {
			lockKey := videoId + RedisLockSuffix
			if m.redis.GetLock(lockKey, hostname, MaxLockDuration) {
				defer m.redis.ReleaseLock(lockKey, hostname)
				// get data from mysql
				if userIds, err := m.db.SelectLikesByVideo(request.VideoId); err != nil {
					log.Println(err.Error())
					return ConvertActionResponse(http.StatusInternalServerError, err.Error(), err)
				} else {
					v := make([]string, len(userIds)+1)
					for i, id := range userIds {
						v[i] = strconv.FormatInt(id, 10)
					}
					v[len(userIds)] = userId
					// flush to redis
					m.redis.AddSetWithExpire(ExpireTime, likeListKey, v...)
				}
			} else {
				// if fail to get redis lock, write mysql directly
				if err = m.db.InsertUserFavorite(claim.Id, request.VideoId); err != nil {
					log.Println(err.Error())
					return ConvertActionResponse(http.StatusInternalServerError, err.Error(), err)
				}
			}
		}
	} else {
		// unlike
		if m.redis.IsTokenExist(likeListKey) {
			m.redis.RemoveSetValue(likeListKey, userId)
		} else {
			m.redis.AddSetWithExpire(ExpireTime, videoId+UnlikeSuffix, userId)
		}
	}
	return ConvertActionResponse(0, "success", nil)
}

func (s *FavoriteServer) GetFavoriteList(ctx context.Context, request *service.FavoriteListRequest) (*service.FavoriteListResponse, error) {
	// verify token
	claim, err := token.ParseToken(request.Token)
	if err != nil || m.redis.IsTokenExist(request.Token) {
		log.Println(err.Error())
		return ConvertListResponse(http.StatusForbidden, "token已过期", nil, err)
	}
	//userId := strconv.FormatInt(claim.Id, 10)
	// what about data in redis

	// todo 分页？
	favorites, err := m.db.SelectLikesByUser(claim.Id)
	if err != nil {
		return ConvertListResponse(http.StatusInternalServerError, err.Error(), nil, err)
	}
	if len(favorites) == 0 {
		return ConvertListResponse(0, "success", nil, nil)
	}
	localVideoList, err := m.db.SelectVideos(favorites)
	if err != nil {
		return ConvertListResponse(http.StatusInternalServerError, err.Error(), nil, err)
	}
	// make map
	var authorIds []int64
	authorMap := make(map[int64]*service.UserFavor)
	for _, v := range localVideoList {
		if _, ok := authorMap[v.UserId]; !ok {
			authorIds = append(authorIds, v.UserId)
			authorMap[v.UserId] = &service.UserFavor{}
		}
	}
	authors, err := m.db.SelectUsers(authorIds)
	if err != nil {
		return ConvertListResponse(http.StatusInternalServerError, err.Error(), nil, err)
	}
	for _, a := range authors {
		authorMap[a.Id] = &a
	}
	var videoList []*service.VideoFavor
	for _, v := range localVideoList {
		videoList = append(videoList, &service.VideoFavor{
			Id:            int64(v.ID),
			Author:        authorMap[v.UserId],
			PlayUrl:       v.PlayURL,
			CoverUrl:      v.CoverURL,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Title:         v.Title,
		})
	}
	return ConvertListResponse(0, "success", videoList, nil)
}

func ConvertActionResponse(status int32, msg string, err error) (*service.FavoriteActionResponse, error) {
	return &service.FavoriteActionResponse{
		StatusCode: status,
		StatusMsg:  &msg,
	}, err
}

func ConvertListResponse(status int32, msg string, videoList []*service.VideoFavor, err error) (*service.FavoriteListResponse, error) {
	return &service.FavoriteListResponse{
		StatusCode: status,
		StatusMsg:  &msg,
		VideoList:  videoList,
	}, err
}
