package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user"
)

type UserRep struct {
	dbConn *sql.DB
}

func NewUserRep(conn *sql.DB) user.UserRep {
	return &UserRep{
		dbConn: conn,
	}
}

func (ur *UserRep) Insert(name string, email string, password string) (*models.User, error) {
	user := &models.User{}
	if err := ur.dbConn.QueryRow(
		"INSERT INTO users (id, name, email, password) VALUES (default, $1, $2, $3) RETURNING id, name, email, password, avatar;",
		name,
		email,
		password,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Avatar,
	); err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRep) Update(user *models.User) error {
	if err := ur.dbConn.QueryRow(
		"UPDATE users SET name = $1, email = $2, password = $3, avatar = $4 WHERE id = $5 RETURNING id, name, email, password, avatar;",
		user.Name,
		user.Email,
		user.Password,
		user.Avatar,
		user.ID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Avatar,
	); err != nil {
		return err
	}
	return nil
}

func (ur *UserRep) SelectById(id uint64) (*models.User, error) {
	user := &models.User{}
	if err := ur.dbConn.QueryRow(
		"SELECT id, name, email, password, avatar FROM users WHERE id = $1;", //TODO: точно ли есть смысл извлекать id, если он известен? везде так
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Avatar,
	); err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRep) SelectByLogin(login string) (*models.User, error) {
	user := &models.User{}
	if err := ur.dbConn.QueryRow(
		"SELECT id, name, email, password, avatar FROM users WHERE name = $1 OR email = $1;",
		login,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Avatar,
	); err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRep) SelectByNameOrEmail(name string, email string) ([]*models.User, error) {
	rows, err := ur.dbConn.Query(
		"SELECT id, name, email, password, avatar FROM users WHERE name = $1 OR email = $2;",
		name,
		email,
	)
	//defer func() {
	//	_ = rows.Close()
	//}()
	if err != nil {
		return nil, err
	}
	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Avatar,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (ur *UserRep) SelectByName(name string, authUserId uint64) (*models.User, error) {
	user := &models.User{}
	var isSubscribed sql.NullInt64
	if err := ur.dbConn.QueryRow(`
	SELECT users.id, users.name, users.avatar, user_subscriber.user_id FROM users
		LEFT JOIN user_subscriber ON user_subscriber.user_subscriber_id = $2 AND user_subscriber.user_id = users.id
	WHERE users.name = $1;`,
		name, authUserId).Scan(&user.ID, &user.Name, &user.Avatar, &isSubscribed); err != nil {
		return nil, err
	}

	if isSubscribed.Valid {
		user.IsSubscribed = true
	}

	return user, nil
}

func (ur *UserRep) SelectIfAdmin(userID uint64) (bool, error) {
	query := "SELECT is_admin FROM users WHERE id = $1"

	isAdmin := false
	err := ur.dbConn.QueryRow(query, userID).Scan(&isAdmin)

	if err != nil {
		return isAdmin, err
	}

	return isAdmin, nil
}

func (ur *UserRep) InsertSubscription(userSubscriberId uint64, userName string) error {
	var insertedUserId uint64
	return ur.dbConn.QueryRow(`INSERT INTO user_subscriber (user_subscriber_id, user_id)
		SELECT $1, users.id FROM users WHERE users.name = $2
		RETURNING user_id;`, userSubscriberId, userName).Scan(&insertedUserId)
}

func (ur *UserRep) RemoveSubscription(userSubscriberId uint64, userName string) error {
	_, err := ur.dbConn.Exec(`DELETE FROM user_subscriber
		WHERE user_subscriber_id = $1 AND user_id = (SELECT users.id FROM users WHERE users.name = $2);`,
		userSubscriberId, userName)
	return err
}

func (ur *UserRep) SelectSubscriptions(id, authUserId uint64) (*models.Subscriptions, error) {
	rows, err := ur.dbConn.Query(`
	SELECT users.id, users.name, users.avatar, user_subscriber2.user_id FROM user_subscriber AS user_subscriber1
		JOIN users ON user_subscriber1.user_subscriber_id = users.id
		LEFT JOIN user_subscriber AS user_subscriber2 ON user_subscriber2.user_subscriber_id = $2 AND user_subscriber2.user_id = users.id
	WHERE user_subscriber1.user_id = $1 ORDER BY users.name;`, id, authUserId)
	defer func() {
		_ = rows.Close()
	}()
	if err != nil {
		return nil, err
	}
	subscriptions := &models.Subscriptions{}
	for rows.Next() {
		user := &models.User{}
		var isSubscribed sql.NullInt64
		if err := rows.Scan(&user.ID, &user.Name, &user.Avatar, &isSubscribed); err != nil {
			return nil, err
		}

		if isSubscribed.Valid {
			user.IsSubscribed = true
		}

		subscriptions.Subscribers = append(subscriptions.Subscribers, user)
	}

	rows, err = ur.dbConn.Query(`
	SELECT users.id, users.name, users.avatar, user_subscriber2.user_id FROM user_subscriber AS user_subscriber1
		JOIN users ON user_subscriber1.user_id = users.id
		LEFT JOIN user_subscriber AS user_subscriber2 ON user_subscriber2.user_subscriber_id = $2 AND user_subscriber2.user_id = users.id
	WHERE user_subscriber1.user_subscriber_id = $1 ORDER BY users.name;`, id, authUserId)
	defer func() {
		_ = rows.Close()
	}()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		user := &models.User{}
		var isSubscribed sql.NullInt64
		if err := rows.Scan(&user.ID, &user.Name, &user.Avatar, &isSubscribed); err != nil {
			return nil, err
		}

		if isSubscribed.Valid {
			user.IsSubscribed = true
		}

		subscriptions.Subscriptions = append(subscriptions.Subscriptions, user)
	}
	return subscriptions, nil
}
