package utils

import (
	"fmt"
	"database/sql"

	"packster/pkg/types"
)

func GetHosts(sqlConn *sql.DB) (map[string]types.Host, error) {
	if sqlConn == nil {
		return nil, fmt.Errorf("sqlConn is nil")
	}

	rows, err := sqlConn.Query(`SELECT id, url, host_type, application_id, secret FROM host`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hosts := make(map[string]types.Host)

	for rows.Next() {
		var h types.Host
		var hostType string
		var appId, secret sql.NullString

		err := rows.Scan(&h.Id, &h.Url, &hostType, &appId, &secret)
		if err != nil {
			return nil, err
		}

		ht, ok := types.HostTypeFromString(hostType)
		if !ok {
			return nil, fmt.Errorf("unknown host type: %s", hostType)
		}

		h.Type = ht
		h.ApplicationId = appId.String
		h.Secret = secret.String

		hosts[h.Url] = h
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

func HostExists(url string, sqlConn *sql.DB) (bool, int) {
	var id int
	err := sqlConn.QueryRow(`SELECT id FROM host WHERE url=$1`, url).Scan(&id)
	if err != nil {
		return false, 0
	}
	return true, id
}

func GetHostById(sqlConn *sql.DB, id int) (*types.Host, error) {
	if sqlConn == nil {
		return nil, fmt.Errorf("sqlConn is nil")
	}

	var h types.Host
	var hostType string
	var appId, secret sql.NullString

	err := sqlConn.QueryRow(`SELECT id, url, host_type, application_id, secret FROM host WHERE id=$1`, id).
		Scan(&h.Id, &h.Url, &hostType, &appId, &secret)
	if err != nil {
		return nil, err
	}

	ht, ok := types.HostTypeFromString(hostType)
	if !ok {
		return nil, fmt.Errorf("unknown host type: %s", hostType)
	}

	h.Type = ht
	h.ApplicationId = appId.String
	h.Secret = secret.String

	return &h, nil
}
