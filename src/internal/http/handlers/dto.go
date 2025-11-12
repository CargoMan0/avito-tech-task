package handlers

type teamDTO struct {
	Name    string    `json:"team_name"`
	Members []userDTO `json:"members"`
}

func (t *teamDTO) IsValid() bool {

}

type userDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
