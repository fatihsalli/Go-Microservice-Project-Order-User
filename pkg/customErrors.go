package pkg

type InternalServerError struct {
	Message    string
	StatusCode int
}

func (err InternalServerError) Error() string {
	return err.Message
}

type NotFoundError struct {
	Message    string
	StatusCode int
}

func (err NotFoundError) Error() string {
	return err.Message
}

type BadRequestError struct {
	Message    string
	StatusCode int
}

func (err BadRequestError) Error() string {
	return err.Message
}

type ClientSideError struct {
	Message    string
	StatusCode int
}

func (err ClientSideError) Error() string {
	return err.Message
}
