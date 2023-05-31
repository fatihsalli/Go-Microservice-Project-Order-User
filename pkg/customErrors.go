package pkg

type InternalServerError struct {
	Message    string
	StatusCode int
}

type NotFoundError struct {
	Message    string
	StatusCode int
}

type BadRequestError struct {
	Message    string
	StatusCode int
}

type ClientSideError struct {
	Message    string
	StatusCode int
}
