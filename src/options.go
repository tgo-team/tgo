package tgo

type Options struct {

}

func NewOptions() *Options {

	return &Options{}
}

type Option func(*Options) error