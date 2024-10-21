package dao

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

var ErrPersonNotFound = errors.New("person not found")

type PersonDAO interface {
	Save(p *Person) error
	GetPersonByName(name string) (*Person, error)
	GetPersonById(id string) (*Person, error)
	GetAllPersons() ([]*Person, error)
	DeletePerson(id string) error
	GetPersonByUID(uid string) (*Person, error)
}

type Person struct {
	ID        int
	UID       string
	Name      string
	Timestamp time.Time
	Traits    json.RawMessage
}

type PersonDAOImpl struct {
	db *sql.DB
}

func NewPersonDAO(db *sql.DB) (*PersonDAOImpl, error) {
	return &PersonDAOImpl{db: db}, nil
}

func (dao *PersonDAOImpl) Save(p *Person) error {
	query := `
		INSERT INTO persons (uid, name, timestamp, traits)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (uid) DO UPDATE
		SET name = EXCLUDED.name, timestamp = EXCLUDED.timestamp, traits = EXCLUDED.traits
		RETURNING id
	`
	return dao.db.QueryRow(query, p.UID, p.Name, p.Timestamp, p.Traits).Scan(&p.ID)
}

func (dao *PersonDAOImpl) GetPersonByName(name string) (*Person, error) {
	p := &Person{}
	query := `SELECT id, uid, name, timestamp, traits FROM persons WHERE name = $1`
	err := dao.db.QueryRow(query, name).Scan(&p.ID, &p.UID, &p.Name, &p.Timestamp, &p.Traits)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (dao *PersonDAOImpl) GetPersonById(id string) (*Person, error) {
	p := &Person{}
	query := `SELECT id, uid, name, timestamp, traits FROM persons WHERE id = $1`
	err := dao.db.QueryRow(query, id).Scan(&p.ID, &p.UID, &p.Name, &p.Timestamp, &p.Traits)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (dao *PersonDAOImpl) GetAllPersons() ([]*Person, error) {
	query := `SELECT id, uid, name, timestamp, traits FROM persons`
	rows, err := dao.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []*Person
	for rows.Next() {
		p := &Person{}
		err := rows.Scan(&p.ID, &p.UID, &p.Name, &p.Timestamp, &p.Traits)
		if err != nil {
			return nil, err
		}
		persons = append(persons, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return persons, nil
}

func (dao *PersonDAOImpl) DeletePerson(id string) error {
	query := `DELETE FROM persons WHERE uid = $1`
	_, err := dao.db.Exec(query, id)
	return err
}

func (dao *PersonDAOImpl) GetPersonByUID(uid string) (*Person, error) {
	p := &Person{}
	query := `SELECT id, uid, name, timestamp, traits FROM persons WHERE uid = $1`
	err := dao.db.QueryRow(query, uid).Scan(&p.ID, &p.UID, &p.Name, &p.Timestamp, &p.Traits)
	if err != nil {
		return nil, err
	}
	return p, nil
}
