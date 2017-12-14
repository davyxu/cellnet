package util

import "io"

// 完整发送所有封包
func WriteFull(writer io.Writer, p []byte) error {

	total := len(p)

	for pos := 0; pos < total; {

		n, err := writer.Write(p[pos:])

		if err != nil {
			return err
		}

		pos += n
	}

	return nil

}
