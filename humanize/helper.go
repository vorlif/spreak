package humanize

func (h *Humanizer) LanguageName(lang string) string {
	return h.loc.Get(lang)
}
