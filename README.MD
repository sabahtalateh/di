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
err := di.Setup[*SomeType](c, 
    // ..
)
```

#### Setup Init function
Define function to initialize component. Usualy call `NewSomeType` here. To call all the functions defined with `di.Init` or `di.InitE` call `Init` function on container
```go

c := di.NewContainer()

err := di.Setup[*SomeType](c,
    // ..
    di.Init(func(*Container) *SomeType) {
        return NewSomeType()
    }
)

err := c.Init()
// ..

```

#### Setup Init with error function
Define function to initialize component that may return error
```go

func NewSomeType() (*SomeType, error) {
    // ..
}

err := di.Setup[*SomeType](c,
    // ..
    di.InitE(func(*Container) (*SomeType, error)) {
        return NewSomeType()
    }
)
```

#### Setup Name
Define component name. Name may be used with `di.Get/di.GetE`
```go
err := di.Setup[*SomeType](c,
    di.Name("SomeName"),
    // ..
)
```

#### Setup Stage
Define function that will be executed when component executing some stage. Can be used to open/close connections, start/stop background workers, fill caches before app ready to work. To execute stage functions use `ExecStage`. Stage functions will be run in parallel
```go
c := NewContainer()

err := di.Setup[*SomeType](c,
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

#### Get component from container
Component can be retrieved from container during initialization or after it. To get component during initialization use `di.Get`, if component will not be found and panic occures while initialization it will be captured within `Init` function. To get component after initialization use `di.GetE`

```go

c := NewContainer()

err := di.Setup[*SomeType](c,
    // ..
    di.InitE(func(c *Container) (*SomeType, error)) {
        return NewSomeType(
            di.Get[config.AppConfig](c),
            di.Get[*SomeOtherType](c, di.Name("SomeOtherType")),
        )
    }
)
// ..

err = c.Init()
// ..

err = c.Stage("start", ctx)
// ..

someService, err := di.GetE[*SomeType](c)
// ..
someService.DoWork()

```

## Setup and Initialization
Call of `di.Setup` adds component init function to internal map. Order of `di.Setup` calls does not matter, all the init function will be called on container init stage. If component not initialized yet, it's init function will be called. If component init function has `di.Get` function inside, component used with this `di.Get` will be initialized too (if not yet) and will not be initialized again
