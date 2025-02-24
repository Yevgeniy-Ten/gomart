package repository

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

type DuplicateError struct {
}

func (de *DuplicateError) Error() string {
	return "Login already exists"
}

func NewDuplicateError() error {
	return &DuplicateError{}
}

type ShouldBePositiveError struct {
}

func (de *ShouldBePositiveError) Error() string {
	return "Sum should be positive"
}
func NewShouldBePositiveError() error {
	return &ShouldBePositiveError{}
}
