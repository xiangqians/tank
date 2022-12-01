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

// JSON序列化
func Serialize(v interface{}) ([]byte, error) {
	var buf, err = json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// JSON反序列化
func Deserialize(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
