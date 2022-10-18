package auth

type authIdentity struct {
	idType, userId, orgId, projectId, name string
}

func (a *authIdentity) Type() string {
	return a.idType
}

func (a *authIdentity) UserId() string {
	return a.userId
}

func (a *authIdentity) Name() string {
	return a.name
}

func (a *authIdentity) OrgId() string {
	return a.orgId
}

func (a *authIdentity) ProjectId() string {
	return a.projectId
}

func AnonymousIdentity() AuthenticatedIdentity {
	return &authIdentity{idType: AuthIdentityTypeAnonymous,
		userId: "anonymous", name: "Anonymous User"}
}
