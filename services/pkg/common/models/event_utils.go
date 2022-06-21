package models

import (
	"encoding/json"
	"reflect"

	event_api "github.com/abhisek/supply-chain-gateway/services/gen"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
)

func (m MetaEventWithAttributes) Serialize() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return []byte{}, err
	} else {
		return bytes, nil
	}
}

func newMetaEventWithAttributes(t string) MetaEventWithAttributes {
	return MetaEventWithAttributes{
		MetaEvent: MetaEvent{
			Type:    t,
			Version: EventSchemaVersion,
		},
		MetaAttributes: MetaAttributes{},
	}
}

func NewArtefactRequestEvent(a Artefact) DomainEvent[Artefact] {
	return DomainEvent[Artefact]{
		MetaEventWithAttributes: newMetaEventWithAttributes(EventTypeArtefactRequestSubject),
		Data:                    a,
	}
}

func NewArtefactResponseEvent(a Artefact) DomainEvent[Artefact] {
	return DomainEvent[Artefact]{
		MetaEventWithAttributes: newMetaEventWithAttributes(EventTypeArtefactResponseSubject),
		Data:                    a,
	}
}

type commonDomainEventBuilder[T any] struct{}

func NewDomainEventBuilder[T any]() DomainEventBuilder[T] {
	return &commonDomainEventBuilder[T]{}
}

func (b *commonDomainEventBuilder[T]) Created(model T) DomainEvent[T] {
	return b.event(model, DomainEventTypeCreated)
}

func (b *commonDomainEventBuilder[T]) Updated(model T) DomainEvent[T] {
	return b.event(model, DomainEventTypeUpdated)
}

func (b *commonDomainEventBuilder[T]) Deleted(model T) DomainEvent[T] {
	return b.event(model, DomainEventTypeDeleted)
}

func (b *commonDomainEventBuilder[T]) event(model T, operation string) DomainEvent[T] {
	return DomainEvent[T]{
		MetaEventWithAttributes: MetaEventWithAttributes{
			MetaEvent: MetaEvent{
				Type:    EventTypeDomainEvent,
				Version: EventSchemaVersion,
			},
			MetaAttributes: MetaAttributes{
				Attributes: map[string]string{
					"model":     reflect.TypeOf(model).Name(),
					"operation": operation,
				},
			},
		},
		Data: model,
	}
}

func (b *commonDomainEventBuilder[T]) From(v interface{}) (DomainEvent[T], error) {
	var event DomainEvent[T]
	err := utils.MapStruct(v, &event)

	if err != nil {
		return DomainEvent[T]{}, err
	}

	return event, nil
}

// Utils for new spec driven events
func eventUid() string {
	return utils.NewUniqueId()
}

func NewSpecEventHeader(event_type, source string) event_api.EventHeader {
	return event_api.EventHeader{
		Type:    event_type,
		Source:  source,
		Id:      eventUid(),
		Context: &event_api.EventContext{},
	}
}
