package v1

type User struct {
	ID          string `json:"id,omitempty"`
	Email       string `json:"email"`
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Anonymous   bool   `json:"anonymous"`
	CreatedAt   int64  `json:"createdAt"`
}
