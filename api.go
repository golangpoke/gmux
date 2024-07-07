package gmux

type Api interface{}
type R struct {
	Code  int         `json:"code"` //业务代码
	Error error       `json:"debug"`
	Data  interface{} `json:"data"` //数据
}

func Result(data interface{}, code int, err error) *R {
	return &R{Code: code, Error: err, Data: data}
}

// global map
var resultMaps = map[int]string{}

func NewMap(maps ...map[int]string) {
	for _, m := range maps {
		for k1, v1 := range m {
			resultMaps[k1] = v1
		}
	}
}

func NewApi(code int, message string) map[int]string {
	return map[int]string{code: message}
}
