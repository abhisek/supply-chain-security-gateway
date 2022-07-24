package route

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoute(t *testing.T) {
	cases := []struct {
		name          string
		pattern, path string
		is_match      bool
		labels        map[string]string
	}{
		{
			"Basic positive case",
			"/a/:something/b", "/a/b/b",
			true, map[string]string{"something": "b"},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			handler, err := NewRouteHandler(test.pattern)

			assert.Nil(t, err)
			assert.NotNil(t, handler)

			m := handler.Match(test.path)
			assert.Equal(t, test.is_match, m.IsMatch())
			assert.ElementsMatch(t, test.labels, m.Labels())
		})
	}
}
