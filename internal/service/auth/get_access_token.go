package auth

import "context"

// GetAccessToken returns access token by refresh token
func (s *serv) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	return "", nil

	// claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(refreshTokenSecretKey))
	// if err != nil {
	// 	return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	// }

	// // Можем слазать в базу или в кэш за доп данными пользователя

	// accessToken, err := utils.GenerateToken(model.UserInfo{
	// 	Username: claims.Username,
	// 	// Это пример, в реальности роль должна браться из базы или кэша
	// 	Role: "admin",
	// },
	// 	[]byte(accessTokenSecretKey),
	// 	accessTokenExpiration,
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// return &descAuth.GetAccessTokenResponse{AccessToken: accessToken}, nil
}
