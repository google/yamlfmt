package hotfix

func StripCRBytes(crlfContent []byte) []byte {
	onlyLf := []byte{}
	for _, b := range crlfContent {
		if b != '\r' {
			onlyLf = append(onlyLf, b)
		}
	}
	return onlyLf
}

func WriteCRLFBytes(lfContent []byte) []byte {
	crlfContent := []byte{}
	for _, b := range lfContent {
		if b == '\n' {
			crlfContent = append(crlfContent, '\r')
		}
		crlfContent = append(crlfContent, b)
	}
	return crlfContent
}
