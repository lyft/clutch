package topology

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
)

const (
	column                   = "column."
	metadata                 = "metadata."
	queryDefaultLimit uint64 = 100
	maxResultLimit           = 1000
)

func paginatedQueryBuilder(
	filter *topologyv1.SearchRequest_Filter,
	sort *topologyv1.SearchRequest_Sort,
	pageToken string,
	limit uint64,
) (sq.SelectBuilder, uint64, error) {
	queryLimit := queryDefaultLimit
	if limit > maxResultLimit {
		return sq.SelectBuilder{}, 0, status.Error(codes.InvalidArgument, "maximum query limit is 1000")
	} else if limit > 0 {
		queryLimit = limit
	}

	// If no page is supplied default to 0
	var pageNum uint64 = 0
	var err error
	if len(pageToken) > 0 {
		pageNum, err = strconv.ParseUint(pageToken, 10, 64)
		if err != nil {
			return sq.SelectBuilder{}, 0, status.Error(codes.InvalidArgument, "unable to parse page_token")
		}
	}

	var queryOffset uint64 = 0
	if pageNum > 0 {
		queryOffset = pageNum * limit
	}

	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "data", "metadata").
		From("topology_cache").
		Limit(uint64(queryLimit)).
		Offset(uint64(queryOffset))

	if filter != nil {
		query, err = filterQueryBuilder(query, filter)
		if err != nil {
			return sq.SelectBuilder{}, 0, err
		}
	}

	if sort != nil {
		query, err = sortQueryBuilder(query, sort)
		if err != nil {
			return sq.SelectBuilder{}, 0, err
		}
	}

	// Blindly increment pageNum by 1 for next_page_token
	return query, pageNum + 1, nil
}

func filterQueryBuilder(query sq.SelectBuilder, f *topologyv1.SearchRequest_Filter) (sq.SelectBuilder, error) {
	if f == nil {
		return query, nil
	}

	if f.Search != nil && len(f.Search.Field) > 0 {
		searchIdentiferExpr := sq.Expr("id LIKE ?", f.Search.Text)
		identifer, err := getFilterSortPrefixIdentifer(f.Search.Field)
		if err != nil {
			return sq.SelectBuilder{}, err
		}

		if identifer == column {
			searchIdentiferExpr = sq.Expr(strings.TrimPrefix(f.Search.Field, column))
		} else if identifer == metadata {
			mdQuery := convertMetadataToQuery(strings.TrimPrefix(f.Search.Field, metadata))
			searchIdentiferExpr = sq.Expr(mdQuery)
		}

		query = query.Where(sq.Expr("quote_literal(?) LIKE ?", searchIdentiferExpr, fmt.Sprintf("%%%s%%", f.Search.Text)))
	}

	if len(f.TypeUrl) > 0 {
		query = query.Where(sq.Eq{"resolver_type_url": f.TypeUrl})
	}

	if f.Metadata != nil {
		metadataJson, err := json.Marshal(f.Metadata)
		if err != nil {
			return sq.SelectBuilder{}, err
		}
		query = query.Where(sq.Expr("metadata @> ?::jsonb", metadataJson))
	}

	return query, nil
}

func sortQueryBuilder(query sq.SelectBuilder, s *topologyv1.SearchRequest_Sort) (sq.SelectBuilder, error) {
	if len(s.Field) > 0 {
		direction := getDirection(s.Direction)
		identifer, err := getFilterSortPrefixIdentifer(s.Field)
		if err != nil {
			return sq.SelectBuilder{}, err
		}

		if identifer == column {
			query = query.OrderByClause(fmt.Sprintf("? %s", direction), strings.TrimPrefix(s.Field, column))
		} else if identifer == metadata {
			mdQuery := convertMetadataToQuery(strings.TrimPrefix(s.Field, metadata))
			query = query.OrderByClause(fmt.Sprintf("? %s", direction), mdQuery)
		}
	} else {
		query = query.OrderBy("ID ASC")
	}

	return query, nil
}

func getFilterSortPrefixIdentifer(identifer string) (string, error) {
	// appending the extra `.` so we can utilize the const defined for column and metadata
	identifer += "."

	switch identifer {
	case column:
		return column, nil
	case metadata:
		return metadata, nil
	default:
		return "", fmt.Errorf("Unsupported identifer: [%s]", identifer)
	}
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

func getDirection(direction topologyv1.SearchRequest_Sort_Direction) string {
	switch direction {
	case topologyv1.SearchRequest_Sort_ASCENDING:
		return "ASC"
	case topologyv1.SearchRequest_Sort_DESCENDING:
		return "DESC"
	default:
		return "ASC"
	}
}
