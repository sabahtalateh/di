# DI Container

```sh
go get github.com/sabahtalateh/di
```

See example service at [https://github.com/sabahtalateh/diexample](https://github.com/sabahtalateh/diexample)

## Setup

Use `di.Setup` function to setup component

```go
c := di.NewContainer()
di.Setup[*SomeType](c, 
    di.Init(func(*Container) *SomeType)
)
```
