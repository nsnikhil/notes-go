package document

type Service interface {
	CreateDocument(name, content string) (string, error)
}

type documentService struct {
	store Store
}

func (ds *documentService) CreateDocument(name, content string) (string, error) {
	return "", nil
}

func NewDocumentService(store Store) Service {
	return &documentService{
		store: store,
	}
}
