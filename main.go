package main

import (
	"cbt/auth"
	"cbt/handler"
	"fmt"
	"sync"

	// "cbt/helper"
	"cbt/institution"
	"cbt/peserta"
	"cbt/result"
	"cbt/soal"
	"cbt/tmpAnswer"
	"cbt/utils"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	ipLimiter := newIPLimiter()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // URL frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dsn := "bara:" + dbPassword + "@tcp(localhost)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, errDB := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if errDB != nil {
		log.Fatal(errDB.Error())
	}

	authService := auth.NewJwtService()

	resultRepository := result.NewRepository(db)
	resultService := result.NewService(resultRepository)
	resultHandler := handler.NewResultsHandler(resultService)

	tmpAnswerRespository := tmpAnswer.NewRepository(db)
	tmpAnswerService := tmpAnswer.NewService(tmpAnswerRespository)
	tmpAnswerHandler := handler.NewTmpAnswerHandler(tmpAnswerService)

	soalRepository := soal.NewRepository(db)
	soalService := soal.NewService(soalRepository)
	soalHandler := handler.NewSoalHandler(soalService)

	pesertaRepository := peserta.NewRepository(db)
	pesertaService := peserta.NewService(pesertaRepository, authService, tmpAnswerService, soalService, resultService)
	pesertaHandler := handler.NewPesertaHandler(pesertaService, authService)

	institutionRepository := institution.NewRepository(db)
	institutionService := institution.NewService(institutionRepository)
	institutionHandler := handler.NewInstitutionHandler(institutionService)

	router := gin.Default()

	authRouter := router.Group("/v2/auth")
	{
		authRouter.POST("/sessions/users/login", userAgentMiddleware(), rateLimitMiddleware(ipLimiter), pesertaHandler.Login)
		authRouter.POST("/peserta/create", authMiddleware(authService, pesertaService), pesertaHandler.SaveNewPeserta)
		authRouter.POST("/peserta/create/password", authMiddleware(authService, pesertaService), pesertaHandler.SaveNewPeserta2)
		authRouter.POST("/admin/create", pesertaHandler.SaveNewAdmin)

	}

	api := router.Group("/v1")
	api.Use(authMiddleware(authService, pesertaService))

	{
		api.GET("/soal/get-all", soalHandler.GetSoal)
		api.POST("/soal/create", soalHandler.CreateSoal)

		api.GET("/auth/session/get-user", authMiddleware(authService, pesertaService), func(c *gin.Context) {
			c.JSON(200, gin.H{"data": c.MustGet("currentUser").(peserta.GetPesertaTokenInput)})
		})

		api.GET("/auth/session/test-cookie", authMiddleware(authService, pesertaService), func(c *gin.Context) {
			c.JSON(200, gin.H{"data": c.MustGet("currentUser").(peserta.GetPesertaTokenInput)})
		})

		api.GET("/soal/get-by-kode-soal/:kode_soal", soalHandler.GetSoalByKodeSoal)
		api.GET("/peserta/get-by-uuid/:uuid", pesertaHandler.GetByUUID)

		api.POST("/peserta/start/quiz", authMiddleware(authService, pesertaService), pesertaHandler.StartQuiz)
		api.POST("/peserta/finish/quiz/submit", rateLimiter, authMiddleware(authService, pesertaService), pesertaHandler.FinishQuiz)

		api.GET("/peserta/get-all", pesertaHandler.GetAll)
		api.PUT("/peserta/update/:uuid", pesertaHandler.UpdatePeserta)
		api.DELETE("/peserta/delete/:uuid", pesertaHandler.DeletePeserta)
		api.POST("/auth/session/logout", authMiddleware(authService, pesertaService), pesertaHandler.Logout)

		api.GET("/soal/tmp-soal/get-all", soalHandler.GetAllTmpSoal)
		api.GET("/soal/tmp-soal/type/get-all", soalHandler.GetAllTmpSoalByType)
		api.GET("/soal/tmp-soal/type/get-all/student", soalHandler.GetAllTmpSoalByTypeStudent)
		api.POST("/soal/tmp-soal/create", soalHandler.CreateTmpSoal)
		api.DELETE("/soal/tmp-soal/delete/:id", soalHandler.DeleteSoal)
		api.PUT("/soal/tmp-soal/update", soalHandler.UpdateTmpSoal)
		api.POST("/tmp-answer/quiz/save", authMiddleware(authService, pesertaService), tmpAnswerHandler.Create)
		api.GET("/tmp-answer/quiz/get-all", tmpAnswerHandler.GetAll)
		api.GET("/tmp-answer/quiz/student/get-all/answers", tmpAnswerHandler.GetAllAnswerByUUID)
		api.GET("/tmp-answer/quiz/get-by-question-id", tmpAnswerHandler.FindByQuestionID)

		// institution
		api.POST("/institutions/create", institutionHandler.Create)
		api.PUT("/institutions/update", institutionHandler.Update)

		// results
		api.GET("/results/get-all", authMiddleware(authService, pesertaService), resultHandler.GetAllResults)
		api.DELETE("/results/delete/:id", resultHandler.DeleteResult)
		api.PUT("/results/update/:id", resultHandler.UpdateResult)

	}
	errRouter := router.Run(":6121")
	if errRouter != nil {
		return
	}
}

// Rate limiter: 5 request per detik, burst hingga 10
var limiter = rate.NewLimiter(1, 10)

func rateLimiter(c *gin.Context) {

	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests (DDoS detected)"})
		c.Abort()
		return
	}

	c.Next()
}

// Struktur untuk menyimpan rate limiter berdasarkan IP
type ipLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
}

func newIPLimiter() *ipLimiter {
	return &ipLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

func (i *ipLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(1, 5) // 1 request per second dengan burst 5
		i.limiters[ip] = limiter
	}
	return limiter
}

func rateLimitMiddleware(i *ipLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := i.getLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func userAgentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		fmt.Println("User-Agent:", userAgent)
		c.Next()
	}
}

func authMiddleware(authService auth.Service, userService peserta.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenStr, errToken := c.Cookie("session_token")
		if errToken != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - No Token"})
			c.Abort()
			return
		}
		splitToken := strings.Split(tokenStr, "|")
		decodedCookie, err := utils.GetEncryptCookies(splitToken[1])

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - No Token"})
			c.Abort()
			return
		}

		token, err := authService.ValidateToken(decodedCookie)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Invalid Token"})
			c.Abort()
			return
		}
		uuid := string(claim["uuid"].(string))

		user, _ := userService.GetByUUID(uuid)

		c.Set("currentUser", user)

		c.Next()
	}
}

// func authMiddlewareHeaders(authService auth.Service, userService peserta.Service) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if !strings.Contains(authHeader, "Bearer") {
// 			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
// 			return
// 		}

// 		tokenString := ""
// 		arrayToken := strings.Split(authHeader, " ")
// 		if len(arrayToken) == 2 {
// 			tokenString = arrayToken[1]
// 		}

// 		token, err := authService.ValidateToken(tokenString)
// 		if err != nil {
// 			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
// 			return
// 		}
// 		claim, ok := token.Claims.(jwt.MapClaims)
// 		if !ok || !token.Valid {
// 			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
// 			return
// 		}
// 		uuid := string(claim["uuid"].(string))

// 		user, err := userService.GetByUUID(uuid)

// 		if err != nil {
// 			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
// 			return
// 		}

// 		c.Set("currentUser", user)

// 	}
// }
