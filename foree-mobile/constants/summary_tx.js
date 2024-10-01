export const TxSummaryStatusInitial      = "Initial"
export const TxSummaryStatusAwaitPayment = "Await Payment"
export const TxSummaryStatusInProgress   = "In Progress"
export const TxSummaryStatusCompleted    = "Completed"
export const TxSummaryStatusCancelled    = "Cancelled"
export const TxSummaryStatusPickup       = "Ready To Pickup"
export const TxSummaryStatusRefunding    = "Refunding"

export const SummaryTxStatuses = {
  [TxSummaryStatusInitial]: {
    label: TxSummaryStatusInitial,
    borderColor: "border-purple-800",
    textColor: "text-purple-800",
    bgColor: "bg-purple-200",
    description: "We are intiating your transaction. Please wait for email to make the payment."
  },
  [TxSummaryStatusAwaitPayment]: {
    label: TxSummaryStatusAwaitPayment,
    borderColor: "border-yellow-800",
    textColor: "text-yellow-800",
    bgColor: "bg-yellow-200",
    description: "We're waiting for Interac funds. This process usually taks minutes, however your bank or Interac might hold the funds loger for verification."
  },
  [TxSummaryStatusInProgress]: {
    label: TxSummaryStatusInProgress,
    borderColor: "border-purple-800",
    textColor: "text-purple-800",
    bgColor: "bg-purple-200",
    description: "Your funds are on the way. You will be notified when funds have been delivered."
  },
  [TxSummaryStatusCompleted]: {
    label: TxSummaryStatusCompleted,
    borderColor: "border-green-800",
    textColor: "text-green-800",
    bgColor: "bg-green-200",
    description: 'Your funds has been delivered.'
  },
  [TxSummaryStatusCancelled]: {
    label: TxSummaryStatusCancelled,
    borderColor: "border-red-800",
    textColor: "text-red-800",
    bgColor: "bg-red-200",
    description: 'This transaction has been cancelled.'
  },
  [TxSummaryStatusPickup]: {
    label: TxSummaryStatusPickup,
    borderColor: "border-yellow-800",
    textColor: "text-yellow-800",
    bgColor: "bg-yellow-200",
    description: `Your funds are available for collection at any NBP branch counter. The recipient should present the Reference # starting with "NP" along with their CNIC at the time of pick-up.`
  },
  [TxSummaryStatusRefunding]: {
    label: TxSummaryStatusRefunding,
    borderColor: "border-purple-800",
    textColor: "text-purple-800",
    bgColor: "bg-purple-200",
    description: "We have received your qurest to cancel this transaction. You will be notified when the transaction has been cancelled."
  },
}