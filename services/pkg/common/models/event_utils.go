package models

import "encoding/json"

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
