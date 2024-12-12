package middleware

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

// Unit тест для функции AuthFiler
func TestUnitAuthFilter(t *testing.T) {
	cases := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "page that needs auth case",
			path: "/",
			want: false,
		},
		{
			name: "page that doesn't need auth 'swagger' case",
			path: "/swagger",
			want: true,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			mockApp := fiber.New()

			c := mockApp.AcquireCtx(&fasthttp.RequestCtx{})
			defer mockApp.ReleaseCtx(c)

			c.Path(cs.path)
			got := AuthFiler(c)

			assert.Equal(t, cs.want, got)
		})
	}
}

// Unit тест для функции KeyauthValidator
func TestUnitKeyauthValidator(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    bool
		wantErr error
	}{
		{
			name:    "valid key case",
			input:   "valid",
			want:    true,
			wantErr: nil,
		},
		{
			name:    "wrong key case",
			input:   "wrong",
			want:    false,
			wantErr: fiber.ErrUnauthorized,
		},
		{
			name:    "empty key case",
			input:   "",
			want:    false,
			wantErr: fiber.ErrUnauthorized,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			os.Setenv("API_KEY", "valid")
			defer os.Unsetenv("API_KEY")

			mockApp := fiber.New()

			c := mockApp.AcquireCtx(&fasthttp.RequestCtx{})
			defer mockApp.ReleaseCtx(c)

			got, gotErr := KeyauthValidator(c, cs.input)

			assert.Equal(t, cs.want, got)
			assert.Equal(t, cs.wantErr, gotErr)
		})
	}
}
