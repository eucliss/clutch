package store

type Store interface {
	InsertDocument(index string, body map[string]interface{})
	Query(index string, query string) (r map[string]interface{})
	CreateIndices(indicies ...Index)
	DeleteIndex(index string)
	Initialize()
	GetResults(searchResult map[string]interface{}) (res []map[string]interface{})
}
