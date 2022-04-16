//+build example

package search

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/exluap/kit"
	"github.com/exluap/kit/log"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type Tweet struct {
	User     string
	Message  string
	Retweets int64
}

func Test_Simple(t *testing.T) {
	// SetSniff(false) used to worked properly with docker environment on the local machine(macOS)
	opts := []elastic.ClientOptionFunc{elastic.SetSniff(false)}
	// Obtain a client. You can also provide your own HTTP client here.
	client, err := elastic.NewClient(opts...)
	// Trace request and response details like this
	// client, err := elastic.NewClient(elastic.SetTraceLog(log.New(os.Stdout, "", 0)))
	if err != nil {
		// Handle error
		t.Fatal(err)
	}

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("http://127.0.0.1:9200").Do(context.Background())
	if err != nil {
		// Handle error
		t.Fatal(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
	if err != nil {
		// Handle error
		t.Fatal(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("twitter").Do(context.Background())
	if err != nil {
		// Handle error
		t.Fatal(err)
	}

	if exists {
		// Delete an index.
		_, err := client.DeleteIndex("twitter").Do(context.Background())
		if err != nil {
			// Handle error
			t.Fatal(err)
		}
	}

	// Create a new index.
	mapping := `
{
	"mappings":{
		"properties":{
			"user":{
				"type":"keyword"
			},
			"message":{
				"type":"text"
			},
			"retweets":{
				"type":"long"
			}
		}
	}
}
`
	_, err = client.CreateIndex("twitter").Body(mapping).Do(context.Background())
	if err != nil {
		// Handle error
		t.Fatal(err)
	}

	// Index a tweet (using JSON serialization)
	tweet1 := Tweet{User: "olivere", Message: "Take Five", Retweets: 0}
	put1, err := client.Index().
		Index("twitter").
		Id("1").
		BodyJson(tweet1).
		Do(context.Background())
	if err != nil {
		// Handle error
		t.Fatal(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	// Index a second tweet (by string)
	tweet2 := `{"user" : "olivere", "message" : "It's a Raggy Waltz"}`
	put2, err := client.Index().
		Index("twitter").
		Id("2").
		BodyString(tweet2).
		Do(context.Background())
	if err != nil {
		// Handle error
		t.Fatal(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put2.Id, put2.Index, put2.Type)

	// Get tweet with specified ID
	get1, err := client.Get().
		Index("twitter").
		Id("1").
		Do(context.Background())
	if err != nil {
		switch {
		case elastic.IsNotFound(err):
			t.Fatal(fmt.Sprintf("Document not found: %v", err))
		case elastic.IsTimeout(err):
			t.Fatal(fmt.Sprintf("Timeout retrieving document: %v", err))
		case elastic.IsConnErr(err):
			t.Fatal(fmt.Sprintf("Connection problem: %v", err))
		default:
			// Some other kind of error
			t.Fatal(err)
		}
	}
	fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)

	// Refresh to make sure the documents are searchable.
	_, err = client.Refresh().Index("twitter").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// Search with a term query
	termQuery := elastic.NewTermQuery("user", "olivere")
	searchResult, err := client.Search().
		Index("twitter").        // search in index "twitter"
		Query(termQuery).        // specify the query
		Sort("user", true).      // sort by "user" field, ascending
		From(0).Size(10).        // take documents 0-9
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		// Handle error
		t.Fatal(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var ttyp Tweet
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		t := item.(Tweet)
		fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
	}
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

	// Here's how you iterate through results with full control over each step.
	if searchResult.TotalHits() > 0 {
		fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Tweet
			_ = json.Unmarshal(hit.Source, &t)
			// Work with tweet
			fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
	} else {
		// No hits
		fmt.Print("Found no tweets\n")
	}

	// Update a tweet by the update API of Elasticsearch.
	// We just increment the number of retweets.
	script := elastic.NewScript("ctx._source.retweets += params.num").Param("num", 1)
	update, err := client.Update().Index("twitter").Id("1").
		Script(script).
		Upsert(map[string]interface{}{"retweets": 0}).
		Do(context.Background())
	if err != nil {
		// Handle error
		t.Fatal(err)
	}
	fmt.Printf("New version of tweet %q is now %d", update.Id, update.Version)

	// ...

}

func LF() log.CLoggerFunc {
	return func() log.CLogger {
		return log.L(log.Init(&log.Config{Level: log.TraceLevel}))
	}
}

func Test_ChangingMapping_ModelMappingWithCfgSettings(t *testing.T) {
	es, err := NewEs(&Config{
		Host:     "localhost",
		Port:     "9200",
		Trace:    true,
		Sniff:    true,
		Shards:   2,
		Replicas: 2,
	}, LF())
	if err != nil {
		t.Fatal(err)
	}

	type model struct {
		Field string `json:"field" es:"type:keyword"`
	}

	index := kit.NewRandString()
	err = es.BuildIndexWithModel(index, &model{})
	if err != nil {
		t.Fatal(err)
	}
	type modelNew struct {
		Field  string `json:"field" es:"type:keyword"`
		Field2 string `json:"field2" es:"type:keyword"`
	}

	err = es.BuildIndexWithModel(index, &modelNew{})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ChangingMapping_ModelMapping(t *testing.T) {
	es, err := NewEs(&Config{
		Host:  "localhost",
		Port:  "9200",
		Trace: false,
		Sniff: true,
	}, LF())
	if err != nil {
		t.Fatal(err)
	}

	type model struct {
		Field string `json:"field" es:"type:keyword"`
	}

	index := kit.NewRandString()
	err = es.BuildIndexWithModel(index, &model{})
	if err != nil {
		t.Fatal(err)
	}
	type modelNew struct {
		Field  string `json:"field" es:"type:keyword"`
		Field2 string `json:"field2" es:"type:keyword"`
	}

	err = es.BuildIndexWithModel(index, &modelNew{})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ChangingMapping_ExistentFields(t *testing.T) {
	es, err := NewEs(&Config{
		Host:  "localhost",
		Port:  "9200",
		Trace: false,
		Sniff: true,
	}, LF())
	if err != nil {
		t.Fatal(err)
	}

	type model struct {
		Field string `json:"field" es:"type:keyword"`
	}

	index := kit.NewRandString()
	err = es.BuildIndexWithModel(index, &model{})
	if err != nil {
		t.Fatal(err)
	}
	type modelNew struct {
		Field string `json:"field" es:"type:text"`
	}

	err = es.BuildIndexWithModel(index, &modelNew{})
	assert.NotNil(t, err)
}

func Test_ChangingMapping_ExplicitMapping(t *testing.T) {
	es, err := NewEs(&Config{
		Host:     "localhost",
		Port:     "9200",
		Trace:    true,
		Sniff:    true,
		Shards:   2,
		Replicas: 2,
	}, LF())
	if err != nil {
		t.Fatal(err)
	}

	mapping := `
{
	"mappings": {
		"properties": {
			"field1": {
				"type": "keyword"
			}
		}
	}
}
`
	index := kit.NewRandString()
	err = es.BuildIndexWithExplicitMapping(index, mapping)
	if err != nil {
		t.Fatal(err)
	}
	newMapping := `
{
	"mappings": {
		"properties": {
			"field1": {
				"type": "keyword"
			},
			"field2": {
				"type": "keyword"
			}
		}
	}
}
`
	err = es.BuildIndexWithExplicitMapping(index, newMapping)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Mapping_WhenNotIndexField(t *testing.T) {
	es, err := NewEs(&Config{
		Host:  "localhost",
		Port:  "9200",
		Trace: false,
		Sniff: true,
	}, LF())
	if err != nil {
		t.Fatal(err)
	}

	type model struct {
		Field1 string `json:"field1" es:"type:keyword"`
		Field2 string `json:"field2" es:"type:keyword;-"`
	}

	index := kit.NewRandString()
	err = es.BuildIndexWithModel(index, &model{})
	if err != nil {
		t.Fatal(err)
	}
}
