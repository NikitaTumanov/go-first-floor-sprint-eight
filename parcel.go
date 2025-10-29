package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at);",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return -1, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(lastID), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	res := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number;",
		sql.Named("number", number))

	p := Parcel{}
	err := res.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	res, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client;",
		sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var result []Parcel
	for res.Next() {
		p := Parcel{}
		err := res.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	return result, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number;",
		sql.Named("number", number),
		sql.Named("status", status))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = 'registered';",
		sql.Named("number", number),
		sql.Named("address", address))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = 'registered';",
		sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}
