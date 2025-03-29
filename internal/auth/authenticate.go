package auth

type Authentication struct {
	Ok      bool
	IdToken string
}

func Authenticate(userName string) *Authentication {
	return &Authentication{
		Ok:      true,
		IdToken: "token",
	}
}
