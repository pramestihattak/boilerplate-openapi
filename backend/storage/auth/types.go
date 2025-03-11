package auth_storage

type Register struct {
	UserID            string `db:"user_id"`
	FullName          string `db:"full_name"`
	Email             string `db:"email"`
	Password          string `db:"password"`
	VerificationToken string `db:"verification_token"`
	PhoneNumber       string `db:"phone_number"`
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	UserID   string
	Email    string
	FullName string
	Password string
	Verified bool
}
