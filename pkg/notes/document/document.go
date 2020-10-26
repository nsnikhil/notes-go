package document

import "time"

type Document struct {
	id   string
	name string

	content string

	createdAt time.Time
	updatedAt time.Time
}
