package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Cursor struct {
	ID        string `json:"id"`
	SortValue string `json:"sort_value"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	NextCursor string      `json:"next_cursor,omitempty"`
	HasMore    bool        `json:"has_more"`
}

func EncodeCursor(id, sortValue string) string {
	c := Cursor{ID: id, SortValue: sortValue}
	b, _ := json.Marshal(c)
	return base64.URLEncoding.EncodeToString(b)
}

func DecodeCursor(cursor string) (*Cursor, error) {
	b, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("decode cursor: %w", err)
	}
	var c Cursor
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("unmarshal cursor: %w", err)
	}
	return &c, nil
}
