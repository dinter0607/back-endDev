package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "postgres"
)

// User struct
type User struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
	FullName string `json:"full_name"`
}

func dbConn() (db *sql.DB) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	return db
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Static("/public", "./public")

	client := r.Group("/api")
	{
		client.GET("/user/:id", Read)
		client.POST("/user/create", Create)
		client.PATCH("/user/update/:id", Update)
		client.DELETE("/user/:id", Delete)
	}

	return r
}

// Read func
func Read(c *gin.Context) {

	db := dbConn()
	rows, err := db.Query("SELECT id, user_name, full_name FROM account WHERE id = " + c.Param("id"))

	if err != nil {
		c.JSON(500, gin.H{
			"messages": "User  not found",
		})
	}

	user := User{}

	for rows.Next() {
		var id int
		var user_name, full_name string

		err = rows.Scan(&id, &user_name, &full_name)
		if err != nil {
			panic(err.Error())
		}

		user.ID = id
		user.UserName = user_name
		user.FullName = full_name
	}
	c.JSON(200, user)
	defer db.Close() // Hoãn lại việc close database connect cho đến khi hàm Read() thực hiệc xong
}

// Create func
func Create(c *gin.Context) {
	db := dbConn()

	type CreateUser struct {
		UserName string `form:"user_name" json:"user_name" binding:"required"`
		FullName string `form:"full_name" json:"full_name" binding:"required"`
	}

	var json CreateUser

	if err := c.ShouldBindJSON(&json); err == nil {
		insUser, err := db.Prepare("INSERT INTO account(user_name, full_name) VALUES($1, $2)")
		if err != nil {
			c.JSON(500, gin.H{
				"messages": err,
			})
		}

		insUser.Exec(json.UserName, json.FullName)
		c.JSON(200, gin.H{
			"messages": "Thêm mới thành công",
		})

	} else {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	defer db.Close()
}

// Update func
func Update(c *gin.Context) {
	db := dbConn()

	type UpdateUser struct {
		UserName string `form:"user_name" json:"user_name" binding:"required"`
		FullName string `form:"full_name" json:"full_name" binding:"required"`
	}

	var json UpdateUser
	if err := c.ShouldBindJSON(&json); err == nil {
		edit, err := db.Prepare("UPDATE account SET user_name=$1, full_name=$2 WHERE id= " + c.Param("id"))
		if err != nil {
			panic(err.Error())
		}
		edit.Exec(json.UserName, json.FullName)

		c.JSON(200, gin.H{
			"messages": "Cập nhật thành công",
		})
	} else {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	defer db.Close()
}

// Delete func
func Delete(c *gin.Context) {
	db := dbConn()

	delete, err := db.Prepare("DELETE FROM account WHERE id= " + c.Param("id"))
	if err != nil {
		// panic(err.Error())
		fmt.Println(err)
	}

	delete.Exec(c.Param("id"))
	c.JSON(200, gin.H{
		"messages": "Xóa thành công",
	})

	defer db.Close()
}
func main() {
	r := setupRouter()
	r.Run(":8080") // Ứng dụng chạy tại cổng 8080
}
