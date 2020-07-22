package queries

import (
	"net/url"
	"testing"
)

type PaginationQuerySchema struct {
	PageNum  *uint64 `json:"pageNumber,omitempty"`
	PageSize *uint64 `json:"pageSize,omitempty"`
}

type PaginationQuerySchemaRequired struct {
	PageNum  uint64 `json:"pageNumber"`
	PageSize uint64 `json:"pageSize"`
}

func TestUnmarshal(t *testing.T) {
	t.Run("happy path none required", func(t *testing.T) {
		paginationQuery := PaginationQuerySchema{}

		values := url.Values{
			"pageNumber": []string{"1"},
			"pageSize":   []string{"1000"},
		}

		err := Unmarshal(&paginationQuery, values)
		if err != nil {
			t.Errorf("failed to unmarshal: %+v", err)
		}

		pageSize := paginationQuery.PageSize
		pageNum := paginationQuery.PageNum

		if pageSize == nil || *pageSize != 1000 {
			t.Errorf("pageSize not as expected: %v", pageSize)
		}

		if pageNum == nil || *pageNum != 1 {
			t.Errorf("pageNum not as expected: %v", pageNum)
		}
	})

	t.Run("happy path some required", func(t *testing.T) {
		paginationQuery := PaginationQuerySchemaRequired{}

		values := url.Values{
			"pageNumber": []string{"1"},
			"pageSize":   []string{"1000"},
		}

		err := Unmarshal(&paginationQuery, values)
		if err != nil {
			t.Errorf("failed to unmarshal: %+v", err)
		}

		pageSize := paginationQuery.PageSize
		pageNum := paginationQuery.PageNum

		if pageSize != 1000 {
			t.Errorf("pageSize not as expected: %v", pageSize)
		}

		if pageNum != 1 {
			t.Errorf("pageNum not as expected: %v", pageNum)
		}
	})

}
