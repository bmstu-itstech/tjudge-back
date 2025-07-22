package uuid

import "github.com/jaevor/go-nanoid"

const alphabet = "0123456789ABCDEFGHIJKLMNOPRQSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// const lengthLong = 16 // 21 would be overkill
const lengthShort = 8

// var Generate = nanoid.MustCustomASCII(alphabet, lengthLong)
var GenerateShort = nanoid.MustCustomASCII(alphabet, lengthShort)
