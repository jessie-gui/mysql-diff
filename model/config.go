package model

type Config struct {
	Source Source
	Target Target
}

// Source /**
type Source struct {
	User    string
	Pwd     string
	Address string
	DbBase  string
}

// Target /**
type Target struct {
	User    string
	Pwd     string
	Address string
	DbBase  string
}
