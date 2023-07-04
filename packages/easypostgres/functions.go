package easypostgres

import (
	"database/sql"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"time"
)

// Admin
type Admin struct {
	ID        int64  `db:"id" json:"id"`
	Username  string `db:"username" json:"username"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"password"`
	Avatar    string `db:"avatar" json:"avatar"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type AdminJwt struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"created_at"`
	jwt.RegisteredClaims
}

func (p *PostgreSQL) VerifyExist(email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT * FROM admin WHERE email = $1)"
	err := p.DB.QueryRow(query, email).Scan(&exists)

	if err != nil {
		return exists, err
	}

	return exists, nil
}

func (p *PostgreSQL) SignUp(a *Admin) error {
	query := "INSERT INTO admin (username, email, password, avatar, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := p.DB.Exec(query, a.Username, a.Email, a.Password, a.Avatar, a.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQL) GetById(a *Admin) (*Admin, error) {
	query := "SELECT * FROM admin WHERE email = $1"
	rows := p.DB.QueryRow(query, a.ID)

	err := rows.Scan(a.ID, a.Username, a.Email, a.Password, a.Avatar, a.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("not Found")
		}
	}

	return a, err
}

func (p *PostgreSQL) Login(a *Admin) (*Admin, error) {
	query := "SELECT * FROM admin WHERE email = $1"
	rows := p.DB.QueryRow(query, a.Email)

	err := rows.Scan(&a.ID, &a.Username, &a.Email, &a.Password, &a.Avatar, &a.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("not Found")
		}
	}

	return a, err
}

func (a *Admin) FormatAdminToJWT(expiration time.Time) *AdminJwt {
	aj := &AdminJwt{
		ID:        a.ID,
		Username:  a.Username,
		Email:     a.Email,
		Avatar:    a.Avatar,
		CreatedAt: a.CreatedAt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			Issuer:    "test",
		},
	}

	return aj
}

func (aj *AdminJwt) JwtToReqVar() int64 {
	return aj.ID
}

// Movies

type Movie struct {
	ID          int64    `db:"id" json:"id"`
	Title       string   `db:"title" json:"title"`
	ReleaseDate string   `db:"release_date" json:"release_date"`
	Duration    int      `db:"duration" json:"duration"`
	Synopsis    string   `db:"synopsis" json:"synopsis"`
	Realisator  []string `db:"realisator" json:"realisator"`
	Productor   []string `db:"productor" json:"productor"`
	Actor       []string `db:"actor" json:"actor"`
	Picture     string   `db:"picture" json:"picture"`
	TrailerCode string   `db:"trailer_code" json:"trailer_code"`
}

func (p PostgreSQL) GetMovies() ([]*Movie, error) {
	movies := []*Movie{}

	query := `SELECT * FROM movies ORDER BY release_date DESC`
	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		movie := Movie{}

		err = rows.Scan(&movie.ID, &movie.Title, &movie.ReleaseDate, &movie.Duration, &movie.Synopsis, pq.Array(&movie.Realisator), pq.Array(&movie.Productor), pq.Array(&movie.Actor), &movie.Picture, &movie.TrailerCode)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (p PostgreSQL) GetMovie(id int) (*Movie, error) {
	movie := Movie{}
	query := `SELECT * FROM movies WHERE id = $1`
	rows := p.DB.QueryRow(query, id)

	err := rows.Scan(&movie.ID, &movie.Title, &movie.ReleaseDate, &movie.Duration, &movie.Synopsis, pq.Array(&movie.Realisator), pq.Array(&movie.Productor), pq.Array(&movie.Actor), &movie.Picture, &movie.TrailerCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("not Found")
		}
	}

	return &movie, err
}

func (p PostgreSQL) CreateMovie(m *Movie) (int64, error) {
	query := `INSERT INTO movies (title, release_date, duration, synopsis, realisator, productor, actor, picture, trailer_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	realisator := pq.Array(m.Realisator)
	productor := pq.Array(m.Productor)
	actor := pq.Array(m.Actor)

	var id int64
	err := p.DB.QueryRow(query, m.Title, m.ReleaseDate, m.Duration, m.Synopsis, realisator, productor, actor, m.Picture, m.TrailerCode).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
