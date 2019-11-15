package bufio

import bio "bufio"

func ReadNBytes(reader *bio.Reader, n int) ([]byte, error) {
	buf := make([]byte, n)
	totalRead := 0
	for totalRead < n {
		l, err := reader.Read(buf)
		if err != nil {
			return nil, err
		}
		totalRead += l
	}
	return buf, nil
}
