package ba400

// Repository is a repository for HLSeven
type Repository struct {
	*TCP
}

// NewRepository returns a new HLSeven repository
func NewRepository(tcp *TCP) *Repository {
	return &Repository{
		tcp,
	}
}
