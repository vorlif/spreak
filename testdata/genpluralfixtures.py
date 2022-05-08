#!/usr/bin/python

import json
from gettext import c2py


PLURAL_FORMS = [
    "nplurals=3; plural=(n%10==1 && (n%100<11 || n%100>19) ? 0 : n%10>=2 && n%10<=9 && (n%100<11 || n%100>19) ? 1 : 2);",
    "nplurals=5; plural=(n%10==1 && n%100!=11 && n%100!=71 && n%100!=91 ? 0 : n%10==2 && n%100!=12 && n%100!=72 && n%100!=92 ? 1 : ((n%10>=3 && n%10<=4) || n%10==9) && (n%100<10 || n%100>19) && (n%100<70 || n%100>79) && (n%100<90 || n%100>99) ? 2 : n!=0 && n%1000000==0 ? 3 : 4);",
    "nplurals=6; plural=(n==0 ? 0 : n==1 ? 1 : n==2 ? 2 : n==3 ? 3 : n==6 ? 4 : 5);",
    "nplurals=3; plural=(n==0 ? 0 : n==1 ? 1 : 2);",
    "nplurals=2; plural=(n==0 || n==1);",
    "nplurals=1; plural=0;",
    "nplurals=2; plural=(n<=1 || (n>=11 && n<=99));",
    "nplurals=3; plural=(n%10==0 || (n%100>=11 && n%100<=19) ? 0 : n%10==1 && n%100!=11 ? 1 : 2);",
    "nplurals=4; plural=(n==1 ? 0 : n==0 || (n%100>=2 && n%100<=10) ? 1 : n%100>=11 && n%100<=19 ? 2 : 3);",
    "nplurals=6; plural=(n==0 ? 0 : n==1 ? 1 : n==2 ? 2 : n%100>=3 && n%100<=10 ? 3 : n%100>=11 && n%100<=99 ? 4 : 5);",
    "nplurals=3; plural=(n==1 ? 0 : n==0 || (n!=1 && n%100>=1 && n%100<=19) ? 1 : 2);",
    "nplurals=5; plural=(n==1 ? 0 : n==2 ? 1 : n>=3 && n<=6 ? 2 : n>=7 && n<=10 ? 3 : 4);",
    "nplurals=3; plural=(n==1 ? 0 : n>=2 && n<=4 ? 1 : 2);",
    "nplurals=4; plural=(n%100==1 ? 0 : n%100==2 ? 1 : n%100>=3 && n%100<=4 ? 2 : 3);",
    "nplurals=3; plural=(n==1 ? 0 : n%10>=2 && n%10<=4 && (n%100<12 || n%100>14) ? 1 : 2);",
    "nplurals=2; plural=(n != 1);",
    "nplurals=4; plural=(n==1 || n==11 ? 0 : n==2 || n==12 ? 1 : (n>=3 && n<=10) || (n>=13 && n<=19) ? 2 : 3);",
    "nplurals=3; plural=(n==0 || n==1 ? 0 : n>=2 && n<=10 ? 1 : 2);",
    "nplurals=2; plural=(n==1 || n==2 || n==3 || (n%10!=4 && n%10!=6 && n%10!=9));",
    "nplurals=2; plural=(n%10==1 && n%100!=11);",
    "nplurals=4; plural=(n%10==1 ? 0 : n%10==2 ? 1 : n%100==0 || n%100==20 || n%100==40 || n%100==60 || n%100==80 ? 2 : 3);",
    "nplurals=3; plural=(n==1 ? 0 : n==2 ? 1 : 2);",
    "nplurals=2; plural=(n > 1);",
    "nplurals=4; plural=(n==1 ? 0 : n==2 ? 1 : n>10 && n%10==0 ? 2 : 3);",
    "nplurals=3; plural=(n==0 ? 0 : (n==0 || n==1) && n!=0 ? 1 : 2);",
    "nplurals=3; plural=(n%10==1 && n%100!=11 ? 0 : n%10>=2 && n%10<=4 && (n%100<12 || n%100>14) ? 1 : 2);",
]

NUM = 1000


def gen():
    tests = []
    for plural_form in PLURAL_FORMS:
        parts = plural_form.split(";")
        rule = parts[1].removeprefix(" plural=")
        expr = c2py(rule)

        tests.append({
            'pluralform': plural_form,
            'fixture': [expr(n) for n in range(NUM + 1)]
        })
    return json.dumps(tests)


if __name__ == "__main__":
    print(gen())
