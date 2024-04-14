package http_server

import (
	"awesomeProject/config"
	get "awesomeProject/internal/app/acc/get-funcs"
	post "awesomeProject/internal/app/acc/post-funcs"
	getHandle "awesomeProject/internal/app/handlers/get"
	postHandle "awesomeProject/internal/app/handlers/post"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
)

type Server struct {
	Cfg         *config.Config
	Router      *gin.Engine
	Db          *sql.DB
	RedisClient *redis.Client
}

func NewServer(cfg *config.Config) *Server {
	db, rd, err := InitDb(cfg)
	if err != nil {
		log.Fatalf("Cannot create new db. Error: {%s}", err)
	}
	return &Server{
		Cfg:         cfg,
		Router:      gin.Default(),
		Db:          db,
		RedisClient: rd,
	}
}

type App struct {
	post       *post.Post
	get        *get.Get
	postHandle *postHandle.PostF
	getHandle  *getHandle.GetF
}

func (s *Server) Run() error {
	app := &App{}
	app.post = post.New()
	app.get = get.New()
	app.postHandle = postHandle.New(app.post)
	app.getHandle = getHandle.New(app.get)

	s.Router.GET("/user_banner", app.getHandle.GetBannerHandler(s.Db, s.RedisClient)) // OK http://localhost:8080/user_banner?tag_id=8&feature_id=15&use_last_revision=true&admin_token=25
	s.Router.POST("/set_admin", app.postHandle.PostAdminStateHandler(s.Db))           // OK http://localhost:8080/set_admin?id=25&state=true
	s.Router.POST("/banner", app.postHandle.CreateNewBannerHandler(s.Db))             // OK http://localhost:8080/banner?admin_token=10&feature_id=15&tag_ids=22,12&content=notebooklovers,simpledescr,http://aboba.com&is_active=true
	s.Router.PATCH("/banner/{id}", app.postHandle.ChangeBannerHandler(s.Db))          // OK http://localhost:8080/banner/{id}?admin_token=25&feature_id=100&tag_ids=100,101&content=avitolovers,descr,http://avito.com&is_active=true&id=3
	s.Router.DELETE("/banner/{id}", app.postHandle.DeleteBannerHandler(s.Db))         // OK http://localhost:8080/banner/{id}?admin_token=10&id=3
	s.Router.GET("/banner", app.getHandle.GetBannerByFilterHandler(s.Db))             // ОК http://localhost:8080/banner?admin_token=25&feature_id=15&tag_id=15&content=5&offset=0
	s.Router.DELETE("/delete", app.postHandle.DeleteFeatureTagHandler(s.Db))          // ОК http://localhost:8080/delete?admin_token=25&feature_id=15 или http://localhost:8080/delete?admin_token=25&tag_id=15&content=1&offset=1
	s.Router.GET("/history", app.getHandle.GetBannersHistoryHandler(s.Db))            // ОК http://localhost:8080/history?admin_token=25&id=15
	s.Router.PATCH("/history", app.postHandle.ChangeBannersHistoryHandler(s.Db))      // OK http://localhost:8080/history?admin_token=25&id=15&number=8

	if err := s.Router.Run(s.Cfg.Server.Port); err != nil {
		log.Fatalf("Cannot listen. Error: {%s}", err)
	}

	return nil
}
