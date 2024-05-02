# DI Container

```sh
go get github.com/sabahtalateh/di
```

See example service at [https://github.com/sabahtalateh/diexample](https://github.com/sabahtalateh/diexample)

## API

List of function you can use to interact with container

### Setup component

Use `di.Setup` function to setup component. If component `A` depends on component `B` it should be setup in corresponding order

```go
c := di.NewContainer()

err := di.Setup[*A](c, 
    di.Init(func(*Container) *A) {
        return NewA()
    }
)

err := di.Setup[*B](c, 
    di.Init(func(*Container) *B) {
        return NewB(di.Get[*A](c))
    }
)
```

#### Setup Init function

Define function to initialize component. Usualy call `NewSomeType` here. To call all the functions defined with `di.Init` or `di.InitE` call `Init` function on container. Functions will be called in order
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

Define function that will be executed when application expiriencing some stage. Can be used to open/close connections, start/stop background workers, fill caches before app started/stopped. To execute stage functions use `ExecStage`. Functions defined on same stage will be executed in parallel

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

Component can be retrieved from container during initialization and after it. To get component during initialization use `di.Get`, if component not found panic occures while initialization that will be captured within `Init` function. To get component after initialization use `di.GetE`

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

### Init Container

After all the components set call `Init`. It will call all the components init functions in order. If component `A` depends on component `B` it should be setup in corresponding order

```go
c := di.NewContainer()
// ..
err := c.Init()
// ..
```

### Execute Stage
Execute stage defined with `Init` function with `ExecStage`
```go
c := di.NewContainer()
// ..
err := c.Init()
// ..

err = c.ExecStage("start")
// ..

err = c.ExecStage("stop")
// ..
```

## Setup and Initialization
Call of `di.Setup` adds component init function to internal initialization list. Order of `di.Setup` calls does matters, all the init function will be called on container init stage in order corresponding setup functions were called
