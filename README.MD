# DI Container

```sh
go get github.com/sabahtalateh/di
```

See example service at [https://github.com/sabahtalateh/diexample](https://github.com/sabahtalateh/diexample)

## API

List of function you can use to interact with container

### Setup component

Use `di.Setup` function to setup component

```go
c := di.NewContainer()
di.Setup[*SomeType](c, 
    // ..
)
```

#### Setup Init function
Define function to initialize component. Usualy call `NewSomeType` here
```go
di.Setup[*SomeType](c,
    // ..
    di.Init(func(*Container) *SomeType)
)
```

#### Setup Name function
Define component name. Name may be used with `di.Get/di.GetE`
```go
di.Setup[*SomeType](c,
    di.Name("SomeName"),
    // ..
)

