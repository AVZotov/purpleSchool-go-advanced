package order

type Service struct {
	Repository Repository
}

func NewService(r Repository) *Service {
	return &Service{
		Repository: r,
	}
}
