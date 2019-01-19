package dto

import "encoding/json"

type Keys struct {
	Keys []map[string]interface{} `json:"keys"`
}

func KeysNew(inbound []byte) (k *Keys) {

	var raw map[string]interface{}
	k = new(Keys)
	json.Unmarshal(inbound, &raw)
	k.Keys = append(k.Keys, raw)

	return
}

type CustomUserData struct {
	Name string
}
