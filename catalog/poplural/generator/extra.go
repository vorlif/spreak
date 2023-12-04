package main

func getExtraRules() []*ruleData {
	list := make([]*ruleData, len(extraRulesTable))
	for i, rawRule := range extraRulesTable {
		list[i] = newRuleData(rawRule)
	}
	return list
}

// Separate list of rules that are not in the JSON file but are frequently used.
var extraRulesTable = []string{
	"nplurals=3; plural=n == 0 ? 0 : (n == 0 || n == 1) && n != 0 ? 1 : 2;",
	"nplurals=4; plural=n % 100 == 1 ? 0 : n % 100 == 2 ? 1 : n % 100 >= 3 && n % 100 <= 4 ? 2 : 3;",
	"nplurals=2; plural=n % 10 == 1 && n % 100 != 11;",
	"nplurals=2; plural=((n == 1 || (n == 2 || n == 3)) || n % 10 != 4 && n % 10 != 6 && n % 10 != 9);",
	"nplurals=2; plural=(n == 0 || n == 1);",
	"nplurals=5; plural=n % 10 == 1 && n % 100 != 11 && n % 100 != 71 && n % 100 != 91 ? 0 : n % 10 == 2 && n % 100 != 12 && n % 100 != 72 && n % 100 != 92 ? 1 : (n % 10 >= 3 && n % 10 <= 4 || n % 10 == 9) && (n % 100 < 10 || n % 100 > 19) && (n % 100 < 70 || n % 100 > 79) && (n % 100 < 90 || n % 100 > 99) ? 2 : n != 0 && n % 1000000 == 0 ? 3 : 4;",
	"nplurals=2; plural=(n <= 1 || n >= 11 && n <= 99);",
}
