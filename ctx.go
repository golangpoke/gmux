package gmux

import (
	"encoding/json"
	"io"
	"net/http"
)

type Map map[string]interface{}

type Ctx struct {
	w http.ResponseWriter
	r *http.Request
}

func newContext(w http.ResponseWriter, r *http.Request) *Ctx {
	return &Ctx{w: w, r: r}
}

func (c *Ctx) SetHeader(key string, value string) {
	c.w.Header().Set(key, value)
}

func (c *Ctx) SetStatus(code int) {
	c.w.WriteHeader(code)
}

func (c *Ctx) JSON(status int, data interface{}) {
	c.w.WriteHeader(status)
	c.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if encoder := json.NewEncoder(c.w); encoder != nil {
		if err := encoder.Encode(data); err != nil {
			c.Error(err)
		}
	}
}

func (c *Ctx) String(status int, data string) {
	c.w.WriteHeader(status)
	c.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := c.w.Write([]byte(data))
	if err != nil {
		c.Error(err)
	}
}

func (c *Ctx) Bind(data interface{}) error {
	body, err := io.ReadAll(c.r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, data)
}

func (c *Ctx) URL() string {
	return c.r.URL.Path
}

func (c *Ctx) Method() string {
	return c.r.Method
}

func (c *Ctx) PathValue(name string) string {
	return c.r.PathValue(name)
}

func (c *Ctx) Error(err error) {
	http.Error(c.w, err.Error(), http.StatusInternalServerError)
}
