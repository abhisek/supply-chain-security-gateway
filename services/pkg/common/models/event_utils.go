package models

import (
	"encoding/json"

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

// Utils for new spec driven events
func eventUid() string {
	return utils.NewUniqueId()
}

func NewSpecEventHeader(tp event_api.EventType, source string) event_api.EventHeader {
	return event_api.EventHeader{
		Type:    tp,
		Source:  source,
		Id:      eventUid(),
		Context: &event_api.EventContext{},
	}
}
