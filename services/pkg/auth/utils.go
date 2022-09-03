package auth

type authCredential struct {
	userId     string
	userSecret string
}

func (a *authCredential) UserId() string {
	return a.userId
}

func (a *authCredential) UserSecret() string {
	return a.userSecret
}

type authIdentity struct {
	idType, id, name string
}

func NewAuthIdentity(idType, id, name string) AuthenticatedIdentity {
	return &authIdentity{idType: idType, id: id, name: name}
}

func (a *authIdentity) Type() string {
	return a.idType
}

func (a *authIdentity) Id() string {
	return a.id
}

func (a *authIdentity) Name() string {
	return a.name
}

func AnonymousIdentity() AuthenticatedIdentity {
	return NewAuthIdentity(AuthIdentityTypeAnonymous,
		"anonymous", "Anonymous Identity")
}
