package topology

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
)

const (
	column   = "column."
	metadata = "metadata."
)

func paginatedQueryBuilder(
	filter *topologyv1.SearchTopologyRequest_Filter,
	sort *topologyv1.SearchTopologyRequest_Sort,
	pageToken,
	limit int,
) (sq.SelectBuilder, error) {
	queryLimit := 100
	if limit > 0 {
		queryLimit = limit
	}

	queryOffset := 0
	if pageToken > 0 {
		queryOffset = pageToken * limit
	}

	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id, data, metadata").
		From("topology_cache").
		Limit(uint64(queryLimit)).
		Offset(uint64(queryOffset))

	query = filterQueryBuilder(query, filter)
	query = sortQueryBuilder(query, sort)

	return query, nil
}

func filterQueryBuilder(query sq.SelectBuilder, f *topologyv1.SearchTopologyRequest_Filter) sq.SelectBuilder {
	if f.Search != nil && len(f.Search.Field) > 0 {
		if strings.HasPrefix(f.Search.Field, column) {
			query = query.Where(sq.Like{strings.TrimPrefix(f.Search.Field, column): fmt.Sprintf("%%%s%%", f.Search.Text)})
		} else if strings.HasPrefix(f.Search.Field, metadata) {
			mdQuery := convertMetadataToQuery(strings.TrimPrefix(f.Search.Field, metadata))
			query = query.Where(sq.Like{mdQuery: fmt.Sprintf("%%%s%%", f.Search.Text)})
		}
	}

	if len(f.TypeUrl) > 0 {
		query = query.Where(sq.Eq{"resolver_type_url": f.TypeUrl})
	}

	// TODO: Support nested objects
	if f.Metadata != nil {
		for k, v := range f.Metadata {
			query = query.Where(sq.Eq{fmt.Sprintf("metadata->>'%s'", k): v})
		}
	}

	return query
}

func sortQueryBuilder(query sq.SelectBuilder, s *topologyv1.SearchTopologyRequest_Sort) sq.SelectBuilder {
	if len(s.Direction.String()) > 0 && len(s.Field) > 0 {
		direction := getDirection(s.Direction.String())

		if strings.HasPrefix(s.Field, column) {
			query = query.OrderBy(fmt.Sprintf("%s %s", strings.TrimPrefix(s.Field, column), direction))
		} else if strings.HasPrefix(s.Field, metadata) {
			mdQuery := convertMetadataToQuery(strings.TrimPrefix(s.Field, metadata))
			query = query.OrderBy(fmt.Sprintf("%s %s", mdQuery, direction))
		}
	} else {
		query = query.OrderBy("ID ASC")
	}

	return query
}

func convertMetadataToQuery(metadata string) string {
	metadataQuery := ""
	splitMetadata := strings.Split(metadata, ".")

	if len(splitMetadata) == 1 {
		metadataQuery = fmt.Sprintf("metadata->>'%s'", splitMetadata[0])
	} else {
		for i := range splitMetadata {
			splitMetadata[i] = fmt.Sprintf("'%s'", splitMetadata[i])
		}
		metadataQuery = fmt.Sprintf("metadata->%s", strings.Join(splitMetadata, "->"))
	}

	return metadataQuery
}

func getDirection(direction string) string {
	switch direction {
	case topologyv1.SearchTopologyRequest_Sort_ASCENDING.String():
		return "ASC"
	case topologyv1.SearchTopologyRequest_Sort_DESCENDING.String():
		return "DESC"
	case topologyv1.SearchTopologyRequest_Sort_UNSPECIFIED.String():
		// Default to ASC
		return "ASC"
	default:
		// Default to ASC
		return "ASC"
	}
}
