package serverutils

import (
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// encodeCursor returns an encoded cursor given an id and position
func encodeCursor(id string, position int) string {
	raw := fmt.Sprintf("%s@%d", id, position)

	encoded := base64.StdEncoding.EncodeToString([]byte(raw))

	return encoded
}

// decodeCursor decodes an opaque cursor returning the id and position
func decodeCursor(cursor string) (string, int, error) { //nolint:unparam
	if cursor == "" {
		return "", 0, fmt.Errorf("invalid cursor")
	}

	bytes, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return "", 0, fmt.Errorf("invalid base64 cursor: %w", err)
	}

	result := strings.Split(string(bytes), "@")

	if len(result) != 2 {
		return "", 0, fmt.Errorf("invalid cursor")
	}

	position, err := strconv.Atoi(result[1])
	if err != nil {
		return "", 0, fmt.Errorf("cursor position not a number: %w", err)
	}

	return result[0], position, nil
}

// PageInfo is used to add pagination information to Relay edges.
type PageInfo struct {
	// Forward pagination
	HasNextPage bool
	EndCursor   string

	// Backward pagination
	HasPreviousPage bool
	StartCursor     string
}

// Pagination represents paging parameters
type PaginationInput struct {
	// Forward pagination arguments
	First *int
	After *string

	// Backward pagination arguments
	Last   *int
	Before *string
}

// ToLimitOffset translates PaginationInput pagination into limit/offset pagination
//
// Rules applied:
//   - forward pagination (first, after) → limit = first,
//     offset = afterOffset+1 (or 0 if after absent).
//   - backward pagination (last, before) → limit = last,
//     offset = max(beforeOffset-last, 0).
//
// Conflicting/invalid combinations return an error.
func (p PaginationInput) ToLimitOffset() (limit, offset int, err error) {
	err = p.Validate()
	if err != nil {
		return 0, 0, fmt.Errorf("invalid pagination input: %w", err)
	}

	switch {
	case p.IsForward():
		var position int

		if p.After != nil {
			_, position, err = decodeCursor(*p.After)
			if err != nil {
				return 0, 0, fmt.Errorf("invalid cursor: %w", err)
			}
		} else {
			position = 0
		}

		var offset int

		if position == 0 {
			offset = position
		} else {
			offset = position + 1
		}

		return *p.First, offset, nil

	case p.IsBackward():
		limit := *p.Last

		_, position, err := decodeCursor(*p.Before)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid cursor: %w", err)
		}

		offset := position - limit
		if offset < 0 {
			offset = 0
		}

		return limit, offset, nil
	}

	return
}

// ToPageNumber translates PaginationInput pagination into pageSize/page pagination
//
// Rules applied:
//   - forward pagination (first, after) → pageSize = first,
//     page = afterPage+1 (or 1 if after absent).
//   - backward pagination (last, before) → pageSize = last,
//     page = max(afterPage-last, 1).
//
// Conflicting/invalid combinations return an error.
func (p PaginationInput) ToPageNumber() (pageSize, page int, err error) {
	err = p.Validate()
	if err != nil {
		return 0, 0, fmt.Errorf("invalid pagination input: %w", err)
	}

	switch {
	case p.IsForward():
		var position int

		if p.After != nil {
			_, position, err = decodeCursor(*p.After)
			if err != nil {
				return 0, 0, fmt.Errorf("invalid cursor: %w", err)
			}
		} else {
			position = 0
		}

		page := position + 1

		return *p.First, page, nil

	case p.IsBackward():
		_, position, err := decodeCursor(*p.Before)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid cursor: %w", err)
		}

		page := position - 1

		if page < 1 {
			page = 1
		}

		return *p.Last, page, nil
	}

	return
}

// IsForward indicates whether it is forward pagination
func (p PaginationInput) IsForward() bool {
	return p.First != nil
}

// IsBackward indicates whether it is backward pagination
func (p PaginationInput) IsBackward() bool {
	return p.Last != nil
}

// Validate validates the pagination input
func (p *PaginationInput) Validate() error {
	// pagination has not been provided
	// set defaults for sanity :-)
	if !p.IsForward() && !p.IsBackward() {
		first := 20
		after := encodeCursor("", 0)

		p.First = &first
		p.After = &after

		return nil
	}

	if p.IsForward() && p.IsBackward() {
		return fmt.Errorf("cannot paginate both forward and backward in one request")
	}

	if p.IsForward() {
		if *p.First < 0 {
			return fmt.Errorf("first must be non-negative")
		}
	}

	if p.IsBackward() {
		if *p.Last < 0 {
			return fmt.Errorf("last must be non-negative")
		}

		if p.Before == nil {
			return fmt.Errorf("before cursor must be provided")
		}
	}

	return nil
}

// Represents interface for a Node
type Node interface {
	NodeID() string
}

// Edge wraps a node plus its cursor.
type Edge[T Node] struct {
	Cursor string `json:"cursor"`
	Node   T      `json:"node"`
}

// Connection wraps edges plus its metadata.
type Connection[T Node] struct {
	TotalCount int
	Edges      []Edge[T]
	PageInfo   PageInfo
}

// BuildLimitOffsetConnection turns a slice of items into Relay-compatible edges and PageInfo
// uses limit/offset/total counts
func BuildLimitOffsetConnection[T Node](items []T, offset, limit, total int) Connection[T] {
	edges := make([]Edge[T], len(items))

	for i, item := range items {
		cursor := encodeCursor(item.NodeID(), offset+i)
		edges[i] = Edge[T]{Cursor: cursor, Node: item}
	}

	var (
		startCur, endCur string
	)

	if len(edges) > 0 {
		startCur = edges[0].Cursor
		endCur = edges[len(edges)-1].Cursor
	}

	return Connection[T]{
		TotalCount: total,
		Edges:      edges,
		PageInfo: PageInfo{
			StartCursor:     startCur,
			EndCursor:       endCur,
			HasNextPage:     offset+limit < total,
			HasPreviousPage: offset > 0,
		},
	}
}

// BuildPageNumberConnection turns a slice of items into Relay-compatible edges and PageInfo
// uses currentPage/pageSize/total counts
func BuildPageNumberConnection[T Node](items []T, currentPage, pageSize, total int) Connection[T] {
	edges := make([]Edge[T], len(items))

	for i, item := range items {
		cursor := encodeCursor(item.NodeID(), currentPage)
		edges[i] = Edge[T]{Cursor: cursor, Node: item}
	}

	var (
		startCur, endCur string
	)

	if len(edges) > 0 {
		startCur = edges[0].Cursor
		endCur = edges[len(edges)-1].Cursor
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return Connection[T]{
		TotalCount: total,
		Edges:      edges,
		PageInfo: PageInfo{
			StartCursor:     startCur,
			EndCursor:       endCur,
			HasNextPage:     currentPage < totalPages,
			HasPreviousPage: currentPage > 1,
		},
	}
}
