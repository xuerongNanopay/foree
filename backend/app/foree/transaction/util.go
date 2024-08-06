package transaction

import "fmt"

func generateNbpId(prefix string, transactionID int64) string {
	return fmt.Sprintf("%s%012d", prefix, transactionID)
}
