package orders

type NotFoundError struct {
	Number string
}

func (ne *NotFoundError) Error() string {
	return "Order not found"
}

func NewNotFoundError(number string) error {
	return &NotFoundError{
		Number: number,
	}
}
