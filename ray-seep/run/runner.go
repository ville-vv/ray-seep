package main

type Runner interface {
	Run()
	Scheme() string
	Stop()
}
