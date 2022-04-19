package parser

const (
	operators = "|*."
	Epsilon   = "_"
	lower     = "abcdefghijklmnopqrstuvwxyz"
	num       = "0123456789"
	literals  = lower + num + Epsilon
)
