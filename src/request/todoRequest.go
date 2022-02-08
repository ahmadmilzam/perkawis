package request

type (
	CreateTodo struct {
		Title           string `json:"title"`
		ActivityGroupID uint64 `json:"activity_group_id"`
		IsActive        bool   `json:"is_active"`
		Priority        string `json:"priority"`
	}

	UpdateTodo struct {
		ID              uint64 `json:"id"`
		Title           string `json:"title"`
		ActivityGroupID uint64 `json:"activity_group_id"`
		IsActive        bool   `json:"is_active"`
		Priority        string `json:"priority"`
	}
)
