package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	elastic "gopkg.in/olivere/elastic.v7"
)

const (
	elasticIndexName = "Customer"
	elasticTypeName  = "user"
)

// User struct
type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}

// CustomerRequest struct
type CustomerRequest struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Age    string `json:"age"`
	Gender string `json:"gender"`
}

// SearchResponse struct
type SearchResponse struct {
	Time     string            `json:"time"`
	Hits     string            `json:"hits"`
	customer []CustomerRequest `json:"customer"`
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
	r.GET("/search", Search)
	if err = r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// CreateCustomer func
func CreateCustomer(c *gin.Context) {
	var docs []CustomerRequest
	if err := c.BindJSON(&docs); err != nil {
		errorResponse(c, http.StatusBadRequest, "Malformed request body")
		return
	}
	ctx := context.Background()
	// Use the IndexExists service to check if a specified index exists.
	exists, err := elasticClient.IndexExists(elasticIndexName).Do(ctx)
	if err != nil {
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := elasticClient.CreateIndex("customer").Do(ctx)
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	// Index a customer (using JSON serialization)
	db := dbConn()
	rows, _ := db.Query("SELECT id, name, age, gender FROM customer")

	user := User{}

	for rows.Next() {
		var id, age int
		var name, gender string

		err = rows.Scan(&id, &name, &age, &gender)
		if err != nil {
			errorResponse(c, http.StatusBadRequest, "Malformed request body")
			return
		}

		user.ID = id
		user.Name = name
		user.Age = age
		user.Gender = gender

		put, err := elasticClient.Index().
			Index(elasticIndexName).
			Type(elasticTypeName).
			Id(strconv.Itoa(id)).
			BodyJson(user).
			Do(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Indexed customer %s to index %s, type %s\n", put.Id, put.Index, put.Type)
	}
	defer db.Close()
}

// Search func
func Search(c *gin.Context) {
	// Parse request
	query := c.Query("query")
	if query == "" {
		errorResponse(c, http.StatusBadRequest, "Query not specified")
		return
	}
	skip := 0
	take := 10
	if i, err := strconv.Atoi(c.Query("skip")); err == nil {
		skip = i
	}
	if i, err := strconv.Atoi(c.Query("take")); err == nil {
		take = i
	}
	// Perform search
	esQuery := elastic.NewMultiMatchQuery(query, "id", "name", "age", "gender").
		Fuzziness("4").
		MinimumShouldMatch("4")
	result, err := elasticClient.Search().
		Index(elasticIndexName).
		Query(esQuery).
		From(skip).Size(take).
		Do(c.Request.Context())
	if err != nil {
		log.Println(err)
		errorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	res := SearchResponse{
		Time: fmt.Sprintf("%d", result.TookInMillis),
		Hits: fmt.Sprintf("%d", result.Hits.TotalHits),
	}
	// Transform search results before returning them
	users := make([]CustomerRequest, 0)
	for _, hit := range result.Hits.Hits {
		var cus CustomerRequest
		json.Unmarshal(*hit.Source, &cus)
		users = append(users, cus)
	}
	res.customer = users
	c.JSON(http.StatusOK, res)
}

func errorResponse(c *gin.Context, code int, err string) {
	c.JSON(code, gin.H{
		"error": err,
	})
}
