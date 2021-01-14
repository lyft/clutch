package topology

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
)

func TestPaginatedQueryBuilder(t *testing.T) {
	testCases := []struct {
		id        string
		filter    *topologyv1.SearchRequest_Filter
		sort      *topologyv1.SearchRequest_Sort
		pageToken string
		limit     int
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
			expect:    "SELECT id, data, metadata FROM topology_cache WHERE metadata->'search'->'field' LIKE $1 AND resolver_type_url = $2 AND metadata->>'label' = $3 ORDER BY metadata->'meow'->'iam'->'a'->'cat' ASC LIMIT 5 OFFSET 50",
		},
	}
	for _, test := range testCases {
		output, err := paginatedQueryBuilder(test.filter, test.sort, test.pageToken, test.limit)
		assert.NoError(t, err)

		sql, _, err := output.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, test.expect, sql)
	}
}

func TestFilterQueryBuilder(t *testing.T) {
	testCases := []struct {
		id     string
		input  *topologyv1.SearchRequest_Filter
		expect string
	}{
		{
			id:     "No Input",
			input:  &topologyv1.SearchRequest_Filter{},
			expect: "SELECT * FROM topology_cache",
		},
		{
			id: "Search by column",
			input: &topologyv1.SearchRequest_Filter{
				Search: &topologyv1.SearchRequest_Filter_Search{
					Field: "column.id",
					Text:  "cat",
				},
			},
			expect: "SELECT * FROM topology_cache WHERE id LIKE $1",
		},
		{
			id: "Search by Metadata",
			input: &topologyv1.SearchRequest_Filter{
				Search: &topologyv1.SearchRequest_Filter_Search{
					Field: "metadata.label",
					Text:  "cat",
				},
			},
			expect: "SELECT * FROM topology_cache WHERE metadata->>'label' LIKE $1",
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
			expect: "SELECT * FROM topology_cache WHERE metadata->>'label' LIKE $1 AND resolver_type_url = $2 AND metadata->>'label' = $3 AND metadata->>'label2' = $4",
		},
	}

	for _, test := range testCases {
		selectBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
			Select("*").
			From("topology_cache")

		output := filterQueryBuilder(selectBuilder, test.input)
		sql, _, err := output.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, test.expect, sql)
	}
}

func TestSortQueryBuilder(t *testing.T) {
	testCases := []struct {
		id     string
		input  *topologyv1.SearchRequest_Sort
		expect string
	}{
		{
			id:     "Default to ID ASC",
			input:  &topologyv1.SearchRequest_Sort{},
			expect: "SELECT * FROM topology_cache ORDER BY ID ASC",
		},
		{
			id: "Sort by custom column and direction",
			input: &topologyv1.SearchRequest_Sort{
				Field:     "column.cat",
				Direction: topologyv1.SearchRequest_Sort_DESCENDING,
			},
			expect: "SELECT * FROM topology_cache ORDER BY cat DESC",
		},
		{
			id: "Sort by custom metadata and direction",
			input: &topologyv1.SearchRequest_Sort{
				Field:     "metadata.meow",
				Direction: topologyv1.SearchRequest_Sort_ASCENDING,
			},
			expect: "SELECT * FROM topology_cache ORDER BY metadata->>'meow' ASC",
		},
		{
			id: "Sort by custom metadata deeply nested",
			input: &topologyv1.SearchRequest_Sort{
				Field:     "metadata.meow.iam.a.cat",
				Direction: topologyv1.SearchRequest_Sort_ASCENDING,
			},
			expect: "SELECT * FROM topology_cache ORDER BY metadata->'meow'->'iam'->'a'->'cat' ASC",
		},
	}

	for _, test := range testCases {
		selectBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
			Select("*").
			From("topology_cache")

		output := sortQueryBuilder(selectBuilder, test.input)
		sql, _, err := output.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, test.expect, sql)
	}
}

func TestConvertMetadataToQuery(t *testing.T) {
	testCases := []struct {
		id     string
		input  string
		expect string
	}{
		{
			id:     "top level field",
			input:  "toplevel",
			expect: "metadata->>'toplevel'",
		},
		{
			id:     "one level deep",
			input:  "toplevel.level1",
			expect: "metadata->'toplevel'->'level1'",
		},
		{
			id:     "two levels deep",
			input:  "toplevel.level1.level2",
			expect: "metadata->'toplevel'->'level1'->'level2'",
		},
	}

	for _, test := range testCases {
		output := convertMetadataToQuery(test.input)
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
