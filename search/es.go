package search

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/exluap/kit"
	"github.com/exluap/kit/log"
	"github.com/olivere/elastic/v7"
	"reflect"
	"strings"
	"time"
)

// Config - model of ES configuration
type Config struct {
	Host     string // Host - ES host
	Port     string // Port - ES port
	Trace    bool   // Trace enables tracing mode
	Sniff    bool   // Sniff - read https://github.com/olivere/elastic/issues/387
	Shards   int    // Shards - how many shards to be created for index
	Replicas int    // Replicas - how many replicas to eb created for index
}

// Search allows indexing and searching with ES
type Search interface {
	// BuildIndexWithExplicitMapping takes index name and mapping as string(json)
	// see ES doc about https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping.html
	// it checks if index exists, if not creates it
	BuildIndexWithExplicitMapping(index string, mapping string) error
	// BuildIndexWithModel takes model with tags and builds index based on provided info
	// if index doesn't exist, a new index is created
	// if index exists it checks whether existent mapping is modified and if it is, it fails. If only new fields added, it handles them as PUT
	// Note, "json" tag must be specified together with "es" tag
	// example:
	// type IndexModel struct {
	//   Field1 string `json:"field1" es:"type:text"`   	// field is mapped with text type
	//   Field2 string `json:"field2" es:"type:keyword"` 	// field is mapped with keyword type
	//   Field3 time.Time `json:"field3" es:"type:date"` 	// field is mapped with date type
	//   Field4 time.Time `json:"field4" es:"-"` 			// field is mapped with "Index=false"
	// }
	//
	// model must be pointer type
	BuildIndexWithModel(index string, model interface{}) error
	// Index indexes a document
	Index(index string, id string, data interface{}) error
	// IndexAsync indexes a document async
	IndexAsync(index string, id string, data interface{})
	// IndexBulkAsync allows indexing bulk of documents in one hit
	IndexBulkAsync(index string, docs map[string]interface{})
	// GetClient provides an access to ES client
	GetClient() *elastic.Client
	// Close closes client
	Close()
	// Exists checks if a document exists in the index
	Exists(index, id string) (bool, error)
	// AwaitDocExists periodically hits index for a document with a given id within a timeout
	// if it results in nil, then document exists, otherwise it's either ES error or timeout
	AwaitDocExists(index, id string, timeout time.Duration) <-chan error
}

type esImpl struct {
	client *elastic.Client
	logger log.CLoggerFunc
	cfg    *Config
}

func (s *esImpl) l() log.CLogger {
	return s.logger().Cmp("es")
}

func NewEs(cfg *Config, logger log.CLoggerFunc) (Search, error) {

	s := &esImpl{
		logger: logger,
		cfg:    cfg,
	}
	l := s.l().Mth("new").F(log.FF{"host": cfg.Host, "sniff": cfg.Sniff})

	url := fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port)

	opts := []elastic.ClientOptionFunc{elastic.SetURL(url), elastic.SetSniff(cfg.Sniff)}
	if cfg.Trace {
		opts = append(opts, elastic.SetTraceLog(s.l().Mth("es-trace")))
	}

	cl, err := elastic.NewClient(opts...)
	if err != nil {
		return nil, ErrEsNewClient(err)
	}
	s.client = cl
	l.Inf("ok")
	return s, nil
}

const (
	T_KEYWORD            = "keyword"
	T_TEXT               = "text"
	T_DATE               = "date"
	T_BOOL               = "boolean"
	T_SEARCH_AS_YOU_TYPE = "search_as_you_type"
)

var typesMap = map[string]struct{}{
	T_KEYWORD:            {},
	T_TEXT:               {},
	T_DATE:               {},
	T_BOOL:               {},
	T_SEARCH_AS_YOU_TYPE: {},
}

type EsProperty struct {
	Type  string `json:"type,omitempty"`  // Type specifies a datatype
	Index *bool  `json:"index,omitempty"` // Index - if false, field isn't indexed
}

type EsProperties map[string]*EsProperty

type EsSettings struct {
	NumberOfShards   int `json:"number_of_shards"`
	NumberOfReplicas int `json:"number_of_replicas"`
}

type EsMapping struct {
	Settings EsSettings `json:"settings"`
	Mappings struct {
		Properties EsProperties `json:"properties"`
	} `json:"mappings"`
}

func (s *esImpl) setSettings(mapping *EsMapping) {
	if mapping.Settings.NumberOfReplicas == 0 {
		mapping.Settings.NumberOfReplicas = s.cfg.Replicas
		if mapping.Settings.NumberOfReplicas == 0 {
			mapping.Settings.NumberOfReplicas = 1
		}
	}
	if mapping.Settings.NumberOfShards == 0 {
		mapping.Settings.NumberOfShards = s.cfg.Shards
		if mapping.Settings.NumberOfShards == 0 {
			mapping.Settings.NumberOfShards = 1
		}
	}
}

// modelToMapping creates ES mapping based on model tag
// check model_mapping_test for usage details
func (s *esImpl) modelToMapping(modelObj interface{}) (*EsMapping, error) {

	s.l().Mth("model-to-mapping")

	if modelObj == nil {
		return nil, nil
	}

	type params map[string]string

	if reflect.ValueOf(modelObj).Kind() != reflect.Ptr || reflect.TypeOf(modelObj).Elem().Kind() != reflect.Struct {
		return nil, ErrEsInvalidModelType()
	}

	// takes type description
	r := reflect.TypeOf(modelObj).Elem()
	mappingProperties := make(EsProperties)

	// build mapping fields map
	// go through fields
	for i := 0; i < r.NumField(); i++ {
		field := r.Field(i)

		// check json tag
		// we use index field name from json mapping
		// if json tag missing, field is skipped
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {

			jsonParams := strings.Split(jsonTag, ",")
			// if there is no field name in json tag
			if len(jsonParams) == 0 {
				return nil, ErrEsInvalidModel()
			}
			indexFieldName := jsonParams[0]

			// take es tag
			esTag := field.Tag.Get("es")
			// if es tag missing, skip the field
			if esTag != "" {
				esTagParams := make(params)
				// take params separated by ;
				params := strings.Split(esTag, ";")
				for _, p := range params {
					kv := strings.Split(p, ":")
					if len(kv) == 2 {
						esTagParams[kv[0]] = kv[1]
					} else {
						esTagParams[kv[0]] = ""
					}
				}
				// populate mapping params
				mappingProperties[indexFieldName] = &EsProperty{}
				for esTag, esTagValue := range esTagParams {
					switch esTag {
					case "type":
						if _, ok := typesMap[esTagValue]; !ok {
							return nil, ErrEsInvalidModel()
						} else {
							mappingProperties[indexFieldName].Type = esTagValue
						}
					case "-":
						f := false
						mappingProperties[indexFieldName].Index = &f
					}
				}
			}
		}
	}

	// return ES mapping if specified
	if len(mappingProperties) == 0 {
		return nil, ErrEsInvalidModel()
	} else {
		r := &EsMapping{}
		r.Mappings.Properties = mappingProperties
		return r, nil
	}

}

func (s *esImpl) BuildIndexWithExplicitMapping(index string, mapping string) error {

	l := s.l().Mth("create-idx-mapping")

	exists, err := s.client.IndexExists(index).Do(context.Background())
	if err != nil {
		return ErrEsIdxExists(err, index)
	}

	if exists {
		// we allow adding new fields to mapping, but don't allow changing existent ones
		l.DbgF("index %s exists", index)
		// get current mapping
		curMappings, err := s.client.GetMapping().Index(index).Do(context.Background())
		if err != nil {
			return ErrEsGetMapping(err, index)
		}
		curMapping, ok := curMappings[index]
		if !ok {
			return ErrEsNoMappingFound(index)
		}

		mappingJson, _ := json.Marshal(curMapping)
		currentMapping := &EsMapping{}
		err = json.Unmarshal(mappingJson, currentMapping)
		if err != nil {
			return ErrEsMappingSchemaNotExpected(err, index)
		}

		// new mapping
		newMapping := &EsMapping{}
		err = json.Unmarshal([]byte(mapping), newMapping)
		if err != nil {
			return ErrEsMappingSchemaNotExpected(err, index)
		}
		// check if there are changes in existent fields
		if v := s.checkExistentFieldsMappingModified(currentMapping, newMapping); len(v) > 0 {
			return ErrEsMappingExistentFieldsModified(index, v)
		}
		// extract added fields
		if addedFieldsMapping := s.addedFieldsMapping(currentMapping, newMapping); len(addedFieldsMapping.Mappings.Properties) > 0 {
			addedFieldsMappingJson, _ := json.Marshal(addedFieldsMapping.Mappings)
			_, err = s.client.PutMapping().Index(index).BodyString(string(addedFieldsMappingJson)).Do(context.Background())
			if err != nil {
				return ErrEsPutMapping(err, index)
			}
			l.DbgF("fields added: %+v", addedFieldsMapping.Mappings.Properties)
		}
		return nil
	} else {
		// new mapping
		newMapping := &EsMapping{}
		err = json.Unmarshal([]byte(mapping), newMapping)
		if err != nil {
			return ErrEsMappingSchemaNotExpected(err, index)
		}
		s.setSettings(newMapping)
		newMappingJson, _ := json.Marshal(newMapping)
		_, err = s.client.CreateIndex(index).BodyString(string(newMappingJson)).Do(context.Background())
		if err != nil {
			return ErrEsIdxCreate(err, index)
		}
		return nil
	}
}

func (s *esImpl) addedFieldsMapping(currentMapping, newMapping *EsMapping) *EsMapping {
	addedFieldsMapping := &EsMapping{}
	addedFieldsMapping.Mappings.Properties = make(EsProperties)
	for f, v := range newMapping.Mappings.Properties {
		if _, found := currentMapping.Mappings.Properties[f]; !found {
			addedFieldsMapping.Mappings.Properties[f] = v
		}
	}
	return addedFieldsMapping
}

// checkExistentFieldsMappingModified compares current and provided mapping and returns true if there are changes in existent fields
func (s *esImpl) checkExistentFieldsMappingModified(currentMapping, newMapping *EsMapping) []string {
	var modifiedFields []string
	for curFieldName, curField := range currentMapping.Mappings.Properties {
		for newFieldName, newField := range newMapping.Mappings.Properties {
			if curFieldName == newFieldName && curField.Type != newField.Type {
				modifiedFields = append(modifiedFields, curFieldName)
			}
		}
	}
	return modifiedFields
}

func (s *esImpl) BuildIndexWithModel(index string, model interface{}) error {

	l := s.l().Mth("create-idx-mapping")

	// check if index exists
	exists, err := s.client.IndexExists(index).Do(context.Background())
	if err != nil {
		return ErrEsIdxExists(err, index)
	}

	if exists {
		// we allow adding new fields to mapping, but don't allow changing existent ones

		l.DbgF("index %s exists", index)

		// get current mapping
		mappings, err := s.client.GetMapping().Index(index).Do(context.Background())
		if err != nil {
			return ErrEsGetMapping(err, index)
		}
		mapping, ok := mappings[index]
		if !ok {
			return ErrEsNoMappingFound(index)
		}

		mappingJson, _ := json.Marshal(mapping)
		currentMapping := &EsMapping{}
		err = json.Unmarshal(mappingJson, currentMapping)
		if err != nil {
			return ErrEsMappingSchemaNotExpected(err, index)
		}
		// new mapping
		newMapping, err := s.modelToMapping(model)
		if err != nil {
			return err
		}
		// check if there are changes in existent fields
		if v := s.checkExistentFieldsMappingModified(currentMapping, newMapping); len(v) > 0 {
			return ErrEsMappingExistentFieldsModified(index, v)
		}
		// extract added fields
		if addedFieldsMapping := s.addedFieldsMapping(currentMapping, newMapping); len(addedFieldsMapping.Mappings.Properties) > 0 {
			addedFieldsMappingJson, _ := json.Marshal(addedFieldsMapping.Mappings)
			_, err = s.client.PutMapping().Index(index).BodyString(string(addedFieldsMappingJson)).Do(context.Background())
			if err != nil {
				return ErrEsPutMapping(err, index)
			}
			l.DbgF("fields added: %+v", addedFieldsMapping.Mappings.Properties)
		}
		return nil
	} else {
		// transform model to mapping
		mapping, err := s.modelToMapping(model)
		if err != nil {
			return err
		}
		// set settings
		s.setSettings(mapping)
		// build mapping json
		mappingJson, _ := json.Marshal(mapping)
		// create a new index
		_, err = s.client.CreateIndex(index).BodyString(string(mappingJson)).Do(context.Background())
		if err != nil {
			return ErrEsIdxCreate(err, index)
		}
		return nil
	}
}

func (s *esImpl) Index(index string, id string, doc interface{}) error {

	s.l().Mth("indexation").F(log.FF{"index": index, "id": id}).Dbg().Trc(kit.Json(doc))

	svc := s.client.Index().
		Index(index).
		Id(id).
		BodyJson(doc)

	_, err := svc.Do(context.Background())
	if err != nil {
		return ErrEsIdx(err, index, id)
	}
	return nil
}

func (s *esImpl) IndexAsync(index string, id string, doc interface{}) {

	go func() {

		l := s.l().Mth("indexation").F(log.FF{"index": index, "id": id}).Dbg().Trc(kit.Json(doc))

		svc := s.client.Index().
			Index(index).
			Id(id).
			BodyJson(doc)

		_, err := svc.Do(context.Background())
		if err != nil {
			l.E(ErrEsIdxAsync(err, index, id)).Err()
		}

	}()

}

func (s *esImpl) IndexBulkAsync(index string, docs map[string]interface{}) {
	go func() {
		l := s.l().Mth("bulk-indexation").F(log.FF{"index": index, "docs": len(docs)}).Dbg()

		bulk := s.client.Bulk().Index(index)
		for id, doc := range docs {
			bulk.Add(elastic.NewBulkIndexRequest().Id(id).Doc(doc))
		}

		_, err := bulk.Do(context.Background())
		if err != nil {
			l.E(ErrEsBulkIdxAsync(err, index)).Err()
		}
	}()
}

// Exists checks if a document exists in the index
func (s *esImpl) Exists(index, id string) (bool, error) {
	l := s.l().Mth("exists").F(log.FF{"index": index, "id": id})
	res, err := s.client.Exists().Index(index).Id(id).Do(context.Background())
	if err != nil {
		return false, ErrEsExists(err, index, id)
	}
	l.DbgF("res: %v", res)
	return res, nil
}

func (s *esImpl) AwaitDocExists(index, id string, timeout time.Duration) <-chan error {
	resChan := make(chan error)
	go func() {
		c, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		ticker := time.NewTicker(time.Millisecond * 100)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				exists, err := s.Exists(index, id)
				if err != nil {
					resChan <- err
					return
				}
				if exists {
					resChan <- nil
					return
				}
			case <-c.Done():
				resChan <- ErrEsAwaitExistsTimeout(index, id)
				return
			}
		}
	}()
	return resChan
}

func (s *esImpl) GetClient() *elastic.Client {
	return s.client
}

func (s *esImpl) Close() {
	s.client.Stop()
}
