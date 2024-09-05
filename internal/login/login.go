package login

func AuthToken() (string, error) {
	tokenSource := GetDebrickedTokenSource()
	token, err := tokenSource.Token()
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}
