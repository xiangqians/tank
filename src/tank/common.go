// common
// @author xiangqian
// @date 10:58 2022/11/27
package tank

import (
	"encoding/json"
	"github.com/golang/freetype/truetype"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"golang.org/x/image/font"
	"log"
	"math/rand"
	"reflect"
	"time"
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

// 反射执行方法
// name: 方法名
// in: 入参，如果没有参数可以传 nil 或者空切片 make([]reflect.Value, 0)
func CallMethod(i any, name string, in []reflect.Value) interface{} {
	ref := reflect.ValueOf(i)
	method := ref.MethodByName(name)
	if method.IsValid() {
		r := method.Call(in)
		return r[0].Interface()
	}
	return nil
}

// 反射执行属性
func CallField(i any, name string, in []reflect.Value) []reflect.Value {
	ref := reflect.ValueOf(i)
	field := ref.FieldByName(name)
	if field.IsValid() {
		r := field.Call(in)
		//return r[0].Interface()
		//return r[0].Elem()
		return r
	}
	return nil
}

// 字体表面
// font.Face
func CreateFontFace(size, dpi float64) font.Face {
	tt, err := truetype.Parse(fonts.ArcadeN_ttf)
	if err != nil {
		log.Fatal(err)
	}
	return truetype.NewFace(tt, &truetype.Options{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

func RandIntn(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}
