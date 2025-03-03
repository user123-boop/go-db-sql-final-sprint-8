package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p

	//db := s.db
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		//sql.Named("number", number),
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))

	if err != nil {
		fmt.Println(err)
	}

	// верните идентификатор последней добавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка

	stmt := `SELECT number, client, status, address, created_at FROM parcel WHERE number = :number`
	row := s.db.QueryRow(stmt, number)
	p := Parcel{}
	//row, err := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		fmt.Println(err)
	}

	//p := Parcel{}

	//err = rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	// заполните объект Parcel данными из таблицы

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	//stmt:=`SELECT number, client, status, address, created_at FROM parcel WHERE client = :client`
	// заполните срез Parcel данными из таблицы
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client", client)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	//срез для возвращаемых строк
	var res []Parcel

	for rows.Next() {
		ress := Parcel{}
		err := rows.Scan(&ress.Number, &ress.Client, &ress.Status, &ress.Address, &ress.CreatedAt)
		if err != nil {
			fmt.Println(err)
		}
		res = append(res, ress)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET number = :number, address = :address WHERE status = :ParcelStatusRegistered", number, address)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number))
	if err != nil {
		fmt.Println(err)
	}
	return err
}
