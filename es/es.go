package es

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ESCLIENT struct {
	*elasticsearch.Client
}

func NewConnect(host string) *ESCLIENT {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://" + host,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return &ESCLIENT{es}
}

func (es *ESCLIENT) HandleESDefine(index string, body *bytes.Buffer) error {
	req := esapi.IndicesCreateRequest{
		Index: index, // Index name
		Body:  body,  // Document body
	}

	res, err := req.Do(context.Background(), es)
	result, err := ioutil.ReadAll(res.Body)
	log.Printf("HandleESDefine result: %s\n", result)
	defer res.Body.Close()
	return err
}

func (es *ESCLIENT) HandleESCreate(index string, body *bytes.Buffer, id string) error {
	res, err := es.Index(
		index,                        // Index name
		body,                         // Document body
		es.Index.WithDocumentID(id),  // Document ID
		es.Index.WithRefresh("true"), // Refresh
	)
	result, err := ioutil.ReadAll(res.Body)
	log.Printf("es create result: %s\n", result)
	defer res.Body.Close()
	return err
}

func (es *ESCLIENT) HandleESSearch(index string, body *bytes.Buffer) ([]byte, error) {
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithBody(body),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	defer res.Body.Close()
	if err != nil {
		return []byte{}, err
	}
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return result, nil
}

func (es *ESCLIENT) HandleESUpdate(index string, doc_id string, body *bytes.Buffer) (bool, error) {
	res, err := es.Update(
		index,
		doc_id,
		body,
		es.Update.WithRefresh(`true`),
		es.Update.WithPretty(),
	)
	defer res.Body.Close()
	if err != nil {
		return false, err
	}
	if res.StatusCode == 200 {
		return true, err
	}
	return false, errors.New("es update fail")
}

func (es *ESCLIENT) HandleESUpdateByQuery(indexes []string, body *bytes.Buffer) (bool, error) {
	res, err := es.UpdateByQuery(
		indexes,
		es.UpdateByQuery.WithBody(body),
		es.UpdateByQuery.WithRefresh(true),
	)
	defer res.Body.Close()
	if err != nil {
		return false, err
	}
	if res.StatusCode == 200 {
		return true, err
	}
	return false, errors.New("es update fail")
}

func (es *ESCLIENT) HandleESDeleteById(index string, doc_id string) (bool, error) {
	res, err := es.Delete(
		index,
		doc_id,
	)
	defer res.Body.Close()
	if err != nil {
		return false, err
	}
	if res.StatusCode == 200 {
		return true, err
	}
	return false, errors.New("es delete fail")
}

func (es *ESCLIENT) HandleESDeleteByQuery(indexes []string, body *bytes.Buffer) (bool, error) {
	res, err := es.DeleteByQuery(
		indexes,
		body,
	)
	defer res.Body.Close()
	if err != nil {
		return false, err
	}
	if res.StatusCode == 200 {
		return true, err
	}
	return false, errors.New("es delete fail")
}

func (es *ESCLIENT) HandleESGet(index string, doc_id string) ([]byte, error) {
	res, err := es.Get(
		index,
		doc_id,
		es.Get.WithPretty(),
	)
	defer res.Body.Close()
	if err != nil {
		return []byte{}, err
	}
	result, err := ioutil.ReadAll(res.Body)
	if err == nil && res.StatusCode == 200 {
		return result, nil
	}
	return []byte{}, errors.New("no data")
}

func (es *ESCLIENT) HandleESCount(index string, body *bytes.Buffer) int {
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

	defer res.Body.Close()
	if err != nil {
		return 0
	}
	res_body, err := ioutil.ReadAll(res.Body)
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

func (es *ESCLIENT) HandleESState(index string, field string) (*[]byte, error) {
	res, err := es.Indices.Stats(
		es.Indices.Stats.WithIndex(index),
		es.Indices.Stats.WithMetric(field),
	)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 200 {
		result, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
	return nil, errors.New("status code error")
}

func (es *ESCLIENT) IsHaveValue(index string, field string, value string) bool {
	req := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				field: value,
			},
		},
	}

	req_json, _ := json.Marshal(req)
	body := bytes.NewBuffer(req_json)
	count := es.HandleESCount(index, body)
	if count > 0 {
		return true
	}

	return false
}
