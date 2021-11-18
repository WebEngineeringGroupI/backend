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

	urls := make([]string, 0, len(records))
	for _, line := range records {
		urls = append(urls, line[0])
	}
	return urls, nil
}

func NewCSV() *CSV {
	return &CSV{}
}
