package repository

import (
	"context"
	"os"
	"testing"
	"time"

	cfg "simple-backend-nongki-go/config"
	auth "simple-backend-nongki-go/features/auth"
	pg "simple-backend-nongki-go/utils/driver/postgresql"
	generator "simple-backend-nongki-go/utils/generator"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	repoTest auth.RepositoryInterface
	pool     *pgxpool.Pool
	ctx      context.Context
)

func TestMain(m *testing.M) {
	os.Setenv("DB_USERNAME", "dwiw")
	os.Setenv("DB_PASSWORD", "secret")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "nongki_db")

	envConfig := &cfg.EnvConfig{
		DB_USERNAME: os.Getenv("DB_USERNAME"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_NAME:     os.Getenv("DB_NAME"),
	}

	pool = pg.ConnectToPg(envConfig)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	repoTest = NewAuthRepository(pool, ctx)

	os.Exit(m.Run())
}

func createRandomUser(t *testing.T) (user auth.User) {
	require.NotNil(t, pool)
	require.NotNil(t, ctx)

	user.FirstName = generator.CreateRandomString(int(generator.RandomInt(3, 13)))
	user.MiddleName = generator.CreateRandomString(int(generator.RandomInt(3, 13)))
	user.LastName = generator.CreateRandomString(int(generator.RandomInt(3, 13)))
	user.Email = generator.CreateRandomEmail(user.FirstName)
	user.Address = generator.CreateRandomString(int(generator.RandomInt(20, 50)))
	user.Gender = generator.CreateRandomGender()
	user.MaritalStatus = generator.CreateRandomMaritalStatus()
	user.HashedPassword = generator.CreateRandomString(int(generator.RandomInt(5, 10)))

	assert.NotEmpty(t, user.FirstName)
	assert.NotEmpty(t, user.LastName)
	assert.NotEmpty(t, user.LastName)
	assert.NotEmpty(t, user.Email)
	assert.NotEmpty(t, user.Address)
	assert.NotEmpty(t, user.Gender)
	assert.NotEmpty(t, user.MaritalStatus)
	assert.NotEmpty(t, user.HashedPassword)

	query := `
	INSERT INTO users(
		email,
		first_name,
		middle_name,
		last_name,
		address,
		gender,
		marital_status,
		hashed_password
	) VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8) 
	RETURNING id;`

	row := pool.QueryRow(ctx, query, user.Email, user.FirstName, user.MiddleName, user.LastName, user.Address, user.Gender, user.MaritalStatus, user.HashedPassword)
	err := row.Scan(&user.ID)
	require.NoError(t, err)

	assert.NotZero(t, user.ID)
	return
}

func TestCheckEmail(t *testing.T) {
	for i := 0; i < 5; i++ {
		t.Run("success", func(t *testing.T) {
			user := createRandomUser(t)

			res, err := repoTest.CheckEmail(user.Email)
			require.NoError(t, err)
			assert.NotZero(t, res)
		})
	}

	tests := []struct {
		name  string
		email string
	}{
		{
			name:  "error_empty_email",
			email: "",
		}, {
			name:  "error_invalid_email",
			email: "av088@mail.com",
		}, {
			name:  "error_typo_email",
			email: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			account := createRandomUser(t)

			if test.name == "error_typo_email" {
				test.email = account.Email + "m"
			}
			res, err := repoTest.ReadUser(test.email)
			require.Error(t, err)
			require.Zero(t, res)
		})
	}
}

func TestReadUser(t *testing.T) {
	for i := 0; i < 5; i++ {
		t.Run("success", func(t *testing.T) {
			user := createRandomUser(t)

			res, err := repoTest.ReadUser(user.Email)
			require.NoError(t, err)
			assert.NotZero(t, &res.ID)
			assert.Equal(t, user.Email, res.Email)
			assert.Equal(t, user.FirstName, res.FirstName)
			assert.Equal(t, user.MiddleName, res.MiddleName)
			assert.Equal(t, user.LastName, res.LastName)
			assert.Equal(t, user.Address, res.Address)
			assert.Equal(t, user.Gender, res.Gender)
			assert.Equal(t, user.MaritalStatus, res.MaritalStatus)
			assert.Equal(t, user.HashedPassword, res.HashedPassword)
		})
	}

	tests := []struct {
		name  string
		email string
	}{
		{
			name:  "error_empty_email",
			email: "",
		}, {
			name:  "error_invalid_email",
			email: "av088@mail.com",
		}, {
			name:  "error_typo_email",
			email: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			account := createRandomUser(t)

			if test.name == "error_typo_email" {
				test.email = account.Email + "m"
			}
			res, err := repoTest.ReadUser(test.email)
			require.Error(t, err)
			require.Nil(t, res)
		})
	}
}

func TestInsertUser(t *testing.T) {
	email := generator.CreateRandomEmail(generator.CreateRandomString(5))
	tests := []struct {
		name  string
		input auth.SignupRequest
		err   bool
	}{
		{
			name: "success",
			input: auth.SignupRequest{
				FirstName:      generator.CreateRandomString(5),
				LastName:       generator.CreateRandomString(7),
				Email:          email,
				Address:        generator.CreateRandomString(20),
				Gender:         generator.CreateRandomGender(),
				MaritalStatus:  generator.CreateRandomMaritalStatus(),
				HashedPassword: generator.CreateRandomString(60),
			},
			err: false,
		}, {
			name: "error_nil_first_name",
			input: auth.SignupRequest{
				LastName:       generator.CreateRandomString(7),
				Email:          generator.CreateRandomEmail(generator.CreateRandomString(5)),
				Address:        generator.CreateRandomString(20),
				Gender:         generator.CreateRandomGender(),
				MaritalStatus:  generator.CreateRandomMaritalStatus(),
				HashedPassword: generator.CreateRandomString(60),
			},
			err: true,
		}, {
			name: "error_empty_address",
			input: auth.SignupRequest{
				FirstName:      generator.CreateRandomEmail(generator.CreateRandomString(5)),
				LastName:       generator.CreateRandomString(7),
				Email:          generator.CreateRandomEmail(generator.CreateRandomString(5)),
				Address:        "",
				Gender:         generator.CreateRandomGender(),
				MaritalStatus:  generator.CreateRandomMaritalStatus(),
				HashedPassword: generator.CreateRandomString(60),
			},
			err: true,
		}, {
			name: "error_duplicate_email",
			input: auth.SignupRequest{
				FirstName:      generator.CreateRandomString(5),
				LastName:       generator.CreateRandomString(7),
				Email:          email,
				Address:        generator.CreateRandomString(20),
				Gender:         generator.CreateRandomGender(),
				MaritalStatus:  generator.CreateRandomMaritalStatus(),
				HashedPassword: generator.CreateRandomString(60),
			},
			err: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := repoTest.InsertUser(test.input)
			if !test.err {
				require.NoError(t, err)
				assert.Equal(t, test.input.FirstName, res.FirstName)
				assert.Equal(t, "", res.MiddleName)
				assert.Equal(t, test.input.LastName, res.LastName)
				assert.Equal(t, test.input.Email, res.Email)
				assert.Equal(t, test.input.Address, res.Address)
				assert.Equal(t, test.input.Gender, res.Gender)
				assert.Equal(t, test.input.MaritalStatus, res.MaritalStatus)
				assert.Equal(t, test.input.HashedPassword, res.HashedPassword)
			} else {
				require.Error(t, err)
			}
		})
	}
}
