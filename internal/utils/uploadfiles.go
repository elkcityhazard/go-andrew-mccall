package utils

import "crypto/rand"

type Tools struct{}

// the random string source

const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+"

func (t *Tools) RandomString(n int) string {
	// define two variables
	//	s is a slice of rune of n length
	//  r uses the random string source to seed a slice of runes
	s, r := make([]rune, n), []rune(randomStringSource)

	// range over s which is a fixed length slice
	// use rand crypto package Prime method
	// generate a random number with x,y values where
	// x gets converted from a big int to a Uint64, and y is normalized to uint64 from the length of the random string source slice
	//	finally update each index of the slice of runes with the index of r based on the modulus of x%y

	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}

	// return the random string

	return string(s)
}
