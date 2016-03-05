package search

import (
	"errors"
	"fmt"

	"bleuvanille/eatable"
	"bleuvanille/log"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzers/keyword_analyzer"
	"github.com/blevesearch/bleve/analysis/language/fr"
)

const (
	INDEX_NAME   = "bleveindex.eatable"
	FIELDNAME_ID = "_id"
)

var index bleve.Index

type SearchService struct {
}

//Index an eatable in the bleve index
func (search *SearchService) Index(eatableVal *eatable.Eatable) error {
	if eatableVal == nil {
		return errors.New("Cannot index a nil eatable")
	}

	indexLocal, err := getIndex()

	if err != nil {
		return err
	}

	if indexLocal == nil {
		return errors.New("nil Index after getting for indexation")
	}

	//Load the parent
	parent, errParent := eatable.GetParent(eatableVal)

	if errParent != nil {
		log.Error(nil, "Error while fetching the parent for "+eatableVal.Id+" : "+errParent.Error())
	} else {
		eatableVal.Parent = parent
	}

	errIndex := indexLocal.Index(eatableVal.Id, eatableVal)

	if errIndex != nil {
		return errIndex
	}

	return nil
}

//indexAll reset the index and indexes every Eatable from the database
func (search *SearchService) indexAll() (int, error) {
	indexLocal, err := getIndex()

	if err != nil {
		return 0, err
	}

	if indexLocal == nil {
		return 0, errors.New("nil Index after getting for indexation")
	}

	//Reset the index
	eatablesids, _, _ := search.SearchForAllEatable(0, 0)

	for _, id := range eatablesids {
		errDel := indexLocal.Delete(id)

		if errDel != nil {
			log.Debug(nil, "Error while delete from index  "+id+" : "+errDel.Error())
		}
	}

	eatables, _ := eatable.FindAll("", "")

	for _, eatable := range eatables {
		errIndex := search.Index(&eatable)

		if errIndex != nil {
			log.Debug(nil, "Error while index "+eatable.Id+" : "+errIndex.Error())
		}
	}

	return len(eatables), nil
}

//Delete removes an eatable from the index. Use for test
func (search *SearchService) Delete(eatable *eatable.Eatable) error {
	if eatable == nil {
		return errors.New("Cannot delete a nil eatable")
	}

	indexLocal, err := getIndex()

	if err != nil {
		return err
	}

	if indexLocal == nil {
		return errors.New("nil Index for deletion")
	}

	errDelete := indexLocal.Delete(eatable.Id)

	if errDelete != nil {
		return errDelete
	}

	return nil
}

//DeleteFromKey removes a eatable of the index from its id
func (search *SearchService) DeleteFromId(id string) error {
	indexLocal, err := getIndex()

	if err != nil {
		return err
	}

	if indexLocal == nil {
		return errors.New("nil Index for deletion")
	}

	errDelete := indexLocal.Delete(id)

	if errDelete != nil {
		return errDelete
	}

	return nil
}

//SearchFromQueryString searches in the current index given a full query string
func (searchVal *SearchService) SearchFromQueryString(querystring string, offset int, limit int) ([]string, int, error) {
	q := bleve.NewQueryStringQuery(querystring)
	return search(q, offset, limit)
}

//SearchForEatable searches for an eatable in the current index given its name
func (searchVal *SearchService) SearchForEatable(name string, offset int, limit int) ([]string, int, error) {
	qString := name + "^4 " + name + "~2^2" + " parent.name:" + name + "~2"
	q := bleve.NewQueryStringQuery(qString)

	return search(q, offset, limit)
}

//SearchForAllEatable returns all eatable contains in the index
func (searchVal *SearchService) SearchForAllEatable(offset int, limit int) ([]string, int, error) {
	q := bleve.NewMatchAllQuery()
	return search(q, offset, limit)
}

func (searchVal *SearchService) SearchPrefix(name string, offset int, limit int) ([]string, int, error) {
	qString := `` + name
	q := bleve.NewPrefixQuery(qString)
	//q.FieldVal = "FieldVal"
	return search(q, offset, limit)
}

//getIndex returns an internal index open once
func getIndex() (bleve.Index, error) {
	if index == nil {

		indexReal, errOpenIndex := bleve.Open(INDEX_NAME)

		if errOpenIndex != nil || indexReal == nil {
			return createIndex()
		} else {
			if indexReal == nil {
				return nil, errors.New("nil Index after openning without error")
			}

			count, errCount := indexReal.DocCount()

			if errCount != nil {
				return nil, errors.New("Error while openning and counting the index")
			}
			logMessage := fmt.Sprintf("Openning of the existing search index : %d documents", count)

			log.Debug(nil, logMessage)

			index = indexReal
			return indexReal, nil
		}

	}

	return index, nil
}

//createIndex define the search index, create the associated directory structure and open it
func createIndex() (bleve.Index, error) {
	log.Debug(nil, "Creation of the search index")

	mapping := bleve.NewIndexMapping()

	eatableMapping := bleve.NewDocumentMapping()
	eatableMapping.Dynamic = false
	eatableMapping.DefaultAnalyzer = fr.AnalyzerName

	//Field Id : only kept to retrieve the object from database
	idFieldMapping := bleve.NewTextFieldMapping()
	idFieldMapping.Store = true
	idFieldMapping.Index = false
	idFieldMapping.Analyzer = keyword_analyzer.Name
	eatableMapping.AddFieldMappingsAt(FIELDNAME_ID, idFieldMapping)

	//Field name : name of the eatable. It is the main text source.
	nameFieldMapping := bleve.NewTextFieldMapping()
	nameFieldMapping.Index = true
	nameFieldMapping.Store = false
	nameFieldMapping.IncludeInAll = false
	nameFieldMapping.Name = "name"
	eatableMapping.AddFieldMappingsAt("name", nameFieldMapping)

	descriptionFieldMapping := bleve.NewTextFieldMapping()
	descriptionFieldMapping.Index = true
	descriptionFieldMapping.Store = false
	descriptionFieldMapping.IncludeInAll = false
	eatableMapping.AddFieldMappingsAt("description", descriptionFieldMapping)

	//Field name : name of the parents of the eatable. To be able to find the eatable through its category
	parent := bleve.NewDocumentMapping()
	parentNameFieldMapping := bleve.NewTextFieldMapping()
	parentNameFieldMapping.Index = true
	parentNameFieldMapping.Store = false
	parentNameFieldMapping.IncludeInAll = false
	parent.AddFieldMappingsAt("name", parentNameFieldMapping)
	eatableMapping.AddSubDocumentMapping("parent", parent)

	mapping.DefaultField = "name"
	mapping.AddDocumentMapping("Eatable", eatableMapping)
	mapping.DefaultMapping = eatableMapping
	indexRealForCreation, errForCreation := bleve.New(INDEX_NAME, mapping)

	if errForCreation != nil {
		return nil, errForCreation
	}

	if indexRealForCreation == nil {
		return nil, errors.New("nil Index after creation without error")
	}

	index = indexRealForCreation
	return indexRealForCreation, nil
}

func search(q bleve.Query, offset int, limit int) ([]string, int, error) {
	indexLocal, errIndex := getIndex()

	if errIndex != nil {
		return nil, 0, errIndex
	}

	req := bleve.NewSearchRequest(q)
	req.Fields = []string{FIELDNAME_ID}
	req.From = offset

	if limit > 0 {
		req.Size = limit
	}

	res, err := indexLocal.Search(req)

	if err != nil {
		log.Fatal(err)
		return nil, 0, err
	}

	returnedCount := int(res.Total)

	if returnedCount > req.Size {
		returnedCount = req.Size
	}

	result := make([]string, returnedCount)

	for indexHit, value := range res.Hits {
		id := value.Fields[FIELDNAME_ID]

		if valueid, ok := id.(string); ok {
			result[indexHit] = valueid
		}
	}

	return result, int(res.Total), nil
}
