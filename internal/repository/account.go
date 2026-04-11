package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"packster/internal/utils"
	"packster/pkg/types"
)

var ErrAccountExists = errors.New("account already exists")

type IAccountRepo interface {
	CreateAccount(request types.AuthRequest) (*types.Account, error)
	AccountExists(username, sso string, host int) (*types.Account, error)
	GetDB() *sql.DB
}

type AccountRepo struct {
	SqlConn *sql.DB
}

func NewAccountRepo(sqlConn *sql.DB) *AccountRepo {
	return &AccountRepo{
		SqlConn: sqlConn,
	}
}

func (r *AccountRepo) GetDB() *sql.DB {
	return r.SqlConn
}

func (r *AccountRepo) CreateAccount(request types.AuthRequest) (*types.Account, error) {
	hostExists, hostId := utils.HostExists(request.Host, r.SqlConn)
	if !hostExists {
		return nil, fmt.Errorf("%s isnt a valid host", request.Host)
	}

	existing, err := r.AccountExists(request.Username, request.SsoId, hostId)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return existing, ErrAccountExists
	}

	tx, err := r.SqlConn.Begin()
	if err != nil {
		return nil, err
	}

	var accountID int
	err = tx.QueryRow(`INSERT INTO account (display_name, last_login, created_at) VALUES ($1, $2, $3) RETURNING id`,
		request.Username,
		time.Now(),
		time.Now(),
	).Scan(&accountID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec(`INSERT INTO auth (account, username, sso_id, host) VALUES ($1, $2, $3, $4)`,
		accountID,
		request.Username,
		request.SsoId,
		hostId,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	sso, err := strconv.Atoi(request.SsoId)
	if err != nil {
		return nil, err
	}

	account := &types.Account{
		DisplayName: request.Username,
		AuthData: types.Auth{
			Account: accountID,
			Username: request.Username,
			SsoId: sso,
			Host: hostId,
		},
	}
	return account, nil
}

func (r *AccountRepo) AccountExists(username, sso string, host int) (*types.Account, error) {
	var account types.Account
	err := r.SqlConn.QueryRow(
		`SELECT a.display_name, auth.account, auth.username, auth.sso_id, auth.host
		 FROM auth
		 JOIN account a ON a.id = auth.account
		 WHERE auth.username=$1 AND auth.sso_id=$2 AND auth.host=$3`,
		username,
		sso,
		host,
	).Scan(
		&account.DisplayName,
		&account.AuthData.Account,
		&account.AuthData.Username,
		&account.AuthData.SsoId,
		&account.AuthData.Host,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &account, nil
}
