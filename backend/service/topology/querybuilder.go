package topology

import (
	"fmt"
	"log"
	"strings"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
)

func paginatedQueryBuilder(
	filter *topologyv1.SearchTopologyRequest_Filter,
	sort *topologyv1.SearchTopologyRequest_Sort,
	skip, limit int,
) (string, error) {
	query := "SELECT id, data, metadata FROM topology_cache WHERE"

	queryLimit := 100
	if limit >= 0 {
		queryLimit = limit
	}

	queryOffset := 0
	if skip >= 0 {
		queryOffset = skip
	}

	f := filterQueryBuilder(filter)
	s := sortQueryBuilder(sort)
	l := fmt.Sprintf("LIMIT %d", queryLimit)
	o := fmt.Sprintf("OFFSET %d", queryOffset)

	fullQuery := strings.Join([]string{query, f, s, l, o, ";"}, " ")

	log.Printf("%s", fullQuery)
	log.Printf("%s", fullQuery)
	log.Printf("%s", fullQuery)

	return fullQuery, nil
}

func filterQueryBuilder(f *topologyv1.SearchTopologyRequest_Filter) string {
	filterQuery := ""

	if f == nil {
		return ""
	}

	if f.Search != nil && f.Search.Field != nil {
		if len(f.Search.Field.GetId()) > 0 {
			filterQuery += fmt.Sprintf("id like '%%%s%%'", f.Search.Text)
		} else if f.Metadata != nil {
			mdQuery := convertMetadataToQuery(f.Search.Field.GetMetadata())
			filterQuery += fmt.Sprintf("%s like '%%%s%%'", mdQuery, f.Search.Text)
		}
	}

	if len(f.TypeUrl) > 0 {
		if len(filterQuery) > 0 {
			filterQuery += " AND "
		}
		filterQuery += fmt.Sprintf("resolver_type_url = '%s'", f.TypeUrl)
	}

	return filterQuery
}

func sortQueryBuilder(s *topologyv1.SearchTopologyRequest_Sort) string {
	// Default is to sort by id
	sortQuery := "ORDER BY id ASC"

	if len(s.Direction.String()) > 0 && s.Field.Field != nil {
		direction := getDirection(s.Direction.String())

		if len(s.Field.GetId()) > 0 {
			sortQuery = fmt.Sprintf("ORDER BY %s %s", s.Field.GetId(), direction)
		} else if len(s.Field.GetMetadata()) > 0 {
			mdQuery := convertMetadataToQuery(s.Field.GetMetadata())
			sortQuery = fmt.Sprintf("ORDER BY %s %s", mdQuery, direction)
		}
	}

	return sortQuery
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
