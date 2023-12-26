package common

type successRess struct {
	Data interface{} `json:"data"`
}

func NewSuccessResponsePlayer(data interface{}) *successRess {
	return &successRess{Data: data}
}

func SimpleSuccessResponsePlayer(data interface{}) *successRess {
	return NewSuccessResponsePlayer(data)
}
