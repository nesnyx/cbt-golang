package helper

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Answers struct {
	Data AnswersResponse `json:"data"`
}

type AnswersResponse struct {
	QuestionId string `json:"question_id"`
	Answer     string `json:"answer"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}
	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

func PesertaAnswersResponse(questionId string, answer string) Answers {
	data := AnswersResponse{
		QuestionId: questionId,
		Answer:     answer,
	}
	jsonResponse := Answers{Data: data}

	return jsonResponse

}

//func FormatValidationError(err error) []string {
//	var errors []string
//
//	for _, e := range err.(validator.ValidationErrors) {
//		errors = append(errors, e.Error())
//	}
//
//	return errors
//
//}
