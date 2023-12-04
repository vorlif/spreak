// This file is generated by cldrplural/generator/main.go; DO NOT EDIT
package poplural

import (
	"fmt"
	"testing"

	"golang.org/x/text/language"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuiltinBmBoDzHnjIdIgIiJaJboJvJwKdeKeaKmKoLktLoMsMyNqoOsaSahSesSgSuThToTpiViWoYoYueZhZh_HansZh_Hant(t *testing.T) {
	for _, lang := range []string{"bm", "bo", "dz", "hnj", "id", "ig", "ii", "ja", "jbo", "jv", "jw", "kde", "kea", "km", "ko", "lkt", "lo", "ms", "my", "nqo", "osa", "sah", "ses", "sg", "su", "th", "to", "tpi", "vi", "wo", "yo", "yue", "zh", "zh_Hans", "zh_Hant"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

	}
}

func TestBuiltinAfAnAsaAstAzBalBemBezBgBrxCeCggChrCkbDaDeDe_ATDe_CHDvEeElEnEn_AUEn_CAEn_GBEn_USEoEtEuFiFoFurFyGlGswHaHawHuIaIoJgoJmcKaKajKcgKkKkjKlKsKsbKuKyLbLgLijMasMgoMlMnMrNahNbNdNeNlNl_BENnNnhNoNrNyNynOmOrOsPapPsRmRofRwkSaqScScnSdSdhSehSnSoSqSsSsyStSvSwSw_CDSyrTaTeTeoTigTkTnTrTsUgUrUzVeVoVunWaeXhXogYi(t *testing.T) {
	for _, lang := range []string{"af", "an", "asa", "ast", "az", "bal", "bem", "bez", "bg", "brx", "ce", "cgg", "chr", "ckb", "da", "de", "de_AT", "de_CH", "dv", "ee", "el", "en", "en_AU", "en_CA", "en_GB", "en_US", "eo", "et", "eu", "fi", "fo", "fur", "fy", "gl", "gsw", "ha", "haw", "hu", "ia", "io", "jgo", "jmc", "ka", "kaj", "kcg", "kk", "kkj", "kl", "ks", "ksb", "ku", "ky", "lb", "lg", "lij", "mas", "mgo", "ml", "mn", "mr", "nah", "nb", "nd", "ne", "nl", "nl_BE", "nn", "nnh", "no", "nr", "ny", "nyn", "om", "or", "os", "pap", "ps", "rm", "rof", "rwk", "saq", "sc", "scn", "sd", "sdh", "seh", "sn", "so", "sq", "ss", "ssy", "st", "sv", "sw", "sw_CD", "syr", "ta", "te", "teo", "tig", "tk", "tn", "tr", "ts", "ug", "ur", "uz", "ve", "vo", "vun", "wae", "xh", "xog", "yi"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"1"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"0", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

	}
}

func TestBuiltinAkAmAsBhoBnDoiFaFa_AFFfGuGuwHiHi_LatnHyKabKnLnMgNsoPaPcmSiTiWaZu(t *testing.T) {
	for _, lang := range []string{"ak", "am", "as", "bho", "bn", "doi", "fa", "fa_AF", "ff", "gu", "guw", "hi", "hi_Latn", "hy", "kab", "kn", "ln", "mg", "nso", "pa", "pcm", "si", "ti", "wa", "zu"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"0", "1"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

	}
}

func TestBuiltinCebFilTl(t *testing.T) {
	for _, lang := range []string{"ceb", "fil", "tl"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"0", "1", "2", "3", "5", "7", "8", "10", "11", "12", "13", "15", "17", "18", "20", "21", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"4", "6", "9", "14", "16", "19", "24", "26", "104", "1004"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

	}
}

func TestBuiltinIsMk(t *testing.T) {
	for _, lang := range []string{"is", "mk"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"1", "21", "31", "41", "51", "61", "71", "81", "101", "1001"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"0", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

	}
}

func TestBuiltinTzm(t *testing.T) {
	for _, lang := range []string{"tzm"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"0", "1", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "100", "101", "102", "103", "104", "105", "106", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

	}
}

func TestBuiltinBeBsHrRuShSrSr_MEUk(t *testing.T) {
	for _, lang := range []string{"be", "bs", "hr", "ru", "sh", "sr", "sr_ME", "uk"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"1", "21", "31", "41", "51", "61", "71", "81", "101", "1001"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"2", "3", "4", "22", "23", "24", "32", "33", "34", "42", "43", "44", "52", "53", "54", "62", "102", "1002"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"0", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

	}
}

func TestBuiltinCaEsEs_419Es_ESEs_MXItPt_PTVec(t *testing.T) {
	for _, lang := range []string{"ca", "es", "es_419", "es_ES", "es_MX", "it", "pt_PT", "vec"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"1"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"1000000", "1e6", "2e6", "3e6", "4e6", "5e6", "6e6"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"0", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "100", "1000", "10000", "100000", "1e3", "2e3", "3e3", "4e3", "5e3", "6e3"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

	}
}

func TestBuiltinCsSk(t *testing.T) {
	for _, lang := range []string{"cs", "sk"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"1"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"2", "3", "4"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"0", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

	}
}

func TestBuiltinFrFr_CAFr_CHPtPt_BR(t *testing.T) {
	for _, lang := range []string{"fr", "fr_CA", "fr_CH", "pt", "pt_BR"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"0", "1"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"1000000", "1e6", "2e6", "3e6", "4e6", "5e6", "6e6"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "100", "1000", "10000", "100000", "1e3", "2e3", "3e3", "4e3", "5e3", "6e3"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

	}
}

func TestBuiltinHeIuNaqSatSeSmaSmiSmjSmnSms(t *testing.T) {
	for _, lang := range []string{"he", "iu", "naq", "sat", "se", "sma", "smi", "smj", "smn", "sms"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"1"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"2"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"0", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

	}
}

func TestBuiltinKshLag(t *testing.T) {
	for _, lang := range []string{"ksh", "lag"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"0"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"1"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

	}
}

func TestBuiltinLt(t *testing.T) {
	for _, lang := range []string{"lt"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"1", "21", "31", "41", "51", "61", "71", "81", "101", "1001"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"2", "3", "4", "5", "6", "7", "8", "9", "22", "23", "24", "25", "26", "27", "28", "29", "102", "1002"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"0", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "30", "40", "50", "60", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

	}
}

func TestBuiltinLvPrg(t *testing.T) {
	for _, lang := range []string{"lv", "prg"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"0", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "30", "40", "50", "60", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"1", "21", "31", "41", "51", "61", "71", "81", "101", "1001"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"2", "3", "4", "5", "6", "7", "8", "9", "22", "23", "24", "25", "26", "27", "28", "29", "102", "1002"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

	}
}

func TestBuiltinMoRoRo_MD(t *testing.T) {
	for _, lang := range []string{"mo", "ro", "ro_MD"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"1"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"0", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "101", "1001"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31", "32", "33", "34", "35", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

	}
}

func TestBuiltinPl(t *testing.T) {
	for _, lang := range []string{"pl"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"1"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"2", "3", "4", "22", "23", "24", "32", "33", "34", "42", "43", "44", "52", "53", "54", "62", "102", "1002"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"0", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

	}
}

func TestBuiltinShi(t *testing.T) {
	for _, lang := range []string{"shi"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"0", "1"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"2", "3", "4", "5", "6", "7", "8", "9", "10"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

	}
}

func TestBuiltinDsbHsbSl(t *testing.T) {
	for _, lang := range []string{"dsb", "hsb", "sl"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"1", "101", "201", "301", "401", "501", "601", "701", "1001"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"2", "102", "202", "302", "402", "502", "602", "702", "1002"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"3", "4", "103", "104", "203", "204", "303", "304", "403", "404", "503", "504", "603", "604", "703", "704", "1003"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

		for _, example := range []string{"0", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 3, form, fmt.Sprintf("rule.Evaluate(%s) should be 3", example))
		}

	}
}

func TestBuiltinGd(t *testing.T) {
	for _, lang := range []string{"gd"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"1", "11"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"2", "12"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"3", "4", "5", "6", "7", "8", "9", "10", "13", "14", "15", "16", "17", "18", "19"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

		for _, example := range []string{"0", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31", "32", "33", "34", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 3, form, fmt.Sprintf("rule.Evaluate(%s) should be 3", example))
		}

	}
}

func TestBuiltinGv(t *testing.T) {
	for _, lang := range []string{"gv"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"1", "11", "21", "31", "41", "51", "61", "71", "101", "1001"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"2", "12", "22", "32", "42", "52", "62", "72", "102", "1002"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"0", "20", "40", "60", "80", "100", "120", "140", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

		for _, example := range []string{"3", "4", "5", "6", "7", "8", "9", "10", "13", "14", "15", "16", "17", "18", "19", "23", "103", "1003"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 3, form, fmt.Sprintf("rule.Evaluate(%s) should be 3", example))
		}

	}
}

func TestBuiltinBr(t *testing.T) {
	for _, lang := range []string{"br"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"1", "21", "31", "41", "51", "61", "81", "101", "1001"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"2", "22", "32", "42", "52", "62", "82", "102", "1002"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"3", "4", "9", "23", "24", "29", "33", "34", "39", "43", "44", "49", "103", "1003"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

		for _, example := range []string{"1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 3, form, fmt.Sprintf("rule.Evaluate(%s) should be 3", example))
		}

		for _, example := range []string{"0", "5", "6", "7", "8", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "100", "1000", "10000", "100000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 4, form, fmt.Sprintf("rule.Evaluate(%s) should be 4", example))
		}

	}
}

func TestBuiltinGa(t *testing.T) {
	for _, lang := range []string{"ga"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"1"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"2"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"3", "4", "5", "6"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

		for _, example := range []string{"7", "8", "9", "10"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 3, form, fmt.Sprintf("rule.Evaluate(%s) should be 3", example))
		}

		for _, example := range []string{"0", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 4, form, fmt.Sprintf("rule.Evaluate(%s) should be 4", example))
		}

	}
}

func TestBuiltinMt(t *testing.T) {
	for _, lang := range []string{"mt"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"1"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"2"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"0", "3", "4", "5", "6", "7", "8", "9", "10", "103", "104", "105", "106", "107", "108", "109", "1003"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

		for _, example := range []string{"11", "12", "13", "14", "15", "16", "17", "18", "19", "111", "112", "113", "114", "115", "116", "117", "1011"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 3, form, fmt.Sprintf("rule.Evaluate(%s) should be 3", example))
		}

		for _, example := range []string{"20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31", "32", "33", "34", "35", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 4, form, fmt.Sprintf("rule.Evaluate(%s) should be 4", example))
		}

	}
}

func TestBuiltinArAr_001Ars(t *testing.T) {
	for _, lang := range []string{"ar", "ar_001", "ars"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"0"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"1"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"2"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

		for _, example := range []string{"3", "4", "5", "6", "7", "8", "9", "10", "103", "104", "105", "106", "107", "108", "109", "110", "1003"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 3, form, fmt.Sprintf("rule.Evaluate(%s) should be 3", example))
		}

		for _, example := range []string{"11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "111", "1011"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 4, form, fmt.Sprintf("rule.Evaluate(%s) should be 4", example))
		}

		for _, example := range []string{"100", "101", "102", "200", "201", "202", "300", "301", "302", "400", "401", "402", "500", "501", "502", "600", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 5, form, fmt.Sprintf("rule.Evaluate(%s) should be 5", example))
		}

	}
}

func TestBuiltinCy(t *testing.T) {
	for _, lang := range []string{"cy"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"0"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"1"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"2"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

		for _, example := range []string{"3"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 3, form, fmt.Sprintf("rule.Evaluate(%s) should be 3", example))
		}

		for _, example := range []string{"6"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 4, form, fmt.Sprintf("rule.Evaluate(%s) should be 4", example))
		}

		for _, example := range []string{"4", "5", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "100", "1000", "10000", "100000", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 5, form, fmt.Sprintf("rule.Evaluate(%s) should be 5", example))
		}

	}
}

func TestBuiltinKw(t *testing.T) {
	for _, lang := range []string{"kw"} {
		rule := getBuiltInForLanguage(language.MustParse(lang).String())
		require.NotNil(t, rule)

		for _, example := range []string{"0"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 0, form, fmt.Sprintf("rule.Evaluate(%s) should be 0", example))
		}

		for _, example := range []string{"1"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 1, form, fmt.Sprintf("rule.Evaluate(%s) should be 1", example))
		}

		for _, example := range []string{"2", "22", "42", "62", "82", "102", "122", "142", "1000", "10000", "100000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 2, form, fmt.Sprintf("rule.Evaluate(%s) should be 2", example))
		}

		for _, example := range []string{"3", "23", "43", "63", "83", "103", "123", "143", "1003"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 3, form, fmt.Sprintf("rule.Evaluate(%s) should be 3", example))
		}

		for _, example := range []string{"21", "41", "61", "81", "101", "121", "141", "161", "1001"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 4, form, fmt.Sprintf("rule.Evaluate(%s) should be 4", example))
		}

		for _, example := range []string{"4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "100", "1004", "1000000"} {
			form, err := rule.Evaluate(example)
			require.NoError(t, err)
			assert.Equal(t, 5, form, fmt.Sprintf("rule.Evaluate(%s) should be 5", example))
		}

	}
}
