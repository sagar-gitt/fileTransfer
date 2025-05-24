package main

import (
	"database/sql"
	"fileTransfer/internal/config"
	"fileTransfer/internal/handlers"
	"fileTransfer/internal/repository"
	"fileTransfer/internal/utils"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	//Loading .env file for credentials
	config.LoadEnv()

	//Connecting to MySQL
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/filetransfer"
	} else {
		// Parse Railway.app connection string
		if strings.HasPrefix(dsn, "mysql://") {
			parsedURL, err := url.Parse(dsn)
			if err != nil {
				log.Fatal("Error parsing MySQL connection string: ", err)
			}

			user := parsedURL.User.Username()
			password, _ := parsedURL.User.Password()
			host := parsedURL.Host
			dbName := strings.TrimPrefix(parsedURL.Path, "/")

			dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbName)
		}
	}

	db := ConnectMySQL(dsn)

	//Creating tables and inserting sample data
	mySqlInit := repository.NewMySQLInitRepo(db)
	err := mySqlInit.CreateUserTableIfNotExist()
	if err != nil {
		log.Fatal("Error Creating User Table: ", err)
	}

	err = mySqlInit.CreateFileTableIfNotExist()
	if err != nil {
		log.Fatal("Error Creating File Table: ", err)
	}

	err = mySqlInit.InsertSampleData()
	if err != nil {
		log.Fatal("Error Inserting Sample Data: ", err)
	}

	//Truncating tables when needed.
	//err := mySqlInit.TruncateAllTables()
	//if err != nil {
	//	log.Fatal("error Truncating All Tables: ", err)
	//}

	//Initializing User Repo-MySQL
	mysqlUserRepo := repository.NewMysqlUserRepo(db)
	mysqlFileRepo := repository.NewMysqlFileRepo(db)

	//Initializing Google Oauth2
	handlers.InitGoogleAuth()

	//Initializing JWT Service
	jwt := utils.NewJWTService()

	//Initializing AWS S3 Service
	awsS3 := utils.NewAwsS3()

	//Initializing Handlers
	h := handlers.NewHandlers(mysqlUserRepo, mysqlFileRepo, jwt, awsS3)

	//Go Routine that deletes the expired AWS files
	go func() {
		for {
			awsS3.DeleteExpiredFiles(mysqlFileRepo)
			time.Sleep(1 * time.Hour) // Run every hour
		}
	}()

	//Creating Gin based Routes
	r := gin.Default()

	// CORS configuration
	allowedOrigins := []string{"http://localhost:5173", "https://your-frontend-url.vercel.app"}
	if os.Getenv("ALLOWED_ORIGINS") != "" {
		allowedOrigins = append(allowedOrigins, os.Getenv("ALLOWED_ORIGINS"))
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})

	authRoutes := r.Group("/auth/google")
	{
		authRoutes.GET("/login", h.GoogleLogin)
		authRoutes.GET("/callback", h.GoogleCallback)
	}

	fileRoutes := r.Group("/file")
	{
		fileRoutes.POST("/upload", h.UploadFileAndSaveInfo)
		fileRoutes.GET("/download", h.DownloadFile)
		fileRoutes.GET("/listFiles", h.ListFile)
		fileRoutes.POST("/sendEmail", h.SendFileDownloadLink)
	}

	//Starting the server
	err = r.Run(":8080")
	if err != nil {
		return
	}
}

func ConnectMySQL(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error Connecting to Database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error Pinging Database: ", err)
	}

	log.Println("Successfully connected to MySQL")
	return db
}
