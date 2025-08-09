package config

type MarkovConfig struct {
	Order int `envconfig:"MARKOV_ORDER" default:"5"`
}
