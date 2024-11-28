package middleware

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

// Unit тест для функции authFiler
//
// Количество кейсов: 5
//
// Тест-кейсы:
// {name: "/ PATH TEST", path: "/", want: false}
// {name: "/test PATH TEST", path: "/test", want: false}
// {name: "/test/testPage PATH TEST", path: "/test/testPage", want: false}
// {name: "/swagger PATH TEST", path: "/swagger", want: true}
// {name: "/swagger/test PATH TEST", path: "/swagger/test", want: true}
func Test_authFilter(t *testing.T) {
	appMock := fiber.New()

	c := appMock.AcquireCtx(&fasthttp.RequestCtx{})
	defer appMock.ReleaseCtx(c)

	caseCount := 0

	cases := []struct {
		name string
		path string
		want bool
	}{
		{name: "/ PATH TEST", path: "/", want: false},
		{name: "/test PATH TEST", path: "/test", want: false},
		{name: "/test/testPage PATH TEST", path: "/test/testPage", want: false},
		{name: "/swagger PATH TEST", path: "/swagger", want: true},
		{name: "/swagger/test PATH TEST", path: "/swagger/test", want: true},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			caseCount++

			path := c.Path(cs.path)

			t.Logf("CASE: %v/%v", caseCount, len(cases))
			t.Logf("NAME: %s", cs.name)
			t.Logf("PATH: %s", path)
			t.Logf("WANT: %v", cs.want)

			got := authFiler(c)

			assert.Equalf(t, cs.want, got, "GOT %v, WANT %v", got, cs.want)
		})
	}
}

// Unit тест для функции authFiler
//
// Количество кейсов: 3
//
// Тест-кейсы:
// {name: "VALID KEY CASE", input: "validKey", want: true, wantErr: nil}
// {name: "WRONG KEY CASE", input: "wrongKey", want: false, wantErr: fiber.ErrUnauthorized}
// {name: "EMPTY KEY CASE", input: "", want: false, wantErr: fiber.ErrUnauthorized}
func Test_keyauthValidator(t *testing.T) {
	os.Setenv("AUF_CITATY_KEY", "validKey")
	defer os.Unsetenv("AUF_CITATY_KEY")

	appMock := fiber.New()

	c := appMock.AcquireCtx(&fasthttp.RequestCtx{})
	defer appMock.ReleaseCtx(c)

	caseCount := 0

	cases := []struct {
		name    string
		input   string
		want    bool
		wantErr error
	}{
		{name: "VALID KEY CASE", input: "validKey", want: true, wantErr: nil},
		{name: "WRONG KEY CASE", input: "wrongKey", want: false, wantErr: fiber.ErrUnauthorized},
		{name: "EMPTY KEY CASE", input: "", want: false, wantErr: fiber.ErrUnauthorized},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			caseCount++

			t.Logf("CASE: %v/%v", caseCount, len(cases))
			t.Logf("NAME: %s", cs.name)
			t.Logf("INPUT: %s", cs.input)
			t.Logf("WANT: %v", cs.want)
			t.Logf("WANT_ERR: %v", cs.wantErr)

			got, gotErr := keyauthValidator(c, cs.input)

			assert.Equalf(t, cs.want, got, "GOT %v, WANT %v", got, cs.want)
			assert.Equalf(t, cs.wantErr, gotErr, "GOT %v, WANT %v", gotErr, cs.wantErr)
		})
	}
}
