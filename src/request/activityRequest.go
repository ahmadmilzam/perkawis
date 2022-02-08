package request

type CreateActivity struct {
	Title string `json:"title"`
	Email string `json:"email"`
}

type UpdateActivity struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Email string `json:"email"`
}
