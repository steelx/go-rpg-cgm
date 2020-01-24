package dice

import (
	"math/rand"
	"regexp"
	"sort"
	"strconv"
)

var regex = regexp.MustCompile(`(?i)(?P<op>[×*/^v+-])?\s*((?P<num>\d*)d(?P<sides>[f%]|\d+)(?P<explode>!)?(?P<max>[<>]\d{1,4})?(?P<keep>k-?\d{1,3})?|(?P<alt>\d{1,5})(?P<fudge>f)?)( for (?P<for>[^,;]+))?`)

const (
	Add      = "+"
	Subtract = "-"
	Multiply = "*"
	Divide   = "/"
	Max      = "^"
	Min      = "v"
)

type Dice struct {
	Operator string
	Number   int
	Sides    int
	Minimum  int
	Maximum  int
	Keep     int
	Rolls    []int
	Removed  []int
	Explode  bool
	Fudge    bool
	Total    int
	For      string
}

func Create(diceStr string) func() int {
	result := Parse(diceStr)
	return func() int {
		for _, roll := range result {
			roll.Roll()
		}
		total := 0
		for _, val := range result {
			total += val.Total
		}
		return total
	}
}

func Parse(text string) []*Dice {
	var rolls []*Dice
	for _, m := range regex.FindAllStringSubmatch(text, 8) {
		dice := &Dice{
			Operator: Add,
			Number:   2,
			Sides:    6,
		}
		for i, name := range regex.SubexpNames() {
			switch name {
			case "op":
				if m[i] == "×" {
					dice.Operator = "*"
				} else if m[i] != "" {
					dice.Operator = m[i]
				}
			case "explode":
				dice.Explode = m[i] != ""
			case "fudge":
				if m[i] != "" {
					dice.Fudge = true
					dice.Sides = 3
				}
			case "alt":
				if m[i] != "" {
					num, _ := strconv.Atoi(m[i])
					dice.Number = num
					dice.Sides = 1
				}
			case "num":
				num, err := strconv.Atoi(m[i])
				if err != nil {
					num = 1
				}
				if num > 100 {
					num = 100
				}
				dice.Number = num
			case "sides":
				if m[i] == "" {
					dice.Sides = 6
				} else if m[i] == "f" || m[i] == "F" {
					dice.Fudge = true
					dice.Sides = 3
				} else if m[i] == "%" {
					dice.Sides = 100
				} else {
					dice.Sides, _ = strconv.Atoi(m[i])
					if dice.Sides < 1 {
						dice.Sides = 0
					}
					if dice.Sides > 1000000 {
						dice.Sides = 1000000
					}
				}
			case "keep":
				if m[i] == "" {
					break
				}
				dice.Keep, _ = strconv.Atoi(m[i][1:])
				if dice.Keep > dice.Number {
					dice.Keep = dice.Number
				} else if dice.Keep < -dice.Number {
					dice.Keep = -dice.Number
				}
			case "max":
				if m[i] == "" {
					break
				}
				if m[i][0] == '>' {
					dice.Minimum, _ = strconv.Atoi(m[i][1:])
					if dice.Minimum >= dice.Sides {
						dice.Minimum = dice.Sides - 1
					}
				} else {
					dice.Maximum, _ = strconv.Atoi(m[i][1:])
					if dice.Maximum < 2 {
						dice.Maximum = 2
					}
				}
			case "for":
				dice.For = m[i]
			}
		}
		rolls = append(rolls, dice)
	}
	if len(rolls) == 0 {
		rolls = append(rolls, &Dice{
			Operator: Add,
			Number:   2,
			Sides:    6,
		})
	}
	return rolls
}

func (r *Dice) Roll() {
	r.Total = 0
	r.Rolls = []int{}
	if r.Sides == 0 {
		r.Total = 0
		return
	}
	if r.Sides == 1 {
		r.Total = r.Number
		return
	}
	num := r.Number
	for i := 0; i < num; i++ {
		n := rand.Intn(r.Sides) + 1
		if r.Fudge {
			n -= 2
		}
		r.Total += n
		r.Rolls = append(r.Rolls, n)
		if r.Explode && n == r.Sides {
			num++
		}
	}
	if r.Keep != 0 {
		sort.Ints(r.Rolls)
		if r.Keep > 0 {
			split := len(r.Rolls) - r.Keep
			r.Removed = r.Rolls[:split]
			r.Rolls = r.Rolls[split:]
		} else {
			split := -r.Keep
			r.Removed = r.Rolls[split:]
			r.Rolls = r.Rolls[:split]
		}
		r.Total = 0
		for _, n := range r.Rolls {
			r.Total += n
		}
	}
}
