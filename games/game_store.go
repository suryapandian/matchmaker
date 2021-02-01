package games

type GameStore interface {
	save(*GameSession)
	delete(id string)
	GetGameByID(id string) (*GameSession, error)
	ActiveGames() []*GameSession
}
