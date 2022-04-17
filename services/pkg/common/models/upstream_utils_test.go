package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPath2Artefact(t *testing.T) {
	cases := []struct {
		// Input
		prefix string
		upType string
		path   string

		// Output
		group, name, version string
		err                  error
	}{
		{
			"/maven2", ArtefactSourceTypeMaven2, "/maven2/com/google/guava/guava/30.1.1-jre/guava-30.1.1-jre.pom",
			"com.google.guava", "guava", "30.1.1-jre", nil,
		},
		{
			"/", ArtefactSourceTypeMaven2, "/com/google/guava/guava/30.1.1-jre/guava-30.1.1-jre.pom",
			"com.google.guava", "guava", "30.1.1-jre", nil,
		},
		{
			"/maven2", ArtefactSourceTypeMaven2, "/maven2/com/google/guava",
			"", "", "", errIncorrectMaven2Path,
		},
		{
			"/maven2", ArtefactSourceTypeMaven2, "",
			"", "", "", errIncorrectPrefix,
		},
		{
			"/maven2", ArtefactSourceTypeMaven2, "/maven2",
			"", "", "", errIncorrectMaven2Path,
		},
		{
			"/maven2", ArtefactSourceTypeMaven2, "/maven2/////",
			"", "", "", errIncorrectMaven2Path,
		},
		{
			"/maven2", ArtefactSourceTypeMaven2, "/maven2/com/google/guava/guava/../guava2/30.1.1-jre/guava-30.1.1-jre.pom",
			"com.google.guava", "guava2", "30.1.1-jre", nil,
		},
		{
			"/maven2", ArtefactSourceTypeMaven2, "/maven2/com/google/guava/guava/../../../../../m/x/y/z",
			"", "", "", errIncorrectPrefix,
		},
	}

	for _, test := range cases {
		upstream := ArtefactUpStream{
			Type: test.upType,
			RoutingRule: ArtefactRoutingRule{
				Prefix: test.prefix,
			},
		}

		artefact, err := upstream.Path2Artefact(test.path)

		if test.err != nil {
			assert.NotNil(t, err)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			}

		} else {
			assert.Nil(t, err)
			assert.Equal(t, test.group, artefact.Group)
			assert.Equal(t, test.name, artefact.Name)
			assert.Equal(t, test.version, artefact.Version)
		}
	}
}
