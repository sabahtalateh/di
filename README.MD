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
    di.Init(func(*Container) *SomeType) {
        return NewSomeType()
    }
)
```

#### Setup Name
Define component name. Name may be used with `di.Get/di.GetE`
```go
di.Setup[*SomeType](c,
    di.Name("SomeName"),
    // ..
)
```

#### Setup Stage
Define function that will be executed when component executing some stage. Can be used to open/close connections, start/stop background workers, fill caches before app ready to work
```go
c := NewContainer()

di.Setup[*SomeType](c,
    // ..
    di.Stage("start", func(ctx context.Context, db *DB) error {
        return db.Open(ctx)
    }),
    di.Stage("stop", func(ctx context.Context, db *DB) error {
        return db.Close(ctx)
    }),
)

// ..
err := c.ExecStage("start", ctx)
// ..
err := c.ExecStage("stop", ctx)
```

