package di

func (c *Container) initComponent(coord coordinate) (any, error) {
	comp, ok := c.components[coord]
	if !ok {
		return nil, ErrNotFound
	}

	if comp.initFn == nil {
		return comp.val, nil
	}

	var (
		err error
	)
	comp.val, err = comp.initFn(c)
	if err != nil {
		return nil, err
	}
	comp.initFn = nil

	c.components[coord] = comp

	return comp.val, nil
}

func (c *Container) enterInit() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return ErrInitialized
	}

	if c.initializing {
		return ErrInitialized
	}

	c.initializing = true

	return nil
}

func (c *Container) exitInit() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.initialized = true
	c.initializing = false
}

func InitE[T any](f func(*Container) (T, error)) withInitE[T] { return f }
func Init[T any](f func(*Container) T) withInit[T]            { return f }

func (c *Container) Init() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = toError(r)
			if !recoverable(err) {
				panic(err)
			}
		}
	}()

	if err = c.enterInit(); err != nil {
		return err
	}

	for coord := range c.components {
		if _, err = c.initComponent(coord); err != nil {
			return err
		}
	}

	c.exitInit()

	return err
}
