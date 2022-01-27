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
	"os"
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

		/* LOCAL DB */
		/* CLUSTER */
		api.GET("/cluster", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{"response": ClusterGetList()})
		})
		api.GET("/cluster/:name", func(context *gin.Context) {
			name := context.Param("name")
			cluster := ClusterGet(name)
			if cluster == nil {
				context.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			} else {
				context.JSON(http.StatusOK, gin.H{"response": cluster})
			}
		})
		api.PUT("/cluster", func(context *gin.Context) {
			var cluster Cluster
			context.ShouldBindJSON(&cluster)
			ClusterUpdate(cluster)
			context.JSON(http.StatusOK, gin.H{"response": cluster})
		})
		api.DELETE("/cluster/:name", func(context *gin.Context) {
			name := context.Param("name")
			ClusterDelete(name)
		})

		// TODO make code independent of patroni
		/* PROXY */
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

func ClusterDelete(key string) {
	DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(ClusterBucket).Delete([]byte(key))
	})
}

func ClusterUpdate(cluster Cluster) {
	DB.Update(func(tx *bolt.Tx) error {
		var buff bytes.Buffer
		gob.NewEncoder(&buff).Encode(cluster.Nodes)
		return tx.Bucket(ClusterBucket).Put([]byte(cluster.Name), buff.Bytes())
	})
}

func ClusterGet(key string) *Cluster {
	var cluster *Cluster
	DB.View(func(tx *bolt.Tx) error {
		value := tx.Bucket(ClusterBucket).Get([]byte(key))
		if value == nil {
			return nil
		}
		buff := bytes.NewBuffer(value)
		var nodes []string
		gob.NewDecoder(buff).Decode(&nodes)
		tmp := Cluster{Name: key, Nodes: nodes}
		cluster = &tmp
		return nil
	})
	return cluster
}

func ClusterGetList() []Cluster {
	list := make([]Cluster, 0)
	DB.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(ClusterBucket).Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			buff := bytes.NewBuffer(value)
			var nodes []string
			gob.NewDecoder(buff).Decode(&nodes)
			list = append(list, Cluster{Name: string(key), Nodes: nodes})
		}
		return nil
	})
	return list
}

func setupDatabase() *bolt.DB {
	os.Mkdir("bolt", os.ModePerm)
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
