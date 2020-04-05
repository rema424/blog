package main

import (
	"blog/sandbox"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/flosch/pongo2"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	nodeEnv = os.Getenv("NODE_ENV")
	port    = os.Getenv("PORT")
)

func main() {
	e := newEcho("")
	e = route(e)

	if port == "" {
		port = "8080"
	}
	// http.Handle("/", e)
	// log.Printf("Listening on localhost:%s", port)
	// if err := http.ListenAndServe(":"+port, nil); err != nil {
	// 	log.Fatal(err)
	// }
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}

func route(e *echo.Echo) *echo.Echo {
	e.GET("/", sandbox.HelloHandler)

	return e
}

func newEcho(tmplDir string) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	// Ref: https://www.ipa.go.jp/security/vuln/websecurity.html
	// Ref: https://www.templarbit.com/blog/jp/2018/07/24/top-http-security-headers-and-how-to-deploy-them/
	// Ref: https://www.slideshare.net/yagihashoo/csp-lv2
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		Skipper:               middleware.DefaultSkipper,
		XSSProtection:         middleware.DefaultSecureConfig.XSSProtection,
		ContentTypeNosniff:    middleware.DefaultSecureConfig.ContentTypeNosniff,
		XFrameOptions:         middleware.DefaultSecureConfig.XFrameOptions,
		HSTSMaxAge:            0,
		HSTSExcludeSubdomains: false,
		ContentSecurityPolicy: "",
		CSPReportOnly:         false,
		HSTSPreloadEnabled:    middleware.DefaultSecureConfig.HSTSPreloadEnabled,
		ReferrerPolicy:        "origin-when-cross-origin",
	}))
	e.Pre(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	e.Debug = nodeEnv == "production"
	e.Validator = newCustomValidator(map[string]validator.Func{
		"isJapaneseZip": IsJapaneseZip,
		"isAlphaSpace":  IsAlphaSpace,
		"isTel":         IsTel,
	})
	e.Renderer = NewPongoRenderer(tmplDir)

	return e
}

type CustomValidator struct {
	validator *validator.Validate
}

func newCustomValidator(m map[string]validator.Func) echo.Validator {
	v := validator.New()
	for key, val := range m {
		v.RegisterValidation(key, val)
	}
	return &CustomValidator{validator: v}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

var reJapaneseZip = regexp.MustCompile(`^[\d]{3}-[\d]{4}$`)

func IsJapaneseZip(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	// go1.11ではロック回避のためにCopy()して使う。go1.12からは必要ない。
	// ref: https://golang.org/doc/go1.12#regexp
	return reJapaneseZip.Copy().Match([]byte(val))
}

var reAlphaSpace = regexp.MustCompile(`^[a-zA-Z\s]+$`)

func IsAlphaSpace(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return reAlphaSpace.Copy().Match([]byte(val))
}

var reTel = regexp.MustCompile(`^[\d]+$`)

func IsTel(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return reTel.Copy().Match([]byte(val))
}

type PongoRenderer struct {
	TmplDirName string
}

func NewPongoRenderer(tmplDirName string) echo.Renderer {
	return &PongoRenderer{tmplDirName}
}

// Errors
var (
	ErrInvalidRenderDataType = errors.New("echoutil: invalid render data type, must be map[string]interface{}")
)

func (r *PongoRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if data == nil {
		data = map[string]interface{}{}
	}
	m, ok := data.(map[string]interface{})
	if !ok {
		return ErrInvalidRenderDataType
	}

	path := filepath.Join(r.TmplDirName, name)
	b, err := pongo2.Must(pongo2.FromCache(path)).ExecuteBytes(m)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
