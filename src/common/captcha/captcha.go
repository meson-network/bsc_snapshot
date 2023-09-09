package captcha

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"strings"
	"time"

	"github.com/coreservice-io/utils/rand_util"
	goredis "github.com/go-redis/redis/v8"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/plugin/redis_plugin"
)

const redis_captcha_prefix = "captcha"

func GenCaptcha() (string, string, error) {
	math_result_str, math_png_base64_str, err := gen_captcha_math()
	if err != nil {
		return "", "", err
	}
	// gen id
	id := rand_util.GenRandStr(24)
	err = set(id, math_result_str)
	if err != nil {
		return "", "", err
	}
	return id, math_png_base64_str, nil
}

func VerifyCaptcha(id, captchaCode string) bool {
	if id == "" || captchaCode == "" {
		return false
	}
	if verify(id, captchaCode, true) {
		return true
	} else {
		return false
	}
}

func set(id string, value string) error {
	key := redis_plugin.GetInstance().GenKey(redis_captcha_prefix, id)
	err := redis_plugin.GetInstance().Set(context.Background(), key, value, time.Minute*5).Err()
	if err != nil {
		basic.Logger.Errorln("captcha RedisStore Set error", "err", err, "id", id, "value", value)
		return err
	}
	return nil
}

// get a capt
func get(id string, clear bool) string {
	key := redis_plugin.GetInstance().GenKey(redis_captcha_prefix, id)
	val, err := redis_plugin.GetInstance().Get(context.Background(), key).Result()
	if err != nil {
		if err != goredis.Nil {
			basic.Logger.Errorln("captcha RedisStore Get error", "err", err, "id", id)
		}
		return ""
	}
	if clear {
		err := redis_plugin.GetInstance().Del(context.Background(), key).Err()
		if err != nil {
			basic.Logger.Errorln("captcha RedisStore Del error", "err", err, "id", id)
		}
	}
	return val
}

// verify a capt
func verify(id, answer string, clear bool) bool {
	v := get(id, clear)
	v = strings.ToLower(v)
	answer = strings.ToLower(answer)
	return v == answer
}

// /////////////////////////////////////////////////////

var font *truetype.Font

func gen_captcha_png_base64(text string) (string, error) {

	if font == nil {
		fontBytes, err := ioutil.ReadFile(basic.AbsPath("/assets/fonts/sans.ttf"))
		if err != nil {
			return "", err
		}
		font, err = freetype.ParseFont(fontBytes)
		if err != nil {
			return "", err
		}
	}

	dpi := 200.0
	size := 12.0

	// Initialize the context.
	fg, bg := image.NewUniform(color.RGBA{127, 128, 128, 255}), image.Transparent
	rgba := image.NewRGBA(image.Rect(0, 0, 220, 50))
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.DrawString(text, freetype.Pt(10, 5+int(c.PointToFixed(size)>>6)))

	buf := new(bytes.Buffer)
	if encode_err := png.Encode(buf, rgba); encode_err != nil {
		return "", encode_err
	}

	base64Encoding := "data:image/png;base64,"
	base64Encoding += base64.StdEncoding.EncodeToString(buf.Bytes())

	return base64Encoding, nil

}
