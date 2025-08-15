package main

import (
	"errors"
	"math/big"
)

type BigMathResult struct {
	addition       *big.Int
	subtraction    *big.Int
	multiplication *big.Int
	division       *big.Int
}

func performBigMath(a string, b string) (*BigMathResult, error) {
	// Parse a
	bigA, ok := new(big.Int).SetString(a, 10)
	if !ok {
		return nil, errors.New("could not parse variable a")
	}
	// Parse b
	bigB, ok := new(big.Int).SetString(b, 10)
	if !ok {
		return nil, errors.New("could not parse variable b")
	}
	// Check for division by zero
	if bigB.Sign() == 0 {
		return nil, errors.New("division by zero")
	}

	return &BigMathResult{
		addition:       new(big.Int).Add(bigA, bigB),
		subtraction:    new(big.Int).Sub(bigA, bigB),
		multiplication: new(big.Int).Mul(bigA, bigB),
		division:       new(big.Int).Div(bigA, bigB),
	}, nil
}

func main() {
	A := "987654321098765432109876543210"
	B := "123456789012345678901234567890"

	result, err := performBigMath(A, B)
	if err != nil {
		println(err.Error())
		return
	}

	println("Result:")
	println(result.addition.String())
	println(result.subtraction.String())
	println(result.multiplication.String())
	println(result.division.String())
}
