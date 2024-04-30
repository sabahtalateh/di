package di

type withName string

func Name(n string) withName { return withName(n) }
