// common
// @author xiangqian
// @date 10:58 2022/11/27
package tank

import (
	"encoding/json"
	"github.com/google/uuid"
)

func Uuid() string {
	return uuid.New().String()
}

// 序列化
func Serialize(i interface{}) (string, error) {
	var jsonByte, err = json.Marshal(i)
	if err != nil {
		return "", err
	}

	return string(jsonByte), nil
}

// 反序列化
func Deserialize(jsonStr string, v any) error {
	return json.Unmarshal([]byte(jsonStr), v)
}
