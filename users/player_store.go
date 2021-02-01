package users

type PlayerStore interface {
	Register() (id string)
	isValidID(id string) bool
	ReRegister(id string) error
	Deactivate(id string)
	JoinGame(ids []string, gameID string)
	LeaveGame(ids []string)
	GetPlayerDetails(id string) (*Player, error)
	GetAvailablePlayers(playerIDs []string) (availablePlayerIDs []string)
	CountActivePlayers() int
	getAllPlayers() (ids []string)
}
