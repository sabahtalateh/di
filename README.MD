# DI Container

```sh
go get github.com/sabahtalateh/di
```

## Setup

Use `di.Setup` function to setup component

```go
c := di.NewContainer()
di.Setup[*SomeType](c, 
    di.Init(func(*Container) *SomeType)
)
```
