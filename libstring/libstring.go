// Package libstring provides string related library functions.
package libstring

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"os/user"
	"strings"
)

// ExpandTilde is a convenience function that expands ~ to full path.
func ExpandTilde(path string) string {
	if path == "" {
		return path
	}

	if path[:2] == "~/" {
		usr, err := user.Current()
		if err != nil || usr == nil {
			return path
		}

		if usr.Name == "root" {
			path = strings.Replace(path, "~", "/root", 1)
		} else {
			path = strings.Replace(path, "~", usr.HomeDir, 1)
		}

	}
	return path
}

// ExpandTilde is a convenience function that expands both ~ and $ENV.
func ExpandTildeAndEnv(path string) string {
	path = ExpandTilde(path)
	return os.ExpandEnv(path)
}

// GeneratePassword returns password.
// size determines length of initial seed bytes.
func GeneratePassword(size int) (string, error) {
	// Force minimum size to 32
	if size < 32 {
		size = 32
	}

	rb := make([]byte, size)
	_, err := rand.Read(rb)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(rb), nil
}

// StringInSlice search exact match in a slice of strings.
func StringInSlice(beingSearched string, list []string) bool {
	for _, b := range list {
		if b == beingSearched {
			return true
		}
	}
	return false
}

// Split r.RemoteAddr, return an IP object (or nil if ParseIP fails)
func GetIP(address string) net.IP {
	// Try to parse it
	splitAddress := strings.Split(address, ":")
	if len(splitAddress) == 0 {
		return nil
	}

	// Convert to IP object
	return net.ParseIP(splitAddress[0])
}

// CSVtoJSON parses CSV to JSON
func CSVtoJSON(csvInput string) ([]byte, error) {

	csvReader := csv.NewReader(strings.NewReader(csvInput))
	lineCount := 0
	var headers []string
	var result bytes.Buffer
	var item bytes.Buffer
	result.WriteString("[")

	for {
		// read just one record, but we could ReadAll() as well
		record, err := csvReader.Read()

		if err == io.EOF {
			result.Truncate(int(len(result.String()) - 1))
			result.WriteString("]")
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return []byte(""), err
		}

		if lineCount == 0 {
			headers = record[:]
			lineCount += 1
		} else {
			item.WriteString("{")
			for i := 0; i < len(headers); i++ {
				item.WriteString("\"" + headers[i] + "\": \"" + record[i] + "\"")
				if i == (len(headers) - 1) {
					item.WriteString("}")
				} else {
					item.WriteString(",")
				}
			}
			result.WriteString(item.String() + ",")
			item.Reset()
			lineCount += 1
		}
	}
	return result.Bytes(), nil
}
