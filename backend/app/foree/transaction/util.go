package transaction

import "fmt"

func GenerateNbpId(prefix string, transactionID int64) string {
	return fmt.Sprintf("%s%012d", prefix, transactionID)
}
