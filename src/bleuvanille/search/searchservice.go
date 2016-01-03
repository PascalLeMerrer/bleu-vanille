package search

import (
	"errors"
	
	"bleuvanille/eatable"
	"bleuvanille/log"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzers/keyword_analyzer"
	"github.com/blevesearch/bleve/analysis/language/fr"
)

const (
	INDEX_NAME = "bleve.eatable"
	FIELDNAME_ID = "_id"
)

var index bleve.Index

//Index an eatable in the bleve index
func Index(eatable *eatable.Eatable) error {
	if eatable == nil {
		return errors.New("Cannot index a nil eatable")	
	}
		
	indexLocal, err := getIndex()
	
	if err != nil {
		return err	
	}
	
	if indexLocal == nil {
		return errors.New("nil Index after getting for indexation")	
	}
	
	errIndex := indexLocal.Index(eatable.Id, eatable)
	
	if errIndex != nil {
		return errIndex
	}
	
	return nil
}

//Delete removes an eatable from the index. Use for test
func Delete(eatable *eatable.Eatable) error {
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

//SearchFromQueryString searches in the current index given an full query string
func SearchFromQueryString(querystring string) ([]string, error) {
	indexLocal, errIndex := getIndex()
	
	if errIndex != nil {
		return nil,errIndex	
	}

	q := bleve.NewQueryStringQuery(querystring)
	req := bleve.NewSearchRequest(q)
	//req.Highlight = bleve.NewHighlightWithStyle("ansi")
	req.Fields = []string{FIELDNAME_ID}
	
	res, err := indexLocal.Search(req)
	
	if err != nil {
		log.Fatal(err)
	}
	
	result := make([]string, res.Total)
	
	for indexHit, value := range res.Hits {
		id := value.Fields[FIELDNAME_ID]
		
		if valueid, ok := id.(string);ok {
			result[indexHit] = valueid
		}
	}

	return result, nil
}

//SearchForEatable searches for an eatable in the current index given its name
func SearchForEatable(name string) ([]string, error) {
	qString := `` + name + `~2^2` + " parent.name:" + name + "~2"
	
	return SearchFromQueryString(qString)	
}

//getIndex returns an internal index open once
func getIndex() (bleve.Index, error) {
	if index == nil {
		
		indexReal, errOpenIndex := bleve.Open(INDEX_NAME)
		
		if errOpenIndex != nil || indexReal == nil{
			
			//testcsa
			log.Debug(nil, "Creation of the search index")
			
			mapping := bleve.NewIndexMapping()
//			mapping.DefaultType="eatable"
//			mapping.TypeField="eatable"
			
			eatableMapping := bleve.NewDocumentMapping()
			eatableMapping.Dynamic = false
			eatableMapping.DefaultAnalyzer = fr.AnalyzerName

			//Field Id : only kept to retreive the object from database
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
			
			mapping.DefaultField="name"
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

//			mapping := bleve.NewIndexMapping()
//		    indexRealForCreation, err := bleve.New("example.bleve", mapping)
//		    if err != nil {
//		        return nil, err
//		    }
//		    
//		    index = indexRealForCreation
//		    
//		    return indexRealForCreation, nil
		} else {
			log.Debug(nil, "Openning of the existing search index")
			if indexReal == nil {
				return nil, errors.New("nil Index after openning without error")
			}

			index = indexReal
			return indexReal, nil
		}

	}

	return index, nil
}
