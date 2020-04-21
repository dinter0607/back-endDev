package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	elastic "gopkg.in/olivere/elastic.v7"
)

const (
	elasticIndexName = "customer"
	elasticTypeName  = "user"
)

// User struct
type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
	Wallet int    `json:"wallet"`
}

var (
	elasticClient *elastic.Client
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "postgres"
)

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
func main() {
	var err error
	// Create Elastic client and wait for Elasticsearch to be ready
	for {
		elasticClient, err = elastic.NewClient(
			elastic.SetURL("http://localhost:9200"),
			elastic.SetSniff(false),
		)
		if err != nil {
			log.Println(err)
			// Retry every 3 seconds
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}
	// Start HTTP server
	r := gin.Default()
	r.POST("/customer", CreateCustomer)
	r.GET("/search", GetUserByID)
	r.PATCH("/update", Update)
	if err = r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// CreateCustomer func
func CreateCustomer(c *gin.Context) {
	ctx := context.Background()
	// Use the IndexExists service to check if a specified index exists.
	exists, err := elasticClient.IndexExists(elasticIndexName).Do(ctx)
	if err != nil {
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := elasticClient.CreateIndex(elasticIndexName).Do(ctx)
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	// Index a customer (using JSON serialization)
	db := dbConn()
	rows, _ := db.Query("SELECT id, name, age, gender, wallet FROM customer")

	user := User{}

	for rows.Next() {
		var id, age, wallet int
		var name, gender string

		err = rows.Scan(&id, &name, &age, &gender, &wallet)
		if err != nil {
			errorResponse(c, http.StatusBadRequest, "Lỗi rồi em ơi :)))")
			return
		}

		user.ID = id
		user.Name = name
		user.Age = age
		user.Gender = gender
		user.Wallet = wallet

		put, err := elasticClient.Index().
			Index(elasticIndexName).
			Type(elasticTypeName).
			Id(strconv.Itoa(id)).
			BodyJson(user).
			Do(ctx)
		if err != nil {
			panic(err)
		}
		c.JSON(200, gin.H{
			"Indexed customer": put.Id,
			"Name":             name,
			"Age":              age,
			"Gender":           gender,
			"Wallet":           wallet,
		})
	}
	defer db.Close()
}

// GetUserByID func
func GetUserByID(c *gin.Context) {
	userID := c.Query("id")
	id, _ := strconv.Atoi(userID)

	ctx := context.Background()
	// Search with a term query
	termQuery := elastic.NewTermQuery("id", id)
	searchResult, err := elasticClient.Search().
		Index(elasticIndexName).
		Query(termQuery).
		Pretty(true).
		Do(ctx)
	if err != nil {
		panic(err)
	}

	var ttyp User
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(User); ok {
			c.JSON(200, gin.H{
				"ID":     t.ID,
				"Name":   t.Name,
				"Age":    t.Age,
				"Gender": t.Gender,
				"Wallet": t.Wallet,
			})
		}
	}
	c.JSON(200, gin.H{
		"Found a total of customer": searchResult.TotalHits(),
	})

}

// Update func
func Update(c *gin.Context) {
	userID := c.Query("id")
	Wallet := c.Query("wallet")

	ctx := context.Background()
	// Update a user by the update API of Elasticsearch.
	update, err := elasticClient.Update().Index(elasticIndexName).Type(elasticTypeName).Id(userID).
		Script(elastic.NewScript("ctx._source.wallet = params.wallet").Lang("painless").Param("wallet", Wallet)).
		Upsert(map[string]interface{}{"age": 0}).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	c.JSON(200, gin.H{
		"New version of user": update.Id,
		"Version":             update.Version,
	})
	// Update user at Postgres
	db := dbConn()
	sqlStatement := `UPDATE customer SET wallet =` + Wallet + `WHERE id =` + userID
	_, err = db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	CreateCustomer(c)

}
func errorResponse(c *gin.Context, code int, err string) {
	c.JSON(code, gin.H{
		"error": err,
	})
}

// }
