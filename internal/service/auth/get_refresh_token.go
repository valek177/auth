package auth

import "context"

func (s *serv) GetRefreshToken(ctx context.Context) error { // req *descAuth.GetRefreshTokenRequest) (*descAuth.GetRefreshTokenResponse, error) {
	return nil

	// claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(refreshTokenSecretKey))
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
