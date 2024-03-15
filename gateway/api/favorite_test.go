package api

import (
	"encoding/json"
	"github.com/JirafaYe/gateway/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFavoriteActionType1(t *testing.T) {
	engine := route()

	url := "http://127.0.0.1:8088/douyin/favorite/action?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MjI2ODM0NDQ3NzEwODE0NjE3NiwidXNlcm5hbWUiOiJxcSIsImV4cCI6MTY3NzU4MDI2NywiaXNzIjoidGlrdG9rLnVzZXIifQ.zU_SepCAKXJDfWR7wflvDBtb5_pjy8R734Qw33uqydg&video_id=1&action_type=1"
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", url, nil)
	engine.ServeHTTP(recorder, request)

	var res service.FavoriteActionResponse
	json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, 0, int(res.StatusCode))
	log.Println(res)
}

func TestFavoriteActionType2(t *testing.T) {
	engine := route()

	url := "http://127.0.0.1:8088/douyin/favorite/action?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MjI2ODM0NDQ3NzEwODE0NjE3NiwidXNlcm5hbWUiOiJxcSIsImV4cCI6MTY3NzU4MDI2NywiaXNzIjoidGlrdG9rLnVzZXIifQ.zU_SepCAKXJDfWR7wflvDBtb5_pjy8R734Qw33uqydg&video_id=1&action_type=2"
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", url, nil)
	engine.ServeHTTP(recorder, request)

	var res service.FavoriteActionResponse
	json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, 0, int(res.StatusCode))
	log.Println(res)
}

func TestFavoriteListType(t *testing.T) {
	engine := route()

	url := "http://127.0.0.1:8088/douyin/favorite/list?user_id=2268344477108146176&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MjI2ODM0NDQ3NzEwODE0NjE3NiwidXNlcm5hbWUiOiJxcSIsImV4cCI6MTY3NzU4MDI2NywiaXNzIjoidGlrdG9rLnVzZXIifQ.zU_SepCAKXJDfWR7wflvDBtb5_pjy8R734Qw33uqydg"
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", url, nil)
	engine.ServeHTTP(recorder, request)

	var res service.FavoriteActionResponse
	json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, 0, int(res.StatusCode))
	log.Println(res)
}

func getEngine() *gin.Engine {
	app := New()
	err := app.loadRoute()
	if err != nil {
		panic(err)
	}
	return app.handler
}

func BenchmarkFavoriteAction(b *testing.B) {
	engine := getEngine()
	for i := 0; i < b.N; i++ {
		url := "http://127.0.0.1:8088/douyin/favorite/action?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MjI2ODM0NDQ3NzEwODE0NjE3NiwidXNlcm5hbWUiOiJxcSIsImV4cCI6MTY3NzU4MDI2NywiaXNzIjoidGlrdG9rLnVzZXIifQ.zU_SepCAKXJDfWR7wflvDBtb5_pjy8R734Qw33uqydg&video_id=1&action_type=1"
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", url, nil)
		engine.ServeHTTP(recorder, request)
	}
}

func BenchmarkFavoriteList(b *testing.B) {
	engine := getEngine()
	for i := 0; i < b.N; i++ {
		url := "http://127.0.0.1:8088/douyin/favorite/list?user_id=2268344477108146176&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MjI2ODM0NDQ3NzEwODE0NjE3NiwidXNlcm5hbWUiOiJxcSIsImV4cCI6MTY3NzU4MDI2NywiaXNzIjoidGlrdG9rLnVzZXIifQ.zU_SepCAKXJDfWR7wflvDBtb5_pjy8R734Qw33uqydg"
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", url, nil)
		engine.ServeHTTP(recorder, request)
	}
}
