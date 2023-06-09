package pkg

type CustomError struct {
	Message    string `json:"Message"`
	StatusCode int    `json:"-"`
}

func (err CustomError) Error() string {
	return err.Message
}
