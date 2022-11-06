package v7

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/olivere/elastic/v7"
)

func TestConn(t *testing.T) {
	var ctx = context.Background()
	conn := Conn()
	info, err := conn.NodesInfo().Do(ctx)
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	fmt.Println("cluster", info.ClusterName, err)
}

// User is a very simple struct
type User struct {
	User  string            `json:"user"`
	Point *elastic.GeoPoint `json:"point"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"user":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"location":{
					"type":"geo_point"
				}
			}
		}
	}
}`

func TestPerf(t *testing.T) {
	t.Skip()
	var ctx = context.Background()
	client := Conn()

	exists, err := client.IndexExists("users").Do(ctx)
	if err != nil {
		panic(err)
	}
	if !exists {
		createIndex, err := client.CreateIndex("users").BodyString(mapping).Do(ctx)
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
			panic(err)
		}
	}
	for i := 0; i < 20; i++ {
		// todo add a point here
		u := User{User: "olivere", Point: &elastic.GeoPoint{Lat: rand.Float64(), Lon: rand.Float64()}}
		put1, err := client.Index().
			Index("users").
			Type("user").
			// Id("1").
			BodyJson(u).
			Do(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	}
	_, err = client.DeleteIndex("users").Do(ctx)
	if err != nil {
		panic(err)
	}
}
