package auth

type Authorization struct {
	Ok bool
}

func Authorize(userName string) *Authorization {
	return &Authorization{
		Ok: true,
	}
}
