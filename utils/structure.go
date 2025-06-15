package utils

type Response struct {
	Ok  bool        `json:"ok"`
	Msg interface{} `json:"msg"`
}

type LoginBody struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RegisterBody struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Aria     string `json:"aria"`
	Secret   string `json:"secret"`
}

type AriaData struct {
	Aria   string `json:"aria"`
	Secret string `json:"secret"`
}
