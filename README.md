
# go-result

Rust-like error and nil value handling for Go


[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)


## Installation

Install go-option via `go get`

```bash
  go get -uv https://github.com/neutrino2211/go-result
```
    
## Documentation

[Documentation](docs.md)


## Examples

- Converting a normal Go function into a function that uses results

    ```go
    // A normal Go function
    func MakeHttpRequest(url string) string {
        res, err := http.Get(url)

        if err != nil {
            panic(err)
        }

        println(res.StatusCode)

        content, err := io.ReadAll(res.Body)

        if err != nil {
            panic(err)
        }

        return string(content)
    }

    // A function that uses results
    func ResultHttpRequest(url string) string {
        res := result.SomePair(http.Get(url)).Unwrap() // Can panic

        println(res.StatusCode)

        content := result.SomePair(io.ReadAll(res.Body)).Unwrap() // Can panic

        return string(content)
    }
    ```

- Writing a function that returns results

    ```go
    func PositiveAtoi(str string) *result.Result[int] {
        value, err := strconv.Atoi(str)

        if err != nil {
            return result.Err[int](err) // Return the conversion error
        }

        if value < 0 {
            return result.Err[int](errors.New("PositiveAtoi: number provided must be >= 0")) // Custom failure case
        }

        return result.Some(value)
    }
    ```

- Trying functions that can panic

    ```go
    func ResultHttpRequest(url string) string {
        res := result.SomePair(http.Get(url)).Unwrap() // Can panic

        println(res.StatusCode)

        content := result.SomePair(io.ReadAll(res.Body)).Unwrap() // Can panic

        return string(content)
    }

    func main() {
        requestResult := result.Try(func () string {
            return ResultHttpRequest("https://notarealwebsite") // Can panic
        })

        html := requestResult.Or("")

        println(len(html)) // 0
    }
    ```

Check the [examples directory](./examples) for more.