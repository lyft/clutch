package topology

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
)

func TestMaxQueryLimit(t *testing.T) {
	testCases := []struct {
		id          string
		input       uint64
		shouldError bool
	}{
		{
			id:          "Under limit",
			input:       999,
			shouldError: false,
		},
		{
			id:          "Equal to limit",
			input:       1000,
			shouldError: false,
		},
		{
			id:          "Above limit",
			input:       1001,
			shouldError: true,
		},
	}

	for _, test := range testCases {
		_, _, err := paginatedQueryBuilder(&topologyv1.SearchRequest_Filter{}, &topologyv1.SearchRequest_Sort{}, "0", test.input)
		if test.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestPaginatedQueryBuilder(t *testing.T) {
	testCases := []struct {
		id        string
		filter    *topologyv1.SearchRequest_Filter
		sort      *topologyv1.SearchRequest_Sort
		pageToken string
		limit     uint64
		expect    string
	}{
		{
			id:        "Default",
			filter:    &topologyv1.SearchRequest_Filter{},
			sort:      &topologyv1.SearchRequest_Sort{},
			pageToken: "0",
			limit:     0,
			expect:    "SELECT id, data, metadata FROM topology_cache ORDER BY ID ASC LIMIT 100 OFFSET 0",
		},
		{
			id:        "No page set",
			filter:    &topologyv1.SearchRequest_Filter{},
			sort:      &topologyv1.SearchRequest_Sort{},
			pageToken: "",
			limit:     5,
			expect:    "SELECT id, data, metadata FROM topology_cache ORDER BY ID ASC LIMIT 5 OFFSET 0",
		},
		{
			id:        "Page 0 with limit set",
			filter:    &topologyv1.SearchRequest_Filter{},
			sort:      &topologyv1.SearchRequest_Sort{},
			pageToken: "0",
			limit:     5,
			expect:    "SELECT id, data, metadata FROM topology_cache ORDER BY ID ASC LIMIT 5 OFFSET 0",
		},
		{
			id:        "Change PageToken and Limits",
			filter:    &topologyv1.SearchRequest_Filter{},
			sort:      &topologyv1.SearchRequest_Sort{},
			pageToken: "10",
			limit:     5,
			expect:    "SELECT id, data, metadata FROM topology_cache ORDER BY ID ASC LIMIT 5 OFFSET 50",
		},
		{
			id: "All Options",
			filter: &topologyv1.SearchRequest_Filter{
				Search: &topologyv1.SearchRequest_Filter_Search{
					Field: "metadata.search.field",
					Text:  "cat",
				},
				TypeUrl: "type.googleapis.com/clutch.aws.ec2.v1.AutoscalingGroup",
				Metadata: map[string]string{
					"label": "value",
				},
			},
			sort: &topologyv1.SearchRequest_Sort{
				Field:     "metadata.meow.iam.a.cat",
				Direction: topologyv1.SearchRequest_Sort_ASCENDING,
			},
			pageToken: "10",
			limit:     5,
			expect:    "SELECT id, data, metadata FROM topology_cache WHERE metadata->'search'->>'field' LIKE $1 AND resolver_type_url = $2 AND metadata @> $3::jsonb ORDER BY $4 ASC LIMIT 5 OFFSET 50",
		},
	}
	for _, test := range testCases {
		output, _, err := paginatedQueryBuilder(test.filter, test.sort, test.pageToken, test.limit)
		assert.NoError(t, err)

		sql, _, err := output.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, test.expect, sql)
	}
}

func TestFilterQueryBuilder(t *testing.T) {
	testCases := []struct {
		id          string
		input       *topologyv1.SearchRequest_Filter
		expect      string
		shouldError bool
	}{
		{
			id:          "No Input",
			input:       &topologyv1.SearchRequest_Filter{},
			expect:      "SELECT * FROM topology_cache",
			shouldError: false,
		},
		{
			id: "Try to sql inject",
			input: &topologyv1.SearchRequest_Filter{
				Search: &topologyv1.SearchRequest_Filter_Search{
					Field: "column.id NOT NULL;--",
					Text:  "cat",
				},
			},
			expect:      "",
			shouldError: true,
		},
		{
			id: "Search by column",
			input: &topologyv1.SearchRequest_Filter{
				Search: &topologyv1.SearchRequest_Filter_Search{
					Field: "column.id",
					Text:  "cat",
				},
			},
			expect:      "SELECT * FROM topology_cache WHERE id LIKE $1",
			shouldError: false,
		},
		{
			id: "Search by Metadata",
			input: &topologyv1.SearchRequest_Filter{
				Search: &topologyv1.SearchRequest_Filter_Search{
					Field: "metadata.label",
					Text:  "cat",
				},
			},
			expect:      "SELECT * FROM topology_cache WHERE metadata->>'label' LIKE $1",
			shouldError: false,
		},
		{
			id: "Search all options",
			input: &topologyv1.SearchRequest_Filter{
				Search: &topologyv1.SearchRequest_Filter_Search{
					Field: "metadata.label",
					Text:  "cat",
				},
				TypeUrl: "type.googleapis.com/clutch.aws.ec2.v1.AutoscalingGroup",
				Metadata: map[string]string{
					"label":  "value",
					"label2": "value2",
				},
			},
			expect:      "SELECT * FROM topology_cache WHERE metadata->>'label' LIKE $1 AND resolver_type_url = $2 AND metadata @> $3::jsonb",
			shouldError: false,
		},
	}

	for _, test := range testCases {
		selectBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
			Select("*").
			From("topology_cache")

		output, err := filterQueryBuilder(selectBuilder, test.input)
		if test.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			sql, _, err := output.ToSql()
			assert.NoError(t, err)
			assert.Equal(t, test.expect, sql)
		}
	}
}

func TestSortQueryBuilder(t *testing.T) {
	testCases := []struct {
		id          string
		input       *topologyv1.SearchRequest_Sort
		expect      string
		shouldError bool
	}{
		{
			id:          "Default to ID ASC",
			input:       &topologyv1.SearchRequest_Sort{},
			expect:      "SELECT * FROM topology_cache ORDER BY ID ASC",
			shouldError: false,
		},
		{
			id: "Try to sql inject",
			input: &topologyv1.SearchRequest_Sort{
				Field:     "column.id NOT NULL;--",
				Direction: topologyv1.SearchRequest_Sort_DESCENDING,
			},
			expect:      "",
			shouldError: true,
		},
		{
			id: "Sort by custom column and direction",
			input: &topologyv1.SearchRequest_Sort{
				Field:     "column.cat",
				Direction: topologyv1.SearchRequest_Sort_DESCENDING,
			},
			expect:      "SELECT * FROM topology_cache ORDER BY $1 DESC",
			shouldError: false,
		},
		{
			id: "Sort by custom metadata and direction",
			input: &topologyv1.SearchRequest_Sort{
				Field:     "metadata.meow",
				Direction: topologyv1.SearchRequest_Sort_ASCENDING,
			},
			expect:      "SELECT * FROM topology_cache ORDER BY $1 ASC",
			shouldError: false,
		},
		{
			id: "Sort by custom metadata deeply nested",
			input: &topologyv1.SearchRequest_Sort{
				Field:     "metadata.meow.iam.a.cat",
				Direction: topologyv1.SearchRequest_Sort_ASCENDING,
			},
			expect:      "SELECT * FROM topology_cache ORDER BY $1 ASC",
			shouldError: false,
		},
	}

	for _, test := range testCases {
		selectBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
			Select("*").
			From("topology_cache")

		output, err := sortQueryBuilder(selectBuilder, test.input)
		if test.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			sql, _, err := output.ToSql()
			assert.NoError(t, err)
			assert.Equal(t, test.expect, sql)
		}
	}
}

func TestValidateFilterSortField(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{
			input:       "column.';DROP TABLE;'",
			shouldError: true,
		},
		{
			input:       "column.id NOT NULL;--",
			shouldError: true,
		},
		{
			input:       "metadata.id.hello.meow",
			shouldError: false,
		},
	}

	for _, test := range testCases {
		err := validateFilterSortField(test.input)
		if test.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestGetFilterSortPrefixIdentifer(t *testing.T) {
	testCases := []struct {
		id          string
		input       string
		output      string
		shouldError bool
	}{
		{
			id:          "Column",
			input:       "column.my.id",
			output:      "column",
			shouldError: false,
		},
		{
			id:          "Metadata",
			input:       "metadata.my.id",
			output:      "metadata",
			shouldError: false,
		},
		{
			id:          "Unsupported identifer",
			input:       "meow.my.id",
			output:      "",
			shouldError: true,
		},
	}

	for _, test := range testCases {
		output, err := getFilterSortPrefixIdentifer(test.input)
		if test.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, test.output, output)
	}
}

func TestConvertMetadataToQuery(t *testing.T) {
	testCases := []struct {
		id          string
		input       string
		expect      string
		shouldError bool
	}{
		{
			id:          "Invalid input",
			input:       "metadata.",
			expect:      "",
			shouldError: true,
		},
		{
			id:          "top level field",
			input:       "metadata.toplevel",
			expect:      "metadata->>'toplevel'",
			shouldError: false,
		},
		{
			id:          "one level deep",
			input:       "metadata.toplevel.level1",
			expect:      "metadata->'toplevel'->>'level1'",
			shouldError: false,
		},
		{
			id:          "two levels deep",
			input:       "metadata.toplevel.level1.level2",
			expect:      "metadata->'toplevel'->'level1'->>'level2'",
			shouldError: false,
		},
		{
			id:          "two levels deep",
			input:       "data.level1.level2.idk:ok:cool",
			expect:      "data->'level1'->'level2'->>'idk:ok:cool'",
			shouldError: false,
		},
	}

	for _, test := range testCases {
		output, err := convertDataOrMetadataToQuery(test.input)
		if test.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, test.expect, output)
	}
}

func TestGetDirection(t *testing.T) {
	testCases := []struct {
		id     string
		input  topologyv1.SearchRequest_Sort_Direction
		expect string
	}{
		{
			id:     "ASCENDING",
			input:  topologyv1.SearchRequest_Sort_ASCENDING,
			expect: "ASC",
		},
		{
			id:     "DESCENDING",
			input:  topologyv1.SearchRequest_Sort_DESCENDING,
			expect: "DESC",
		},
		{
			id:     "Bad input",
			input:  topologyv1.SearchRequest_Sort_UNSPECIFIED,
			expect: "ASC",
		},
	}

	for _, tests := range testCases {
		output := getDirection(tests.input)
		assert.Equal(t, tests.expect, output)
	}
}
