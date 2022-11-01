package v7

import (
	"strings"

	"github.com/HomesNZ/env"
	elastic "github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

func Conn() *elastic.Client {
	urls := strings.Split(env.MustGetString("ELASTICSEARCH_URLS"), ";")
	client, err := elastic.NewClient(
		elastic.SetURL(urls...),
		elastic.SetMaxRetries(2),
		// https://github.com/olivere/elastic/wiki/Using-with-AWS-Elasticsearch-Service
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetGzip(false),
	)
	if err != nil {
		panic(err)
	}
	logrus.WithField("package", "elasticsearch").
		Infoln("using olivere/elastic/v7", urls)
	return client
}
