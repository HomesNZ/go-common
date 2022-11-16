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
	conn, err := Conn()
	if err != nil {
		t.Fatal("failed to initialize", err)
	}
	info, err := conn.NodesInfo().Do(ctx)
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	fmt.Println("cluster", info.ClusterName, err)
}

// User is a very simple struct
type User struct {
	Profile string            `json:"profile"`
	Point   *elastic.GeoPoint `json:"point"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings": {
		"properties": {
			"profile": {
				"type": "keyword"
			},
			"location": {
				"type":"geo_point"
			}
		}
	}
}`

func TestPerf(t *testing.T) {
	// docs: https://pkg.go.dev/github.com/olivere/elastic/v7#pkg-overview
	// t.Skip()
	var ctx = context.Background()
	client, err := Conn()
	if err != nil {
		t.Fatal("failed to initialize", err)
	}

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
			panic("create not ack!")
		}
	}
	for i := 0; i < 10; i++ {
		u := User{
			Profile: "olivere createIndex, err := client.CreateIndex(users).BodyString(mapping).Do(ctx) createIndex, err := client.CreateIndex(users).BodyString(mapping).Do(ctx) Flushing a data stream or index is the process of making sure that any data that is currently only stored in the transaction log is also permanently stored in the Lucene index. When restarting, Elasticsearch replays any unflushed operations from the transaction log in to the Lucene index to bring it back into the state that it was in before the restart. Elasticsearch automatically triggers flushes as needed, using heuristics that trade off the size of the unflushed transaction log against the cost of performing each flush olivere createIndex, err := client.CreateIndex(users).BodyString(mapping).Do(ctx) createIndex, err := client.CreateIndex(users).BodyString(mapping).Do(ctx) Flushing a data stream or index is the process of making sure that any data that is currently only stored in the transaction log is also permanently stored in the Lucene index. When restarting, Elasticsearch replays any unflushed operations from the transaction log in to the Lucene index to bring it back into the state that it was in before the restart. Elasticsearch automatically triggers flushes as needed, using heuristics that trade off the size of the unflushed transaction log against the cost of performing each flush",
			Point:   &elastic.GeoPoint{Lat: rand.Float64(), Lon: rand.Float64()},
		}
		put1, err := client.Index().
			Index("users").
			Id(fmt.Sprintf("%d", i+1)).
			BodyJson(u).
			Do(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	}
	_, err = client.DeleteIndex("users").Do(ctx)
	if err != nil {
		panic(err)
	}
}
