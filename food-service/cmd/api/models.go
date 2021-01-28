package main

import "encoding/json"

type Marshaller interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) (Marshaller, error)
}

type Recipe struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Img  string `json:"img"`
}

func (r *Recipe) Marshal() ([]byte, error) { return json.Marshal(r) }
func (r *Recipe) Unmarshal(data []byte) (Marshaller, error) {
	if err := json.Unmarshal(data, r); err != nil {
		return nil, err
	}
	return r, nil
}
