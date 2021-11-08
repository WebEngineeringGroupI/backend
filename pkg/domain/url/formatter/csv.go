package formatter

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type CSV struct {
}

func (C *CSV) FormatDataToURLs(data []byte) ([]string, error) {
	records, err := csv.NewReader(bytes.NewReader(data)).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", url.ErrUnableToConvertDataToLongURLs, err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("%w: %s", url.ErrUnableToConvertDataToLongURLs, "the list of URLs is empty")
	}

	return C.deduplicateRecords(records), nil
}

func (C *CSV) deduplicateRecords(records [][]string) []string {
	set := map[string]struct{}{}
	for _, record := range records {
		// TODO check what happens if the contents are empty in a record
		set[record[0]] = struct{}{}
	}

	uniqueElements := []string{}
	for element := range set {
		uniqueElements = append(uniqueElements, element)
	}
	return uniqueElements
}

func NewCSV() *CSV {
	return &CSV{}
}
