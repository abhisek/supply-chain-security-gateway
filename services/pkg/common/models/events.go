package models

const (
	EventSchemaVersion               = "1.0.0"
	EventTypeDomainEvent             = "event.domain"
	EventTypeArtefactRequestSubject  = "event.artefact.request"
	EventTypeArtefactResponseSubject = "event.artefact.response"

	DomainEventTypeCreated = "event.type.created"
	DomainEventTypeUpdated = "event.type.updated"
	DomainEventTypeDeleted = "event.type.deleted"
)

type MetaEvent struct {
	Version string `json:"version"`
	Type    string `json:"type"`
}

type MetaAttributes struct {
	Attributes map[string]string `json:"attributes"`
}

type MetaEventWithAttributes struct {
	MetaEvent
	MetaAttributes
}

type DomainEvent[T any] struct {
	MetaEventWithAttributes `json:"meta"`
	Data                    T `json:"data"`
}

type DomainEventBuilder[T any] interface {
	Created(v T) DomainEvent[T]
	Updated(v T) DomainEvent[T]
	Deleted(v T) DomainEvent[T]
	From(v interface{}) (DomainEvent[T], error)
}
