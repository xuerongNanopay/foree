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
    borderColor: "border-slate-600",
    textColor: "text-slate-600",
    bgColor: "bg-slate-300",
  },
  [TxSummaryStatusAwaitPayment]: {
    borderColor: "border-yellow-600",
    textColor: "text-yellow-600",
    bgColor: "bg-yellow-100",
  },
  [TxSummaryStatusInProgress]: {
    borderColor: "border-purple-600",
    textColor: "text-purple-600",
    bgColor: "bg-purple-100",
  },
  [TxSummaryStatusCompleted]: {
    borderColor: "border-green-600",
    textColor: "text-green-600",
    bgColor: "bg-green-100",
  },
  [TxSummaryStatusCancelled]: {
    borderColor: "border-red-600",
    textColor: "text-red-600",
    bgColor: "bg-red-100",
  },
  [TxSummaryStatusPickup]: {
    borderColor: "border-yellow-600",
    textColor: "text-yellow-600",
    bgColor: "bg-yellow-100",
  },
  [TxSummaryStatusRefunding]: {
    borderColor: "border-yellow-600",
    textColor: "text-yellow-600",
    bgColor: "bg-yellow-100",
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