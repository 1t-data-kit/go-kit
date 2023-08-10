package excel_bak

func sliceContainString(list []string, sub string) bool {
	for _, item := range list {
		if item == sub {
			return true
		}
	}
	return false
}
