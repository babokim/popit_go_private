package main

import (
	"github.com/labstack/echo"
	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/labstack/echo/middleware"
	"net/http"
	"log"
	"os"
	"time"
	"context"
)

var (
	xormDb *xorm.Engine
)

func init() {
	dbConn := os.Getenv("DB_CONN")
	if len(dbConn) == 0 {
		dbConn = "root:@tcp(127.0.0.1:3306)/wordpress?charset=utf8&parseTime=True"
	}

	db, err := xorm.NewEngine("mysql", dbConn)
	if err != nil {
		panic(fmt.Errorf("Database open error: %s \n", err))
	}
	db.ShowSQL(false)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(60 * time.Hour)

	xormDb = db

}

type ApiResult struct {
	Data  interface{} 	`json:"data"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
}

func main() {
	defer xormDb.Close()

	e := echo.New()
	//renderer := &TemplateRenderer{
	//	templates: template.Must(template.ParseGlob("./views/*.html")),
	//}
	//e.Renderer = renderer


	e.Pre(middleware.RemoveTrailingSlash())
	//e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(setDbConnContext(xormDb))
	//e.Static("/", "public")

	e.GET("/api/RecentPosts", GetRecentPosts)
	e.GET("/api/ChannelPosts", GetChannelPosts)
	e.GET("/api/RandomAuthorPosts", GetRandomAuthorPosts)

	log.Fatal(e.Start(":8000"))
}

func setDbConnContext(xormDb *xorm.Engine) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			session := xormDb.NewSession()
			defer session.Close()

			req := ctx.Request()
			ctx.SetRequest(req.WithContext(context.WithValue(req.Context(), "DB", session)))
			return next(ctx)
		}
	}
}

func GetDBConn(ctx context.Context) *xorm.Session {
	v := ctx.Value("DB")
	if v == nil {
		panic("DB is not exist")
	}
	if db, ok := v.(*xorm.Session); ok {
		return db
	}
	panic("DB is not exist")
}

func GetRecentPosts(c echo.Context) error {
	posts, err := Post{}.GetRecent(c.Request().Context())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ApiResult{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, ApiResult{
		Success: true,
		Data: posts,
		Message: "",
	})
}

func GetChannelPosts(c echo.Context) error {
	posts, err := Post{}.GetRandomPostsByTerm(c.Request().Context())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ApiResult{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, ApiResult{
		Success: true,
		Data: posts,
		Message: "",
	})
}

func GetRandomAuthorPosts(c echo.Context) error {
	posts, err := Post{}.GetRandomPostsByAuthor(c.Request().Context())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ApiResult{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, ApiResult{
		Success: true,
		Data: posts,
		Message: "",
	})
}

//func GetAvatar(c echo.Context) error {
//	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
//	if err != nil {
//		return c.JSON(http.StatusInternalServerError, ApiResult{
//			Success: false,
//			Message: err.Error(),
//		})
//	}
//
//	author, err := Author{}.GetOne(id)
//	if err != nil {
//		return c.JSON(http.StatusInternalServerError, ApiResult{
//			Success: false,
//			Message: err.Error(),
//		})
//	}
//
//
//	avatar, err :=  author.GetAvatar()
//	if err != nil {
//		return c.JSON(http.StatusInternalServerError, ApiResult{
//			Success: false,
//			Message: err.Error(),
//		})
//	}
//
//	return c.JSON(http.StatusOK, ApiResult{
//		Success: true,
//		Data: avatar,
//		Message: "",
//	})
//}