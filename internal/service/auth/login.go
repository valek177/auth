package auth

import (
	"context"
)

func (s *serv) Login(ctx context.Context, username, password string) (string, error) {
	// Go to database through repo layer (GetUser)

	// id := int64(52) // get by name?
	// user, err := s.userRepository.GetUser(ctx, id)
	// if err != nil {
	// 	return "", errors.New("unable to get user")
	// }

	// isPasswordsEqual := passwordLib.CheckPasswordHash(password, user.Password)

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

	return "", nil // refreshToken, nil
}
