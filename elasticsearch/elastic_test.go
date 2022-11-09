package elasticsearch

import (
	"context"
	"fmt"
	"testing"
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
