package entity

type GenerateTokensPayload struct {
	UserId string
}

type GenerateTokensResult struct {
	AccessToken  string
	RefreshToken string
}
