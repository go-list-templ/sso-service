package vo

type AccessToken struct {
	value string
}

func NewAccessToken() AccessToken {
	return AccessToken{
		value: "",
	}
}

func UnsafeAccessToken(token string) AccessToken {
	return AccessToken{value: token}
}

func (a *AccessToken) Value() string {
	return a.value
}
