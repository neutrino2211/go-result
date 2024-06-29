package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/neutrino2211/go-result"
)

func MakeHttpRequest(url string) string {
	res, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	println(res.StatusCode)

	content, _ := io.ReadAll(res.Body)

	return string(content)
}

func ResultHttpRequest(url string) string {
	res := result.SomePair(http.Get(url)).Unwrap()

	println(res.StatusCode)

	content := result.SomePair(io.ReadAll(res.Body)).Unwrap()

	return string(content)
}

func main() {
	googleHtml := result.Try(func() string {
		return MakeHttpRequest("https://google.com")
	}).Or("")

	println(fmt.Sprintf("The google homepage has %d characters", len(googleHtml)))

	shouldFailHtml := result.Try(func() string {
		return MakeHttpRequest("https://notaurl")
	})

	println(fmt.Sprintf("The html should have 0 characters [%d]", len(shouldFailHtml.Or(""))))
	println(fmt.Sprintf("The actual error is: %s", shouldFailHtml.Error()))

	shouldAlsoFailHtml := result.Try(func() string {
		return ResultHttpRequest("https://alsonotaurl")
	})

	println(fmt.Sprintf("The html should also have 0 characters [%d]", len(shouldAlsoFailHtml.Or(""))))
	println(fmt.Sprintf("The actual error is: %s", shouldAlsoFailHtml.Error()))
}
