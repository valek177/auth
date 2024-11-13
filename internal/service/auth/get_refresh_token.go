package auth

import "context"

// GetRefreshToken returns new refresh token by old refresh token
func (s *serv) GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error) {
	return "", nil

	// claims, err := utils.VerifyToken(oldRefreshToken, []byte(refreshTokenSecretKey))
	// if err != nil {
	// 	return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	// }

	// // Можем слазать в базу или в кэш за доп данными пользователя

	// refreshToken, err := utils.GenerateToken(model.UserInfo{
	// 	Username: claims.Username,
	// 	// Это пример, в реальности роль должна браться из базы или кэша
	// 	Role: "admin",
	// },
	// 	[]byte(refreshTokenSecretKey),
	// 	refreshTokenExpiration,
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// return &descAuth.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
}
