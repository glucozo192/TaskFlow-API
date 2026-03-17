package postgres

// LoginUserRepository is a convenience alias used by the HTTP middleware
// to look up users by ID for token validation.
type LoginUserRepository = UserRepository
