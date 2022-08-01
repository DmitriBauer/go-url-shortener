package urlrep

type inFileSplitter struct {
	scanLinesOffset int
}

func newInFileSplitter() *inFileSplitter {
	return &inFileSplitter{
		scanLinesOffset: 1,
	}
}

func (s *inFileSplitter) scanLinesInReverse(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 {
		return 0, nil, nil
	}

	l := len(data)

	for i := l - s.scanLinesOffset - 2; i >= 0; i-- {
		if data[i] == '\n' {
			d := data[i+1 : l-s.scanLinesOffset]
			s.scanLinesOffset = l - i
			return 0, d, nil
		}
	}

	d := data[0 : l-s.scanLinesOffset]
	return len(data), d, nil
}
