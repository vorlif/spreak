package goextractors

/*
func newTestCache(t *testing.T, dir string) *packageCache {
	testCfg := newTestConfig(dir)
	pkgs, err := loadPackagesForAllDirectories(dir)
	require.NoError(t, err)

	cache := newPackageCache(testCfg, pkgs)
	require.NotNil(t, cache)
	return cache
}

func Test_definitionExtractor_extractStruct(t *testing.T) {
	cache := newTestCache(t, testdataUnit)
	defs := extractDefinitions(cache)

	structName := "github.com/vorlif/xspreakunit.TranslationStruct"
	if assert.Contains(t, defs, structName) {
		definitions := defs[structName]
		assert.Contains(t, definitions, "Singular")
		assert.Contains(t, definitions, "plural")
		assert.Contains(t, definitions, "Domain")
		assert.Contains(t, definitions, "context")
	}
}

func Test_definitionExtractor_extractVar(t *testing.T) {
	cache := newTestCache(t, testdataUnit)
	defs := extractDefinitions(cache)

	funcName := "github.com/vorlif/xspreakunit.TranslateFunc"
	require.Contains(t, defs, funcName)
	definitions := defs[funcName]
	assert.Contains(t, definitions, "one")
	assert.Contains(t, definitions, "tow")
	assert.Contains(t, definitions, "three")
	assert.Contains(t, definitions, "plural")
	assert.Contains(t, definitions, "Domain")
	assert.Contains(t, definitions, "context")
}

func Test_definitionExtractor_extractVar1(t *testing.T) {
	cache := newTestCache(t, testdataUnit)
	defs := extractDefinitions(cache)

	for _, name := range []string{"one", "two", "three"} {
		varName := "github.com/vorlif/xspreakunit." + name
		require.Contains(t, defs, varName)
		assert.NotNil(t, defs[varName][""])
	}
}
*/
