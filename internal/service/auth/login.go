package auth

import "context"

func (s *serv) Login(ctx context.Context) error {
	// Go to database through repo layer (GetUser)
	return nil

	// Лезем в базу или кэш за данными пользователя
	// Сверяем хэши пароля

	// refreshToken, err := utils.GenerateToken(model.UserInfo{
	// 	Username: req.GetUsername(),
	// 	// Это пример, в реальности роль должна браться из базы или кэша
	// 	Role: "admin",
	// },
	// 	[]byte(refreshTokenSecretKey),
	// 	refreshTokenExpiration,
	// )
	// if err != nil {
	// 	return nil, errors.New("failed to generate token")
	// }

	// return &descAuth.LoginResponse{RefreshToken: refreshToken}, nil
}
