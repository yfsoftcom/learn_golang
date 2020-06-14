package main

type Storage interface {
	Shorten(url string, exp int64) (string, error)
	UnShorten(shortlink string) (string, error)
	Detail(shortlink string) (interface{}, error)
}
