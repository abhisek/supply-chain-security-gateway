package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthCredential(t *testing.T) {
	cases := []struct {
		name            string
		inputUserId     string
		inputUserSecret string

		outputUserId     string
		outputOrgId      string
		outputProjectId  string
		outputUserSecret string
	}{
		{
			"Full userId format",
			"projectId/username@orgName",
			"userSecret",

			"username@orgName",
			"orgName",
			"projectId",
			"userSecret",
		},
		{
			"ProjectId is not given",
			"username@orgName",
			"userSecret",

			"username@orgName",
			"orgName",
			"",
			"userSecret",
		},
		{
			"ProjectId and OrgId is not given",
			"username",
			"userSecret",

			"username",
			"",
			"",
			"userSecret",
		},
		{
			"Starts with a Slash",
			"/projectId/username@orgName",
			"userSecret",

			"projectId/username@orgName",
			"orgName",
			"",
			"userSecret",
		},
		{
			"Double Slash in UserId",
			"projectId//username@orgName",
			"userSecret",

			"/username@orgName",
			"orgName",
			"projectId",
			"userSecret",
		},
		{
			"Double @ in UserId",
			"projectId/username@@orgName",
			"userSecret",

			"username@@orgName",
			"@orgName",
			"projectId",
			"userSecret",
		},
		{
			"Username ending with Slash",
			"projectId/username/@orgName",
			"userSecret",

			"username/@orgName",
			"orgName",
			"projectId",
			"userSecret",
		},
		{
			"User Secret has special chars",
			"projectId/username@orgName",
			"@@///@@/@@",

			"username@orgName",
			"orgName",
			"projectId",
			"@@///@@/@@",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			creds := authCredential{userId: test.inputUserId, userSecret: test.inputUserSecret}

			assert.Equal(t, creds.UserId(), test.outputUserId)
			assert.Equal(t, creds.OrgId(), test.outputOrgId)
			assert.Equal(t, creds.ProjectId(), test.outputProjectId)
			assert.Equal(t, creds.UserSecret(), test.outputUserSecret)
		})
	}
}
