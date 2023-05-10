package elasticsearch

import (
	"net/http"
	"strings"
	"sync"

	"github.com/HomesNZ/elastic"
	"github.com/HomesNZ/go-common/env"
	"github.com/sirupsen/logrus"
	awsauth "github.com/smartystreets/go-aws-auth"
)

var (
	conn     *elastic.Client
	initOnce = sync.Once{}
)

func awsAuth() bool {
	key := env.GetString("AWS_ACCESS_KEY_ID", "")
	secret := env.GetString("AWS_SECRET_ACCESS_KEY", "")
	token := env.GetString("AWS_SECURITY_TOKEN", "")
	return key != "" && secret != "" || token != ""
}

func initConn() {
	// Create a client
	var err error
	options := []elastic.ClientOptionFunc{
		elastic.SetURL(strings.Split(env.GetString("ELASTICSEARCH_URLS", ""), ";")...),
		elastic.SetHealthcheck(env.GetBool("ELASTICSEARCH_HEALTH_CHECK", true)),
		elastic.SetSniff(env.GetBool("ELASTICSEARCH_SNIFF", false)), // causes issues within AWS, so off by default
	}

	if retries := env.GetInt("ELASTICSEARCH_MAX_RETRIES", 0); retries > 0 {
		options = append(options, elastic.SetMaxRetries(retries))
	}

	if awsAuth() {
		options = append(options, elastic.SetPrepareRequest(func(req *http.Request) {
			awsauth.Sign(req)
		}))
	}
	conn, err = elastic.NewClient(options...)
	if err != nil {
		// Handle error
		panic(err)
	}

	logrus.WithField("package", "elasticsearch").
		Info("elasticsearch.Conn() is deprecated, use elastic.NewClient() instead")
}

// Conn returns a connection to ElasticSearch
func Conn() *elastic.Client {
	initOnce.Do(initConn)
	return conn
}
