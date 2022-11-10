package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/samber/lo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func IsNameExists(s []Creds, name string) bool {
	creds := lo.Filter(s, func(item Creds, index int) bool {
		fmt.Println(item.Name)
		if strings.EqualFold (strings.ToLower(item.Name), strings.ToLower(name)){
			return true
		} else {
			return false
		}
	})
	if len(creds) != 0 {
		return true
	} else {
		return false
	}
}

type Creds struct {
	gorm.Model
	Name string
	Password string
}

func CorsMiddleware() gin.HandlerFunc {
	return func (ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
	}
}

func main() {	
	db, err := gorm.Open(postgres.Open("dbname=gormtest user=postgres password=postgres port=5432"))
	db.AutoMigrate(&Creds{})
	if err != nil {
		panic(err)
	}
	// load .env file
	godotenv.Load()
	
	if !(strings.EqualFold(os.Getenv("GIN_ENV"), "development")) {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(CorsMiddleware())
	r.SetTrustedProxies([]string{})
	r.POST("/submit", func(ctx *gin.Context) {
		r := ctx.Request
		// get formdata value from context's request
		name := r.FormValue("name")
		pw := r.FormValue("password")

		var users []Creds
		db.Select("name").Find(&users)

		// create data if name doesn't exist (gorm doesn't have unique constraint support)
		if !IsNameExists(users, name) {
			db.Create(&Creds{Name: name, Password: pw})
			db.Commit()
			ctx.String(200, "OK")
		} else {
			fmt.Println("Name already exists")
			ctx.String(400, "Name already exists")
		}
	})
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello World")
	})
	fmt.Println("Server listening at port 9000")
	http.ListenAndServe("localhost:9000", r)
}
