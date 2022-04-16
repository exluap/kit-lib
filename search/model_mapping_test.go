package search

import (
	"github.com/exluap/kit/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

var logger = log.Init(&log.Config{Level: log.TraceLevel})
var logf = func() log.CLogger {
	return log.L(logger)
}

func Test_ModelToMapping_WhenNil(t *testing.T) {
	s := &esImpl{logger: logf}
	mapping, _ := s.modelToMapping(nil)
	assert.Empty(t, mapping)
}

func Test_ModelToMapping_WhenNoJsonMapping(t *testing.T) {
	s := &esImpl{logger: logf}
	model := &struct {
		Field string `es:"type:keyword"`
	}{
		Field: "value",
	}
	_, err := s.modelToMapping(model)
	assert.Error(t, err)
}

func Test_ModelToMapping_WhenJsonEmptyField(t *testing.T) {
	s := &esImpl{logger: logf}
	model := &struct {
		Field string `json:"" es:"type:keyword"`
	}{
		Field: "value",
	}
	_, err := s.modelToMapping(model)
	assert.Error(t, err)
}

func Test_ModelToMapping_WhenNotPointerType(t *testing.T) {
	s := &esImpl{logger: logf}
	model := struct {
		Field string `json:"" es:"type:keyword"`
	}{
		Field: "value",
	}
	_, err := s.modelToMapping(model)
	assert.Error(t, err)
}

func Test_ModelToMapping_WhenNotStructType(t *testing.T) {
	s := &esImpl{logger: logf}
	model := "string"
	_, err := s.modelToMapping(model)
	assert.Error(t, err)
}

func Test_ModelToMapping_WhenWithType(t *testing.T) {
	s := &esImpl{logger: logf}
	model := &struct {
		Field string `json:"field" es:"type:keyword"`
	}{
		Field: "value",
	}
	mapping, err := s.modelToMapping(model)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, mapping)
	assert.Equal(t, 1, len(mapping.Mappings.Properties))
	f, ok := mapping.Mappings.Properties["field"]
	assert.True(t, ok)
	assert.Equal(t, "keyword", f.Type)
}

func Test_ModelToMapping_WhenNoEsTag_Skipped(t *testing.T) {
	s := &esImpl{logger: logf}
	model := &struct {
		Field1 string `json:"field1"`
		Field2 string `json:"field2" es:"type:text"`
	}{
		Field1: "value",
		Field2: "value",
	}
	mapping, err := s.modelToMapping(model)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, mapping)
	assert.Equal(t, 1, len(mapping.Mappings.Properties))
}

func Test_ModelToMapping_WhenEsTagEmpty_Skipped(t *testing.T) {
	s := &esImpl{logger: logf}
	model := &struct {
		Field1 string `json:"field1" es:""`
		Field2 string `json:"field2" es:"type:text"`
	}{
		Field1: "value",
		Field2: "value",
	}
	mapping, err := s.modelToMapping(model)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, mapping)
	assert.Equal(t, 1, len(mapping.Mappings.Properties))
}

func Test_ModelToMapping_WhenNoOneEsTag(t *testing.T) {
	s := &esImpl{logger: logf}
	model := &struct {
		Field string `json:"field"`
	}{
		Field: "value",
	}
	_, err := s.modelToMapping(model)
	assert.Error(t, err)
}

func Test_ModelToMapping_WhenTypeWrong(t *testing.T) {
	s := &esImpl{logger: logf}
	model := &struct {
		Field string `json:"field" es:"type:wrong"`
	}{
		Field: "value",
	}
	_, err := s.modelToMapping(model)
	assert.Error(t, err)
}

func Test_ModelToMapping_WhenTagWrong(t *testing.T) {
	s := &esImpl{logger: logf}
	model := &struct {
		Field string `json:"field es:wrong"`
	}{
		Field: "value",
	}
	_, err := s.modelToMapping(model)
	assert.Error(t, err)
}

func Test_ModelToMapping_WithTwoFields(t *testing.T) {
	s := &esImpl{logger: logf}
	model := &struct {
		Field1 string `json:"field1" es:"type:keyword"`
		Field2 string `json:"field2" es:"type:text"`
	}{
		Field1: "value",
		Field2: "value",
	}
	mapping, err := s.modelToMapping(model)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, mapping)
	assert.Equal(t, 2, len(mapping.Mappings.Properties))
	f, ok := mapping.Mappings.Properties["field1"]
	assert.True(t, ok)
	assert.Equal(t, "keyword", f.Type)
	f, ok = mapping.Mappings.Properties["field2"]
	assert.True(t, ok)
	assert.Equal(t, "text", f.Type)
}

func Test_ModelToMapping_WithNoIndex(t *testing.T) {
	s := &esImpl{logger: logf}
	model := &struct {
		Field1 string `json:"field1" es:"type:keyword"`
		Field2 string `json:"field2" es:"-"`
		Field3 string `json:"field3" es:"type:text;-"`
	}{
		Field1: "value",
		Field2: "value",
		Field3: "value",
	}
	mapping, err := s.modelToMapping(model)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, mapping)
	assert.Equal(t, 3, len(mapping.Mappings.Properties))
	f, ok := mapping.Mappings.Properties["field1"]
	assert.True(t, ok)
	assert.Equal(t, "keyword", f.Type)
	f, ok = mapping.Mappings.Properties["field2"]
	assert.True(t, ok)
	assert.Empty(t, f.Type)
	assert.False(t, *f.Index)
	f, ok = mapping.Mappings.Properties["field3"]
	assert.True(t, ok)
	assert.Equal(t, "text", f.Type)
	assert.False(t, *f.Index)
}
