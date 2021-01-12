package topology

import (
	"fmt"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
)

func paginatedQueryBuilder(
	filter *topologyv1.SearchTopologyRequest_Filter,
	sort *topologyv1.SearchTopologyRequest_Sort,
	pageToken string,
	limit int,
) (sq.SelectBuilder, error) {
	queryLimit := 100
	if limit >= 0 {
		queryLimit = limit
	}

	pageNum, err := strconv.Atoi(pageToken)
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	queryOffset := 0
	if pageNum > 0 {
		queryOffset = pageNum * limit
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
	if f.Search != nil && f.Search.Field != nil {
		if f.Search.Field.GetColumn() > 0 {
			query = query.Where(sq.Like{f.Search.Field.GetColumn().String(): fmt.Sprintf("%%%s%%", f.Search.Text)})
		} else if f.Metadata != nil {
			mdQuery := convertMetadataToQuery(f.Search.Field.GetMetadata())
			query = query.Where(sq.Like{mdQuery: f.Search.Text})
		}
	}

	if len(f.TypeUrl) > 0 {
		query = query.Where(sq.Eq{"resolver_type_url": f.TypeUrl})
	}

	return query
}

func sortQueryBuilder(query sq.SelectBuilder, s *topologyv1.SearchTopologyRequest_Sort) sq.SelectBuilder {
	if len(s.Direction.String()) > 0 && s.Field.Field != nil {
		direction := getDirection(s.Direction.String())

		if s.Field.GetColumn() > 0 {
			query = query.OrderBy(fmt.Sprintf("%s %s", s.Field.GetColumn(), direction))
		} else if len(s.Field.GetMetadata()) > 0 {
			mdQuery := convertMetadataToQuery(s.Field.GetMetadata())
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
	default:
		return ""
	}
}
