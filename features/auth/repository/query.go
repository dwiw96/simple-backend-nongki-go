package repository

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"log"

	auth "simple-backend-nongki-go/features/auth"

	"github.com/jackc/pgx/v5/pgxpool"
)

type authRepository struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func NewAuthRepository(pool *pgxpool.Pool, ctx context.Context) auth.RepositoryInterface {
	return &authRepository{
		pool: pool,
		ctx:  ctx,
	}
}

func (r *authRepository) CheckEmail(email string) (result int, err error) {
	query := "SELECT COUNT(email) FROM users WHERE email=$1"

	row := r.pool.QueryRow(r.ctx, query, email)

	var count int
	err = row.Scan(&count)
	if err != nil {
		errMsg := fmt.Errorf("database error")
		log.Printf("CheckEmail() %v, err: %v\n", errMsg, err)
		return -1, errMsg
	}

	return count, nil
}

func (r *authRepository) ReadUser(email string) (result *auth.User, err error) {
	query := "SELECT * FROM users WHERE email = $1;"

	row := r.pool.QueryRow(r.ctx, query, email)

	var user auth.User
	err = row.Scan(&user.ID, &user.Email, &user.FirstName, &user.MiddleName, &user.LastName, &user.Address, &user.Gender, &user.MaritalStatus, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		errMsg := fmt.Errorf("database error")
		log.Printf("ReadUser() %v, err: %v\n", errMsg, err)
		return nil, errMsg
	}

	return &user, nil
}

func (r *authRepository) InsertUser(input auth.SignupRequest) (result *auth.User, err error) {
	query := `INSERT INTO users(
		email,
		first_name,
		middle_name,
		last_name,
		address,
		gender,
		marital_status,
		hashed_password
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8
	) RETURNING *`

	row := r.pool.QueryRow(r.ctx, query, input.Email, input.FirstName, input.MiddleName, input.LastName, input.Address, input.Gender, input.MaritalStatus, input.HashedPassword)

	var user auth.User
	err = row.Scan(&user.ID, &user.Email, &user.FirstName, &user.MiddleName, &user.LastName, &user.Address, &user.Gender, &user.MaritalStatus, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		errMsg := fmt.Errorf("database error")
		log.Printf("InsertUser() %v, err: %v\n", errMsg, err)
		return nil, errMsg
	}

	return &user, nil
}

func (r *authRepository) LoadKey() (key *rsa.PrivateKey, err error) {
	query := "select private_key from sec_m"
	var keyBytes []byte
	rows, err := r.pool.Query(r.ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&keyBytes)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		log.Println("Load success")

		privateKey, err := x509.ParsePKCS1PrivateKey(keyBytes)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return privateKey, nil
	}

	return nil, errors.New("no private key found in database")
}
