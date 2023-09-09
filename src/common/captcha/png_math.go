package captcha

import (
	"math/rand"
	"strconv"
	"time"
)

func gen_captcha_math() (string, string, error) {

	var result int
	var text string

	a := rand.Intn(100)
	b := rand.Intn(10)
	c := rand.Intn(10)

	switch rand.Intn(3) {
	case 1:
		result = a + b*c
		text = strconv.Itoa(a) + " + " + strconv.Itoa(b) + "*" + strconv.Itoa(c)
	case 2:
		result = a + b + c
		text = strconv.Itoa(a) + " + " + strconv.Itoa(b) + " + " + strconv.Itoa(c)
	default:
		result = b*c + a
		text = strconv.Itoa(b) + "*" + strconv.Itoa(c) + " + " + strconv.Itoa(a)
	}

	text_string, err := gen_captcha_png_base64(text)
	if err != nil {
		return "", "", err
	}

	return strconv.Itoa(result), text_string, nil
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
