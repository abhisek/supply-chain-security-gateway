package auth

import "strings"

// Represents an user supplied credential
type authCredential struct {
	userId     string
	userSecret string
}

// Handle the form `project-id/user@org` ignoring the `project-id`
// while returning the userId
func (a *authCredential) UserId() string {
	userId := a.userId
	if strings.Index(userId, "/") >= 0 {
		userId = strings.SplitN(userId, "/", 2)[1]
	}

	return userId
}

func (a *authCredential) ProjectId() string {
	if strings.Index(a.userId, "/") >= 0 {
		return strings.SplitN(a.userId, "/", 2)[0]
	}

	return ""
}

func (a *authCredential) OrgId() string {
	uId := a.UserId()
	if strings.Index(uId, "@") >= 0 {
		return strings.SplitN(uId, "@", 2)[1]
	}

	return ""
}

func (a *authCredential) UserSecret() string {
	return a.userSecret
}
