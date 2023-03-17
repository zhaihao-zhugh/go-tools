package es

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var es *elasticsearch.Client

func NewConnect(host []string) error {
	cfg := elasticsearch.Config{
		Addresses: host,
	}
	cli, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return err
	}
	es = cli
	return nil
}

func HandleESDefine(index string, body io.Reader) (result []byte, err error) {
	req := esapi.IndicesCreateRequest{
		Index: index, // Index name
		Body:  body,  // Document body
	}
	res, e := req.Do(context.Background(), es)
	if e != nil {
		err = e
		return
	}
	result, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return
}

func HandleESCreate(index string, body io.Reader, id string) (result []byte, err error) {
	res, e := es.Index(
		index,                        // Index name
		body,                         // Document body
		es.Index.WithDocumentID(id),  // Document ID
		es.Index.WithRefresh("true"), // Refresh
	)
	if e != nil {
		err = e
		return
	}
	result, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return
}

func HandleESSearch(index string, body io.Reader) (result []byte, err error) {
	res, e := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithBody(body),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if e != nil {
		err = e
		return
	}
	result, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return
}

func HandleESUpdate(index string, doc_id string, body io.Reader) (result []byte, err error) {
	res, e := es.Update(
		index,
		doc_id,
		body,
		es.Update.WithRefresh(`true`),
		es.Update.WithPretty(),
	)
	if e != nil {
		err = e
		return
	}
	if res.StatusCode != 200 {
		err = errors.New("es update fail")
		return
	}
	result, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return
}

func HandleESUpdateByQuery(indexes []string, body io.Reader) (result []byte, err error) {
	res, e := es.UpdateByQuery(
		indexes,
		es.UpdateByQuery.WithBody(body),
		es.UpdateByQuery.WithRefresh(true),
	)
	if e != nil {
		err = e
		return
	}
	if res.StatusCode != 200 {
		err = errors.New("es update fail")
		return
	}
	result, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return
}

func HandleESDeleteById(index string, doc_id string) (result []byte, err error) {
	res, e := es.Delete(
		index,
		doc_id,
	)
	if e != nil {
		err = e
		return
	}
	if res.StatusCode != 200 {
		err = errors.New("es delete fail")
		return
	}
	result, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return
}

func HandleESDeleteByQuery(indexes []string, body io.Reader) (result []byte, err error) {
	res, e := es.DeleteByQuery(
		indexes,
		body,
	)
	if e != nil {
		err = e
		return
	}
	if res.StatusCode != 200 {
		err = errors.New("es delete fail")
		return
	}
	result, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return
}

func HandleESGet(index string, doc_id string) (result []byte, err error) {
	res, e := es.Get(
		index,
		doc_id,
		es.Get.WithPretty(),
	)
	if e != nil {
		err = e
		return
	}
	if res.StatusCode != 200 {
		err = errors.New("no data")
		return
	}
	result, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return
}

func HandleESState(index string, field string) (result []byte, err error) {
	res, e := es.Indices.Stats(
		es.Indices.Stats.WithIndex(index),
		es.Indices.Stats.WithMetric(field),
	)
	if e != nil {
		err = e
		return
	}
	if res.StatusCode != 200 {
		err = errors.New("status code error")
		return
	}
	result, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return
}

func HandleESCount(index string, body io.Reader) int {
	var res *esapi.Response
	var err error
	if body != nil {
		res, err = es.Count(
			es.Count.WithContext(context.Background()),
			es.Count.WithIndex(index),
			es.Count.WithBody(body),
			es.Count.WithPretty(),
		)
	} else {
		res, err = es.Count(
			es.Count.WithContext(context.Background()),
			es.Count.WithIndex(index),
			es.Count.WithPretty(),
		)
	}
	if err != nil {
		return 0
	}
	res_body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return 0
	}
	result := make(map[string]interface{})
	json.Unmarshal(res_body, &result)

	if count, ok := result["count"].(float64); ok {
		return int(count)
	}
	return 0
}

func IsHaveValue(index string, field string, value string) bool {
	req := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				field: value,
			},
		},
	}
	req_json, _ := json.Marshal(req)
	body := bytes.NewBuffer(req_json)
	count := HandleESCount(index, body)
	if count > 0 {
		return true
	}
	return false
}
