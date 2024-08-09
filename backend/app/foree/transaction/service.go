package transaction

type TransactionService struct {
	FeeRepo  *FeeRepo
	FeeJoint *FeeJoint
}

//Service clean cache.
// Mothods:
//1. Get Limit
//2. Get remaining limit
//3. create foree trasanction
//4. process foree transaction
//5. cancelTransaction
//6. hardCancelTransaction
//7. create foreee refund transaction
//8. process refund transaction
//9. Quote in memory.// Using map with RWLock now. Need to improve the performance.
