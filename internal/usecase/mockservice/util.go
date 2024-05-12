package mockservice

import (
	"configurator/pkg/utils"
	"encoding/json"
)

func jsonTransformImpl(jsonStr string, marshaler func(any) ([]byte, error)) (string, error) {
	if jsonStr == "" {
		return "", nil
	}
	var parsedObj any
	if err := json.Unmarshal(utils.S2B(jsonStr), &parsedObj); err != nil {
		return "", err
	}
	marshaled, err := marshaler(&parsedObj)
	if err != nil {
		return "", err
	}
	return utils.B2S(marshaled), nil
}

func beautifyJson(jsonStr string) (string, error) {
	return jsonTransformImpl(jsonStr, func(a any) ([]byte, error) {
		return json.MarshalIndent(a, "", "\t")
	})
}

func compressJsonBeforeStore(jsonStr string) (string, error) {
	return jsonTransformImpl(jsonStr, json.Marshal)
}
