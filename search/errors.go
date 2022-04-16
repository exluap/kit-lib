package search

import "github.com/exluap/kit/er"

var (
	ErrCodeEsNewClient                     = "ES-001"
	ErrCodeEsIdxExists                     = "ES-002"
	ErrCodeEsIdx                           = "ES-003"
	ErrCodeEsIdxAsync                      = "ES-004"
	ErrCodeEsIdxCreate                     = "ES-007"
	ErrCodeEsBulkIdxAsync                  = "ES-008"
	ErrCodeEsExists                        = "ES-009"
	ErrCodeEsAwaitExistsTimeout            = "ES-010"
	ErrCodeEsInvalidModel                  = "ES-011"
	ErrCodeEsInvalidModelType              = "ES-012"
	ErrCodeEsGetMapping                    = "ES-013"
	ErrCodeEsNoMappingFound                = "ES-014"
	ErrCodeEsMappingSchemaNotExpected      = "ES-015"
	ErrCodeEsMappingExistentFieldsModified = "ES-016"
	ErrCodeEsPutMapping                    = "ES-017"
)

var (
	ErrEsNewClient = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeEsNewClient, "").Err() }
	ErrEsIdxExists = func(cause error, index string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsIdxExists, "").F(er.FF{"idx": index}).Err()
	}
	ErrEsGetMapping = func(cause error, index string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsGetMapping, "").F(er.FF{"idx": index}).Err()
	}
	ErrEsMappingSchemaNotExpected = func(cause error, index string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsMappingSchemaNotExpected, "").F(er.FF{"idx": index}).Err()
	}
	ErrEsPutMapping = func(cause error, index string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsPutMapping, "").F(er.FF{"idx": index}).Err()
	}
	ErrEsIdx = func(cause error, index, id string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsIdx, "").F(er.FF{"idx": index, "id": id}).Err()
	}
	ErrEsIdxAsync = func(cause error, index, id string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsIdxAsync, "").F(er.FF{"idx": index, "id": id}).Err()
	}
	ErrEsBulkIdxAsync = func(cause error, index string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsBulkIdxAsync, "").F(er.FF{"idx": index}).Err()
	}
	ErrEsIdxCreate = func(cause error, index string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsIdxCreate, "").F(er.FF{"idx": index}).Err()
	}
	ErrEsExists = func(cause error, index, id string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsExists, "").F(er.FF{"idx": index, "id": id}).Err()
	}
	ErrEsAwaitExistsTimeout = func(index, id string) error {
		return er.WithBuilder(ErrCodeEsAwaitExistsTimeout, "await timeout").F(er.FF{"idx": index, "id": id}).Err()
	}
	ErrEsNoMappingFound = func(index string) error {
		return er.WithBuilder(ErrCodeEsNoMappingFound, "no mapping found").F(er.FF{"idx": index}).Err()
	}
	ErrEsMappingExistentFieldsModified = func(index string, fields []string) error {
		return er.WithBuilder(ErrCodeEsMappingExistentFieldsModified, "ES doesn't allow changing mapping for existent fields.").F(er.FF{"idx": index, "fields": fields}).Err()
	}
	ErrEsInvalidModel     = func() error { return er.WithBuilder(ErrCodeEsInvalidModel, "invalid model, check tags").Err() }
	ErrEsInvalidModelType = func() error {
		return er.WithBuilder(ErrCodeEsInvalidModelType, "model must be pointer of struct").Err()
	}
)
