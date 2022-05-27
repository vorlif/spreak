// This file is generated by cldrplural/generator/generate.sh; DO NOT EDIT
package cldrplural

import "math"

func init() {

	addRuleSet([]string{"bm", "bo", "dz", "hnj", "id", "ig", "ii", "ja", "jbo", "jv", "jw", "kde", "kea", "km", "ko", "lkt", "lo", "ms", "my", "nqo", "osa", "sah", "ses", "sg", "su", "th", "to", "tpi", "und", "vi", "wo", "yo", "yue", "zh"}, &RuleSet{
		Categories: newCategories(Other),
		FormFunc: func(ops *Operands) Category {

			return Other
		},
	})

	addRuleSet([]string{"af", "an", "asa", "az", "bal", "bem", "bez", "bg", "brx", "ce", "cgg", "chr", "ckb", "dv", "ee", "el", "eo", "eu", "fo", "fur", "gsw", "ha", "haw", "hu", "jgo", "jmc", "ka", "kaj", "kcg", "kk", "kkj", "kl", "ks", "ksb", "ku", "ky", "lb", "lg", "mas", "mgo", "ml", "mn", "mr", "nah", "nb", "nd", "ne", "nn", "nnh", "no", "nr", "ny", "nyn", "om", "or", "os", "pap", "ps", "rm", "rof", "rwk", "saq", "sd", "sdh", "seh", "sn", "so", "sq", "ss", "ssy", "st", "syr", "ta", "te", "teo", "tig", "tk", "tn", "tr", "ts", "ug", "uz", "ve", "vo", "vun", "wae", "xh", "xog"}, &RuleSet{
		Categories: newCategories(One, Other),
		FormFunc: func(ops *Operands) Category {

			// n = 1
			if ops.N == 1 {
				return One
			}

			return Other
		},
	})

	addRuleSet([]string{"ak", "bho", "guw", "ln", "mg", "nso", "pa", "ti", "wa"}, &RuleSet{
		Categories: newCategories(One, Other),
		FormFunc: func(ops *Operands) Category {

			// n = 0..1
			if isFloatInRange(ops.N, 0, 1) {
				return One
			}

			return Other
		},
	})

	addRuleSet([]string{"am", "as", "bn", "doi", "fa", "gu", "hi", "kn", "pcm", "zu"}, &RuleSet{
		Categories: newCategories(One, Other),
		FormFunc: func(ops *Operands) Category {

			// i = 0 or n = 1
			if ops.I == 0 || ops.N == 1 {
				return One
			}

			return Other
		},
	})

	addRuleSet([]string{"ast", "ca", "de", "en", "et", "fi", "fy", "gl", "ia", "io", "lij", "nl", "sc", "scn", "sv", "sw", "ur", "yi"}, &RuleSet{
		Categories: newCategories(One, Other),
		FormFunc: func(ops *Operands) Category {

			// i = 1 and v = 0
			if ops.I == 1 && ops.V == 0 {
				return One
			}

			return Other
		},
	})

	addRuleSet([]string{"ceb", "fil", "tl"}, &RuleSet{
		Categories: newCategories(One, Other),
		FormFunc: func(ops *Operands) Category {

			// v = 0 and i = 1,2,3 or v = 0 and i % 10 != 4,6,9 or v != 0 and f % 10 != 4,6,9
			if ops.V == 0 && (isIntOneOf(ops.I, 1, 2, 3)) || ops.V == 0 && !isIntOneOf(ops.I%10, 4, 6, 9) || ops.V != 0 && !isIntOneOf(ops.F%10, 4, 6, 9) {
				return One
			}

			return Other
		},
	})

	addRuleSet([]string{"da"}, &RuleSet{
		Categories: newCategories(One, Other),
		FormFunc: func(ops *Operands) Category {

			// n = 1 or t != 0 and i = 0,1
			if ops.N == 1 || ops.T != 0 && (isIntOneOf(ops.I, 0, 1)) {
				return One
			}

			return Other
		},
	})

	addRuleSet([]string{"ff", "hy", "kab"}, &RuleSet{
		Categories: newCategories(One, Other),
		FormFunc: func(ops *Operands) Category {

			// i = 0,1
			if isIntOneOf(ops.I, 0, 1) {
				return One
			}

			return Other
		},
	})

	addRuleSet([]string{"is"}, &RuleSet{
		Categories: newCategories(One, Other),
		FormFunc: func(ops *Operands) Category {

			// t = 0 and i % 10 = 1 and i % 100 != 11 or t != 0
			if ops.T == 0 && ops.I%10 == 1 && ops.I%100 != 11 || ops.T != 0 {
				return One
			}

			return Other
		},
	})

	addRuleSet([]string{"mk"}, &RuleSet{
		Categories: newCategories(One, Other),
		FormFunc: func(ops *Operands) Category {

			// v = 0 and i % 10 = 1 and i % 100 != 11 or f % 10 = 1 and f % 100 != 11
			if ops.V == 0 && ops.I%10 == 1 && ops.I%100 != 11 || ops.F%10 == 1 && ops.F%100 != 11 {
				return One
			}

			return Other
		},
	})

	addRuleSet([]string{"si"}, &RuleSet{
		Categories: newCategories(One, Other),
		FormFunc: func(ops *Operands) Category {

			// n = 0,1 or i = 0 and f = 1
			if (isFloatOneOf(ops.N, 0, 1)) || ops.I == 0 && ops.F == 1 {
				return One
			}

			return Other
		},
	})

	addRuleSet([]string{"tzm"}, &RuleSet{
		Categories: newCategories(One, Other),
		FormFunc: func(ops *Operands) Category {

			// n = 0..1 or n = 11..99
			if isFloatInRange(ops.N, 0, 1) || isFloatInRange(ops.N, 11, 99) {
				return One
			}

			return Other
		},
	})

	addRuleSet([]string{"bs", "hr", "sh", "sr"}, &RuleSet{
		Categories: newCategories(One, Few, Other),
		FormFunc: func(ops *Operands) Category {

			// v = 0 and i % 10 = 1 and i % 100 != 11 or f % 10 = 1 and f % 100 != 11
			if ops.V == 0 && ops.I%10 == 1 && ops.I%100 != 11 || ops.F%10 == 1 && ops.F%100 != 11 {
				return One
			}

			// v = 0 and i % 10 = 2..4 and i % 100 != 12..14 or f % 10 = 2..4 and f % 100 != 12..14
			if ops.V == 0 && isIntInRange(ops.I%10, 2, 4) && !(isIntInRange(ops.I%100, 12, 14)) || isIntInRange(ops.F%10, 2, 4) && !(isIntInRange(ops.F%100, 12, 14)) {
				return Few
			}

			return Other
		},
	})

	addRuleSet([]string{"es"}, &RuleSet{
		Categories: newCategories(One, Many, Other),
		FormFunc: func(ops *Operands) Category {

			// n = 1
			if ops.N == 1 {
				return One
			}

			// e = 0 and i != 0 and i % 1000000 = 0 and v = 0 or e != 0..5
			if ops.C == 0 && ops.I != 0 && ops.I%1000000 == 0 && ops.V == 0 || !(isIntInRange(ops.C, 0, 5)) {
				return Many
			}

			return Other
		},
	})

	addRuleSet([]string{"fr"}, &RuleSet{
		Categories: newCategories(One, Many, Other),
		FormFunc: func(ops *Operands) Category {

			// i = 0,1
			if isIntOneOf(ops.I, 0, 1) {
				return One
			}

			// e = 0 and i != 0 and i % 1000000 = 0 and v = 0 or e != 0..5
			if ops.C == 0 && ops.I != 0 && ops.I%1000000 == 0 && ops.V == 0 || !(isIntInRange(ops.C, 0, 5)) {
				return Many
			}

			return Other
		},
	})

	addRuleSet([]string{"it", "pt-PT"}, &RuleSet{
		Categories: newCategories(One, Many, Other),
		FormFunc: func(ops *Operands) Category {

			// i = 1 and v = 0
			if ops.I == 1 && ops.V == 0 {
				return One
			}

			// e = 0 and i != 0 and i % 1000000 = 0 and v = 0 or e != 0..5
			if ops.C == 0 && ops.I != 0 && ops.I%1000000 == 0 && ops.V == 0 || !(isIntInRange(ops.C, 0, 5)) {
				return Many
			}

			return Other
		},
	})

	addRuleSet([]string{"iu", "naq", "sat", "se", "sma", "smi", "smj", "smn", "sms"}, &RuleSet{
		Categories: newCategories(One, Two, Other),
		FormFunc: func(ops *Operands) Category {

			// n = 1
			if ops.N == 1 {
				return One
			}

			// n = 2
			if ops.N == 2 {
				return Two
			}

			return Other
		},
	})

	addRuleSet([]string{"ksh"}, &RuleSet{
		Categories: newCategories(Zero, One, Other),
		FormFunc: func(ops *Operands) Category {

			// n = 0
			if ops.N == 0 {
				return Zero
			}

			// n = 1
			if ops.N == 1 {
				return One
			}

			return Other
		},
	})

	addRuleSet([]string{"lag"}, &RuleSet{
		Categories: newCategories(Zero, One, Other),
		FormFunc: func(ops *Operands) Category {

			// n = 0
			if ops.N == 0 {
				return Zero
			}

			// i = 0,1 and n != 0
			if (isIntOneOf(ops.I, 0, 1)) && ops.N != 0 {
				return One
			}

			return Other
		},
	})

	addRuleSet([]string{"lv", "prg"}, &RuleSet{
		Categories: newCategories(Zero, One, Other),
		FormFunc: func(ops *Operands) Category {

			// n % 10 = 0 or n % 100 = 11..19 or v = 2 and f % 100 = 11..19
			if math.Mod(ops.N, 10) == 0 || isFloatInRange(math.Mod(ops.N, 100), 11, 19) || ops.V == 2 && isIntInRange(ops.F%100, 11, 19) {
				return Zero
			}

			// n % 10 = 1 and n % 100 != 11 or v = 2 and f % 10 = 1 and f % 100 != 11 or v != 2 and f % 10 = 1
			if math.Mod(ops.N, 10) == 1 && math.Mod(ops.N, 100) != 11 || ops.V == 2 && ops.F%10 == 1 && ops.F%100 != 11 || ops.V != 2 && ops.F%10 == 1 {
				return One
			}

			return Other
		},
	})

	addRuleSet([]string{"mo", "ro"}, &RuleSet{
		Categories: newCategories(One, Few, Other),
		FormFunc: func(ops *Operands) Category {

			// i = 1 and v = 0
			if ops.I == 1 && ops.V == 0 {
				return One
			}

			// v != 0 or n = 0 or n % 100 = 2..19
			if ops.V != 0 || ops.N == 0 || isFloatInRange(math.Mod(ops.N, 100), 2, 19) {
				return Few
			}

			return Other
		},
	})

	addRuleSet([]string{"pt"}, &RuleSet{
		Categories: newCategories(One, Many, Other),
		FormFunc: func(ops *Operands) Category {

			// i = 0..1
			if isIntInRange(ops.I, 0, 1) {
				return One
			}

			// e = 0 and i != 0 and i % 1000000 = 0 and v = 0 or e != 0..5
			if ops.C == 0 && ops.I != 0 && ops.I%1000000 == 0 && ops.V == 0 || !(isIntInRange(ops.C, 0, 5)) {
				return Many
			}

			return Other
		},
	})

	addRuleSet([]string{"shi"}, &RuleSet{
		Categories: newCategories(One, Few, Other),
		FormFunc: func(ops *Operands) Category {

			// i = 0 or n = 1
			if ops.I == 0 || ops.N == 1 {
				return One
			}

			// n = 2..10
			if isFloatInRange(ops.N, 2, 10) {
				return Few
			}

			return Other
		},
	})

	addRuleSet([]string{"be"}, &RuleSet{
		Categories: newCategories(One, Few, Many, Other),
		FormFunc: func(ops *Operands) Category {

			// n % 10 = 1 and n % 100 != 11
			if math.Mod(ops.N, 10) == 1 && math.Mod(ops.N, 100) != 11 {
				return One
			}

			// n % 10 = 2..4 and n % 100 != 12..14
			if isFloatInRange(math.Mod(ops.N, 10), 2, 4) && !(isFloatInRange(math.Mod(ops.N, 100), 12, 14)) {
				return Few
			}

			// n % 10 = 0 or n % 10 = 5..9 or n % 100 = 11..14
			if math.Mod(ops.N, 10) == 0 || isFloatInRange(math.Mod(ops.N, 10), 5, 9) || isFloatInRange(math.Mod(ops.N, 100), 11, 14) {
				return Many
			}

			return Other
		},
	})

	addRuleSet([]string{"cs", "sk"}, &RuleSet{
		Categories: newCategories(One, Few, Many, Other),
		FormFunc: func(ops *Operands) Category {

			// i = 1 and v = 0
			if ops.I == 1 && ops.V == 0 {
				return One
			}

			// i = 2..4 and v = 0
			if isIntInRange(ops.I, 2, 4) && ops.V == 0 {
				return Few
			}

			// v != 0
			if ops.V != 0 {
				return Many
			}

			return Other
		},
	})

	addRuleSet([]string{"dsb", "hsb"}, &RuleSet{
		Categories: newCategories(One, Two, Few, Other),
		FormFunc: func(ops *Operands) Category {

			// v = 0 and i % 100 = 1 or f % 100 = 1
			if ops.V == 0 && ops.I%100 == 1 || ops.F%100 == 1 {
				return One
			}

			// v = 0 and i % 100 = 2 or f % 100 = 2
			if ops.V == 0 && ops.I%100 == 2 || ops.F%100 == 2 {
				return Two
			}

			// v = 0 and i % 100 = 3..4 or f % 100 = 3..4
			if ops.V == 0 && isIntInRange(ops.I%100, 3, 4) || isIntInRange(ops.F%100, 3, 4) {
				return Few
			}

			return Other
		},
	})

	addRuleSet([]string{"gd"}, &RuleSet{
		Categories: newCategories(One, Two, Few, Other),
		FormFunc: func(ops *Operands) Category {

			// n = 1,11
			if isFloatOneOf(ops.N, 1, 11) {
				return One
			}

			// n = 2,12
			if isFloatOneOf(ops.N, 2, 12) {
				return Two
			}

			// n = 3..10,13..19
			if isFloatInRange(ops.N, 3, 10, 13, 19) {
				return Few
			}

			return Other
		},
	})

	addRuleSet([]string{"he"}, &RuleSet{
		Categories: newCategories(One, Two, Many, Other),
		FormFunc: func(ops *Operands) Category {

			// i = 1 and v = 0
			if ops.I == 1 && ops.V == 0 {
				return One
			}

			// i = 2 and v = 0
			if ops.I == 2 && ops.V == 0 {
				return Two
			}

			// v = 0 and n != 0..10 and n % 10 = 0
			if ops.V == 0 && !(isFloatInRange(ops.N, 0, 10)) && math.Mod(ops.N, 10) == 0 {
				return Many
			}

			return Other
		},
	})

	addRuleSet([]string{"lt"}, &RuleSet{
		Categories: newCategories(One, Few, Many, Other),
		FormFunc: func(ops *Operands) Category {

			// n % 10 = 1 and n % 100 != 11..19
			if math.Mod(ops.N, 10) == 1 && !(isFloatInRange(math.Mod(ops.N, 100), 11, 19)) {
				return One
			}

			// n % 10 = 2..9 and n % 100 != 11..19
			if isFloatInRange(math.Mod(ops.N, 10), 2, 9) && !(isFloatInRange(math.Mod(ops.N, 100), 11, 19)) {
				return Few
			}

			// f != 0
			if ops.F != 0 {
				return Many
			}

			return Other
		},
	})

	addRuleSet([]string{"mt"}, &RuleSet{
		Categories: newCategories(One, Few, Many, Other),
		FormFunc: func(ops *Operands) Category {

			// n = 1
			if ops.N == 1 {
				return One
			}

			// n = 0 or n % 100 = 2..10
			if ops.N == 0 || isFloatInRange(math.Mod(ops.N, 100), 2, 10) {
				return Few
			}

			// n % 100 = 11..19
			if isFloatInRange(math.Mod(ops.N, 100), 11, 19) {
				return Many
			}

			return Other
		},
	})

	addRuleSet([]string{"pl"}, &RuleSet{
		Categories: newCategories(One, Few, Many, Other),
		FormFunc: func(ops *Operands) Category {

			// i = 1 and v = 0
			if ops.I == 1 && ops.V == 0 {
				return One
			}

			// v = 0 and i % 10 = 2..4 and i % 100 != 12..14
			if ops.V == 0 && isIntInRange(ops.I%10, 2, 4) && !(isIntInRange(ops.I%100, 12, 14)) {
				return Few
			}

			// v = 0 and i != 1 and i % 10 = 0..1 or v = 0 and i % 10 = 5..9 or v = 0 and i % 100 = 12..14
			if ops.V == 0 && ops.I != 1 && isIntInRange(ops.I%10, 0, 1) || ops.V == 0 && isIntInRange(ops.I%10, 5, 9) || ops.V == 0 && isIntInRange(ops.I%100, 12, 14) {
				return Many
			}

			return Other
		},
	})

	addRuleSet([]string{"ru", "uk"}, &RuleSet{
		Categories: newCategories(One, Few, Many, Other),
		FormFunc: func(ops *Operands) Category {

			// v = 0 and i % 10 = 1 and i % 100 != 11
			if ops.V == 0 && ops.I%10 == 1 && ops.I%100 != 11 {
				return One
			}

			// v = 0 and i % 10 = 2..4 and i % 100 != 12..14
			if ops.V == 0 && isIntInRange(ops.I%10, 2, 4) && !(isIntInRange(ops.I%100, 12, 14)) {
				return Few
			}

			// v = 0 and i % 10 = 0 or v = 0 and i % 10 = 5..9 or v = 0 and i % 100 = 11..14
			if ops.V == 0 && ops.I%10 == 0 || ops.V == 0 && isIntInRange(ops.I%10, 5, 9) || ops.V == 0 && isIntInRange(ops.I%100, 11, 14) {
				return Many
			}

			return Other
		},
	})

	addRuleSet([]string{"sl"}, &RuleSet{
		Categories: newCategories(One, Two, Few, Other),
		FormFunc: func(ops *Operands) Category {

			// v = 0 and i % 100 = 1
			if ops.V == 0 && ops.I%100 == 1 {
				return One
			}

			// v = 0 and i % 100 = 2
			if ops.V == 0 && ops.I%100 == 2 {
				return Two
			}

			// v = 0 and i % 100 = 3..4 or v != 0
			if ops.V == 0 && isIntInRange(ops.I%100, 3, 4) || ops.V != 0 {
				return Few
			}

			return Other
		},
	})

	addRuleSet([]string{"br"}, &RuleSet{
		Categories: newCategories(One, Two, Few, Many, Other),
		FormFunc: func(ops *Operands) Category {

			// n % 10 = 1 and n % 100 != 11,71,91
			if math.Mod(ops.N, 10) == 1 && !isFloatOneOf(math.Mod(ops.N, 100), 11, 71, 91) {
				return One
			}

			// n % 10 = 2 and n % 100 != 12,72,92
			if math.Mod(ops.N, 10) == 2 && !isFloatOneOf(math.Mod(ops.N, 100), 12, 72, 92) {
				return Two
			}

			// n % 10 = 3..4,9 and n % 100 != 10..19,70..79,90..99
			if (isFloatInRange(math.Mod(ops.N, 10), 3, 4) || math.Mod(ops.N, 10) == 9) && !isFloatInRange(math.Mod(ops.N, 100), 10, 19, 70, 79, 90, 99) {
				return Few
			}

			// n != 0 and n % 1000000 = 0
			if ops.N != 0 && math.Mod(ops.N, 1000000) == 0 {
				return Many
			}

			return Other
		},
	})

	addRuleSet([]string{"ga"}, &RuleSet{
		Categories: newCategories(One, Two, Few, Many, Other),
		FormFunc: func(ops *Operands) Category {

			// n = 1
			if ops.N == 1 {
				return One
			}

			// n = 2
			if ops.N == 2 {
				return Two
			}

			// n = 3..6
			if isFloatInRange(ops.N, 3, 6) {
				return Few
			}

			// n = 7..10
			if isFloatInRange(ops.N, 7, 10) {
				return Many
			}

			return Other
		},
	})

	addRuleSet([]string{"gv"}, &RuleSet{
		Categories: newCategories(One, Two, Few, Many, Other),
		FormFunc: func(ops *Operands) Category {

			// v = 0 and i % 10 = 1
			if ops.V == 0 && ops.I%10 == 1 {
				return One
			}

			// v = 0 and i % 10 = 2
			if ops.V == 0 && ops.I%10 == 2 {
				return Two
			}

			// v = 0 and i % 100 = 0,20,40,60,80
			if ops.V == 0 && (isIntOneOf(ops.I%100, 0, 20, 40, 60, 80)) {
				return Few
			}

			// v != 0
			if ops.V != 0 {
				return Many
			}

			return Other
		},
	})

	addRuleSet([]string{"ar", "ars"}, &RuleSet{
		Categories: newCategories(Zero, One, Two, Few, Many, Other),
		FormFunc: func(ops *Operands) Category {

			// n = 0
			if ops.N == 0 {
				return Zero
			}

			// n = 1
			if ops.N == 1 {
				return One
			}

			// n = 2
			if ops.N == 2 {
				return Two
			}

			// n % 100 = 3..10
			if isFloatInRange(math.Mod(ops.N, 100), 3, 10) {
				return Few
			}

			// n % 100 = 11..99
			if isFloatInRange(math.Mod(ops.N, 100), 11, 99) {
				return Many
			}

			return Other
		},
	})

	addRuleSet([]string{"cy"}, &RuleSet{
		Categories: newCategories(Zero, One, Two, Few, Many, Other),
		FormFunc: func(ops *Operands) Category {

			// n = 0
			if ops.N == 0 {
				return Zero
			}

			// n = 1
			if ops.N == 1 {
				return One
			}

			// n = 2
			if ops.N == 2 {
				return Two
			}

			// n = 3
			if ops.N == 3 {
				return Few
			}

			// n = 6
			if ops.N == 6 {
				return Many
			}

			return Other
		},
	})

	addRuleSet([]string{"kw"}, &RuleSet{
		Categories: newCategories(Zero, One, Two, Few, Many, Other),
		FormFunc: func(ops *Operands) Category {

			// n = 0
			if ops.N == 0 {
				return Zero
			}

			// n = 1
			if ops.N == 1 {
				return One
			}

			// n % 100 = 2,22,42,62,82 or n % 1000 = 0 and n % 100000 = 1000..20000,40000,60000,80000 or n != 0 and n % 1000000 = 100000
			if (isFloatOneOf(math.Mod(ops.N, 100), 2, 22, 42, 62, 82)) || math.Mod(ops.N, 1000) == 0 && (isFloatInRange(math.Mod(ops.N, 100000), 1000, 20000) || isFloatOneOf(math.Mod(ops.N, 100000), 40000, 60000, 80000)) || ops.N != 0 && math.Mod(ops.N, 1000000) == 100000 {
				return Two
			}

			// n % 100 = 3,23,43,63,83
			if isFloatOneOf(math.Mod(ops.N, 100), 3, 23, 43, 63, 83) {
				return Few
			}

			// n != 1 and n % 100 = 1,21,41,61,81
			if ops.N != 1 && (isFloatOneOf(math.Mod(ops.N, 100), 1, 21, 41, 61, 81)) {
				return Many
			}

			return Other
		},
	})

}
