//go:generate mockery --output=../mocks --name ContactHandler
package handler

type ContactHandler interface {
	List()
	Add()
	Detail()
	Update()
	Delete()
}
