package pkg

type InternalServerError struct {
	Message string
}

type NotFoundError struct {
	Message string
}

type BadRequestError struct {
	Message string
}

type ClientSideError struct {
	Message string
}
