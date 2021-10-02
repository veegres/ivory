package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Cluster struct {
	Name  string   `json:"name"`
	Nodes []string `json:"nodes"`
}

var ClusterBucket = []byte("Cluster")
var DB = setupDatabase()

func main() {
	// todo decide what to do with close
	defer DB.Close()

	// todo move it to some config folder
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/ping", func(context *gin.Context) { context.JSON(http.StatusOK, gin.H{"message": "pong"}) })
		api.GET("/cluster/list", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{"response": ClusterGetList()})
		})
		api.POST("/cluster/create", func(context *gin.Context) {
			var cluster Cluster
			context.ShouldBindJSON(&cluster)
			ClusterUpdate(cluster)
			context.JSON(http.StatusOK, gin.H{"response": cluster})
		})
		api.PUT("/cluster/:name/edit", func(context *gin.Context) {
			name := context.Param("name")
			var nodes []string
			context.ShouldBindJSON(&nodes)
			ClusterUpdate(Cluster{Name: name, Nodes: nodes})
		})
		// TODO make code independent of patroni
		/* PATRONI */
		api.GET("/node/:name/cluster", func(context *gin.Context) {
			name := context.Param("name")
			response, err := http.Get("http://" + name + ":8008/cluster")
			if err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err})
			} else {
				var body interface{}
				json.NewDecoder(response.Body).Decode(&body)
				context.JSON(http.StatusOK, gin.H{"response": body})
			}
		})
		api.GET("/node/:name/patroni", func(context *gin.Context) {
			name := context.Param("name")
			response, err := http.Get("http://" + name + ":8008/patroni")
			if err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err})
			} else {
				var body interface{}
				json.NewDecoder(response.Body).Decode(&body)
				context.JSON(http.StatusOK, gin.H{"response": body})
			}
		})
		api.GET("/node/:name/config", func(context *gin.Context) {
			name := context.Param("name")
			response, err := http.Get("http://" + name + ":8008/config")
			if err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err})
			} else {
				var body interface{}
				json.NewDecoder(response.Body).Decode(&body)
				context.JSON(http.StatusOK, gin.H{"response": body})
			}
		})
		api.PATCH("/node/:name/config", func(context *gin.Context) {
			name := context.Param("name")
			body, _ := ioutil.ReadAll(context.Request.Body)
			req, _ := http.NewRequest(http.MethodPatch, "http://"+name+":8008/config", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))
			response, err := http.DefaultClient.Do(req)
			if err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err})
			} else {
				var body interface{}
				json.NewDecoder(response.Body).Decode(&body)
				context.JSON(http.StatusOK, gin.H{"response": body})
			}
		})
	}

	// todo how can we extract it
	router.Run()
}

func ClusterUpdate(cluster Cluster) {
	DB.Update(func(tx *bolt.Tx) error {
		var buff bytes.Buffer
		gob.NewEncoder(&buff).Encode(cluster.Nodes)
		return tx.Bucket(ClusterBucket).Put([]byte(cluster.Name), buff.Bytes())
	})
}

func ClusterGetList() []Cluster {
	var list []Cluster
	DB.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(ClusterBucket).Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			r := bytes.NewBuffer(v)
			var nodes []string
			gob.NewDecoder(r).Decode(&nodes)
			list = append(list, Cluster{Name: string(k), Nodes: nodes})
		}
		return nil
	})
	return list
}

func setupDatabase() *bolt.DB {
	db, err := bolt.Open("bolt/cluster.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte("Cluster"))
		return nil
	})

	return db
}
