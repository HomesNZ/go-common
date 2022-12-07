# go-common - elasticsearch / opensearch

As of Nov 2022 transitioning to AWS opensearch 1.3.  
For new projects import:

```go
elasticsearch "github.com/HomesNZ/go-common/elasticsearch/v7"
github.com/olivere/elastic/v7
```

## Tests

Tests both for this lib and code that depends on it can be run like so:

Open one terminal with Opensearch: `docker run --privileged=true -p 9200:9200 -p 9600:9600 -e "discovery.type=single-node" -e "plugins.security.disabled=true" opensearchproject/opensearch:1.3.1`

On another terminal:

```
export ELASTICSEARCH_INSECURE=true
export ELASTICSEARCH_URLS="http://127.0.0.1:9200"

go test ./...
```
