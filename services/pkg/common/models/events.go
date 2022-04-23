package models

const (
	EventSchemaVersion               = "1.0.0"
	EventTypeArtefactRequestSubject  = "event.artefact.request"
	EventTypeArtefactResponseSubject = "event.artefact.response"
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
