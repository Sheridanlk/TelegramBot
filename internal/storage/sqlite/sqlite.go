package sqlite

import (
	"TelegramBot/internal/storage"
	"TelegramBot/lib/e"
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, e.Wrap("can't open database", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("can't connect database", err)
	}

	return &Storage{db: db}, nil
}

// Save saves player to storage.
func (s *Storage) Save(ctx context.Context, player *storage.Player) error {
	q := `INSERT INTO players (user_name, date_of_last_visit) VALUES (?, ?)`

	if _, err := s.db.ExecContext(ctx, q, player.Name, player.LastVisit); err != nil {
		return e.Wrap("can't add user", err)
	}

	return nil
}

// Update updates the player's last online date.
func (s *Storage) Update(ctx context.Context, player *storage.Player) error {
	q := `UPDATE players SET date_of_last_visit = ? WHERE user_name = ?`

	if _, err := s.db.ExecContext(ctx, q, player.LastVisit, player.Name); err != nil {
		return e.Wrap("can't update last onlite date", err)
	}

	return nil
}

// Remove removes player form storage.
func (s *Storage) Remove(ctx context.Context, player string) error {
	q := `DELETE FORM players WHERE user_name = ?`
	if _, err := s.db.ExecContext(ctx, q, player); err != nil {
		return e.Wrap("can't remove player", err)
	}
	return nil
}

// IsExists check if player exists in storage.
func (s *Storage) IsExists(ctx context.Context, player string) (bool, error) {
	q := `SELECT COUNT(*) FROM players WHERE user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, player).Scan(&count); err != nil {
		return false, e.Wrap("can't check if pages exist", err)
	}
	return count > 0, nil
}

// GetPlayersLastLogin returns all players and their last online time.
func (s *Storage) GetPlayersLastLogin(ctx context.Context) ([]storage.Player, error) {
	q := `SELECT user_name, date_of_last_visit FROM players`

	rows, err := s.db.QueryContext(ctx, q)

	if err != nil {
		return nil, e.Wrap("can't get players last login", err)
	}

	defer rows.Close()

	var players []storage.Player
	for rows.Next() {
		var (
			name      string
			lastLogin time.Time
		)
		if err := rows.Scan(&name, &lastLogin); err != nil {
			return nil, err
		}
		players = append(players, storage.Player{
			Name:      name,
			LastVisit: lastLogin,
		})
	}
	return players, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS players (user_name TEXT, date_of_last_visit DATE)`
	if _, err := s.db.ExecContext(ctx, q); err != nil {
		return e.Wrap("can't create table", err)
	}
	return nil
}
