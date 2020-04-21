// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	_ "net/http"
	"testing"
)

func TestScan(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun."}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Type("tweet").Id("1").BodyJson(&tweet1).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("2").BodyJson(&tweet2).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("3").BodyJson(&tweet3).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Flush().Index(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}

	// Match all should return all documents
	cursor, err := client.Scan(testIndexName).
		Size(1).
		// Pretty(true).Debug(true).
		Do()
	if err != nil {
		t.Fatal(err)
	}

	if cursor.Results == nil {
		t.Errorf("expected results != nil; got nil")
	}
	if cursor.Results.Hits == nil {
		t.Errorf("expected results.Hits != nil; got nil")
	}
	if cursor.Results.Hits.TotalHits != 3 {
		t.Errorf("expected results.Hits.TotalHits = %d; got %d", 3, cursor.Results.Hits.TotalHits)
	}
	if len(cursor.Results.Hits.Hits) != 0 {
		t.Errorf("expected len(results.Hits.Hits) = %d; got %d", 0, len(cursor.Results.Hits.Hits))
	}

	pages := 0
	numDocs := 0

	for {
		searchResult, err := cursor.Next()
		if err == EOS {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		pages += 1

		for _, hit := range searchResult.Hits.Hits {
			if hit.Index != testIndexName {
				t.Errorf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
			numDocs += 1
		}
	}

	if pages <= 0 {
		t.Errorf("expected to retrieve at least 1 page; got %d", pages)
	}

	if numDocs != 3 {
		t.Errorf("expected to retrieve %d hits; got %d", 3, numDocs)
	}
}

func TestScanWithQuery(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun."}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Type("tweet").Id("1").BodyJson(&tweet1).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("2").BodyJson(&tweet2).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("3").BodyJson(&tweet3).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Flush().Index(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}

	// Return tweets from olivere only
	termQuery := NewTermQuery("user", "olivere")
	cursor, err := client.Scan(testIndexName).
		Size(1).
		Query(termQuery).
		// Pretty(true).Debug(true).
		Do()
	if err != nil {
		t.Fatal(err)
	}

	if cursor.Results == nil {
		t.Errorf("expected results != nil; got nil")
	}
	if cursor.Results.Hits == nil {
		t.Errorf("expected results.Hits != nil; got nil")
	}
	if cursor.Results.Hits.TotalHits != 2 {
		t.Errorf("expected results.Hits.TotalHits = %d; got %d", 2, cursor.Results.Hits.TotalHits)
	}
	if len(cursor.Results.Hits.Hits) != 0 {
		t.Errorf("expected len(results.Hits.Hits) = %d; got %d", 0, len(cursor.Results.Hits.Hits))
	}

	pages := 0
	numDocs := 0

	for {
		searchResult, err := cursor.Next()
		if err == EOS {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		pages += 1

		for _, hit := range searchResult.Hits.Hits {
			if hit.Index != testIndexName {
				t.Errorf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
			numDocs += 1
		}
	}

	if pages <= 0 {
		t.Errorf("expected to retrieve at least 1 page; got %d", pages)
	}

	if numDocs != 2 {
		t.Errorf("expected to retrieve %d hits; got %d", 2, numDocs)
	}
}
