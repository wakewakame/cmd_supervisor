package main

import (
  "os"
  "fmt"
  "log"
  "time"
  "net/http"
  "encoding/base64"
  "gorm.io/gorm"
  "gorm.io/driver/postgres"
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/redis"
)

func main() {
  router := gin.Default()

  // セッション管理用のミドルウェアを登録
  sessionSecretBase64 := os.Getenv("SESSION_SECRET")
  if sessionSecretBase64 == "" { log.Fatal("error: Environment variable SESSION_SECRET is not set.") }
  sessionSecret, err := base64.StdEncoding.DecodeString(sessionSecretBase64)
  if err != nil { log.Fatal("error: Environment variable SESSION_SECRET is not base64.") }
  store, _ := redis.NewStore(10, "tcp", "kvs:6379", "", sessionSecret)
  router.Use(sessions.Sessions("session", store))

  // DBに接続
  dsn := fmt.Sprint(
    "host=db ",
    "user=", os.Getenv("POSTGRES_USER"), " ",
    "password=", os.Getenv("POSTGRES_PASSWORD") ," ",
    "dbname=", os.Getenv("POSTGRES_DB"), " ",
    "port=5432 ",
    "sslmode=disable ",
    "TimeZone=Asia/Shanghai",
  )
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil { log.Fatal("error: Failed to connect db") }

  // マイグレーション
  type User struct {
    Id uint `gorm:"primaryKey;autoIncrement"`
    Name string `gorm:"not null"`
    UserId string `gorm:"not null;unique"`
    HashedPassword string `gorm:"not null"`
  }
  type Commands struct {
    Id uint `gorm:"primaryKey;autoIncrement"`
    CreatedById uint `gorm:"not null;column:created_by"`
    CreatedBy User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    StartedAt time.Time `gorm:"not null"`
    FinishedAt time.Time
    Command string `gorm:"not null"`
    Output string
    ExitCode int
  }
  err = db.AutoMigrate(&User{})
  if err != nil { log.Fatal(err) }
  err = db.AutoMigrate(&Commands{})
  if err != nil { log.Fatal(err) }
  db.Create(&User{Id: 0, Name: "hoge", UserId: "abc123", HashedPassword: "def456"})

  // 静的ファイルのホスティング
  router.Static("/static", "static/")
  router.GET("/", func(c *gin.Context) {
    c.Redirect(http.StatusMovedPermanently, "/static/index.html")
  })

  // apiサーバ
  api := router.Group("/api")

  auth := func(c *gin.Context) {
    session := sessions.Default(c)
    id, ok := session.Get("userId").(int)
    if !ok {
      c.JSON(http.StatusUnauthorized, gin.H{"result": "error", "description": "unauthorized"})
      c.Abort()
      return
    }
    c.Set("userId", id)
    c.Next()
  }
  api.POST("/user/login", func(c *gin.Context) {
    var json struct {
      UserId string `json:"user_id"`
      Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&json); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"result": "error", "description": "invalid json"})
      c.Abort()
      return
    }
    if !(json.UserId == "abc123" && json.Password == "hogepiyo") {
      c.JSON(http.StatusForbidden, gin.H{"result": "error", "description": "incorrect id or password"})
      c.Abort()
      return
    }

    session := sessions.Default(c)
    session.Set("userId", 123)
    session.Save()
    c.JSON(http.StatusOK, gin.H{"result": "success"})
  })
  api.POST("/user/logout", func(c *gin.Context) {
    session := sessions.Default(c)
    session.Delete("userId")
    session.Save()
    c.JSON(http.StatusOK, gin.H{"result": "success"})
  })
  api.GET("/user/me", auth, func(c *gin.Context) {
    id, _ := sessions.Default(c).Get("userId").(int)
    c.JSON(http.StatusOK, gin.H{"id": id})
  })

  router.NoRoute(func(c *gin.Context) {
    c.JSON(http.StatusNotFound, gin.H{
      "error": "not found",
    })
  })
  router.Run(":8080")
}
