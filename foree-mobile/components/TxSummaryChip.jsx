import { View, Text } from 'react-native'
import React from 'react'

const TxSummaryStatusInitial      = "Initial"
const TxSummaryStatusAwaitPayment = "Await Payment"
const TxSummaryStatusInProgress   = "In Progress"
const TxSummaryStatusCompleted    = "Completed"
const TxSummaryStatusCancelled    = "Cancelled"
const TxSummaryStatusPickup       = "Ready To Pickup"
const TxSummaryStatusRefunding    = "Refunding"

const SummaryTxStatusChipStyles = {
  [TxSummaryStatusInitial]: {
    borderColor: "border-slate-800",
    textColor: "text-slate-800",
    bgColor: "bg-slate-200",
  },
  [TxSummaryStatusAwaitPayment]: {
    borderColor: "border-yellow-800",
    textColor: "text-yellow-800",
    bgColor: "bg-yellow-200",
  },
  [TxSummaryStatusInProgress]: {
    borderColor: "border-purple-800",
    textColor: "text-purple-800",
    bgColor: "bg-purple-200",
  },
  [TxSummaryStatusCompleted]: {
    borderColor: "border-green-800",
    textColor: "text-green-800",
    bgColor: "bg-green-200",
  },
  [TxSummaryStatusCancelled]: {
    borderColor: "border-red-800",
    textColor: "text-red-800",
    bgColor: "bg-red-200",
  },
  [TxSummaryStatusPickup]: {
    borderColor: "border-yellow-800",
    textColor: "text-yellow-800",
    bgColor: "bg-yellow-200",
  },
  [TxSummaryStatusRefunding]: {
    borderColor: "border-yellow-800",
    textColor: "text-yellow-800",
    bgColor: "bg-yellow-200",
  },
}

const TxSummaryChip = ({status=TxSummaryStatusInitial}={status}) => {
	const style = SummaryTxStatusChipStyles[status]
	return (
		<View className={`px-2 py-1 rounded-full border-[1px] ${style.borderColor} ${style.bgColor}`}>
			<Text className={`${style.textColor}`}>{status}</Text>
		</View>
	)
}

export default TxSummaryChip