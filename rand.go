package gutils

import (
    "math/rand"
)

var (
    LowerLetterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
    UpperLetterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
    NumberRunes      = []rune("0123456789")

    LetterRunes          []rune
    LetterAndNumberRules []rune
)

func init() {
    LetterRunes = append(LetterRunes, LowerLetterRunes...)
    LetterRunes = append(LetterRunes, UpperLetterRunes...)

    LetterAndNumberRules = append(LetterAndNumberRules, LetterRunes...)
    LetterAndNumberRules = append(LetterAndNumberRules, NumberRunes...)
}

func RandStringLetter(n int) string {
    return RandStringWithSeed(n, LetterRunes)
}

func RandStringWithSeed(n int, seed []rune) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = seed[rand.Intn(len(seed))]
    }
    return string(b)
}
