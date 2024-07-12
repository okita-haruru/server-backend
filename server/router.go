package server

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"sushi/utils/config"
	"sushi/utils/ratelimit"
)

func NewRouter(server *Server, conf config.Config, socketserver *socketio.Server) *gin.Engine {
	gin.SetMode(server.config.GinMode())
	r := gin.Default()

	r.Use(ratelimit.GinMiddleware())
	r.Use(CORSMiddleware())
	r.Use(static.Serve("/", static.LocalFile("/app/Demo-UI", true)))

	//public
	r.GET("/ping", server.controller.HandlePing)
	r.GET("/player_list", server.controller.HandleGetPlayerList)
	r.GET("/ranking/balance", server.controller.HandleGetBalanceRanking)
	r.GET("/ranking/kills/total", server.controller.HandleGetPlayerKillStatsSortByTotal)
	r.GET("/ranking/kills/warden", server.controller.HandleGetPlayerKillStatsSortByWarden)
	r.GET("/ranking/kills/ender_dragon", server.controller.HandleGetPlayerKillStatsSortByEnderDragon)
	r.GET("/ranking/kills/wither", server.controller.HandleGetPlayerKillStatsSortByWither)
	r.GET("/ranking/kills/piglin", server.controller.HandleGetPlayerKillStatsSortByPiglin)
	r.GET("/ranking/kills/phantom", server.controller.HandleGetPlayerKillStatsSortByPhantom)
	r.GET("/ranking/fish_amount", server.controller.HandleGetFishAmount)
	r.GET("/ranking/fish_size", server.controller.HandleGetFishSize)
	r.GET("/profile", server.controller.HandleGetPlayerProfileByName)
	v1 := r.Group("/v1")
	authorizedV1 := v1.Group("/")
	authorizedV1.Use(server.GetAuth())
	WithUserRoutes(authorizedV1, server, conf)
	return r
}

func WithTeamRoutes(r *gin.RouterGroup, server *Server) {
	//r.GET("/", server.controller.team.HandleTeamList)
}

func WithUserRoutes(r *gin.RouterGroup, server *Server, conf config.Config) {
	//authorized := r
	//authorized.POST("/users/profile")
}

func (server Server) GetAuth() gin.HandlerFunc {

	return func(c *gin.Context) {
		//AUTH
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Auth-Token, Authorization, Code, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT , PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
