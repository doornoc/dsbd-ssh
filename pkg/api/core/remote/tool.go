package remote

import (
	"bufio"
	"github.com/google/uuid"
	"os"
)

func CreateCh(uuid2 uuid.UUID) map[uuid.UUID](chan []byte) {
	outCh := make(map[uuid.UUID](chan []byte))
	outCh[uuid2] = make(chan []byte)
	return outCh
}

func CreateCusCancelCh(uuid2 uuid.UUID) map[uuid.UUID](chan struct{}) {
	cusOutCancelCh := make(map[uuid.UUID](chan struct{}))
	cusOutCancelCh[uuid2] = make(chan struct{})
	return cusOutCancelCh
}

// convert to LF
// CR: 13, LF: 10, CRLF: 1310
// CRLF
func convertToLF(data []byte) []byte {
	var convertByte []byte
	var newLineCodeByte []byte
	for i, sliceByte := range data {
		if len(data)-1 == i {
			continue
		}
		// 13,10のどちらかであれば、newLineCodeにいれこむ
		if sliceByte == 13 || sliceByte == 10 {
			newLineCodeByte = append(newLineCodeByte, sliceByte)
			if len(newLineCodeByte) == 2 {
				newLineCodeByte = []byte{}
				continue
			}
		} else {
			convertByte = append(convertByte, sliceByte)
			continue
		}
		// 改行コード置き換え処理
		switch len(newLineCodeByte) {
		case 1:
			switch newLineCodeByte[0] {
			case 10:
				// LF=>LF
				newLineCodeByte = append(newLineCodeByte, 10)
			case 13:
				//CR => LF
				if len(data)-1 == i {
					newLineCodeByte = append(newLineCodeByte, 10)
				} else if data[i] == 10 {
					//CRLF => LF
					newLineCodeByte = append(newLineCodeByte, 10)
				}
			}
		}
	}
	return convertByte
}

func byteEqual(a, b []byte) bool { return string(a) == string(b) }

func InputKeyLines(lines chan string) {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		lines <- s.Text()
	}
}
