package repositories

// func (ur *UserPGImpl) ChangeUser(user *models.User) error {
// 	query := "update users set name = $1, email = $2, password = $3, avatar = $4 where id = $5" +
// 		"returning id"

// 	err := ur.dbConn.QueryRow(query,
// 		user.Name,
// 		user.Email,
// 		user.Password,
// 		user.Avatar,
// 		user.ID).Scan(
// 		&user.ID)

// 	switch err {
// 	case sql.ErrNoRows:
// 		return errors.New("User not exists")
// 	case nil:
// 		return nil
// 	default:
// 		return err
// 	}
// }
