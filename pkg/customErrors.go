package pkg

type InternalServerError struct {
	Message    string `json:"Message"`
	StatusCode int    `json:"-"`
}

func (err InternalServerError) Error() string {
	return err.Message
}

type NotFoundError struct {
	Message    string `json:"Message"`
	StatusCode int    `json:"-"`
}

func (err NotFoundError) Error() string {
	return err.Message
}

type BadRequestError struct {
	Message    string `json:"Message"`
	StatusCode int    `json:"-"`
}

func (err BadRequestError) Error() string {
	return err.Message
}

type ClientSideError struct {
	Message    string `json:"Message"`
	StatusCode int    `json:"-"`
}

func (err ClientSideError) Error() string {
	return err.Message
}
