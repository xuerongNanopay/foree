import { View, Text } from 'react-native'
import React from 'react'
import { SummaryTxStatuses, TxSummaryStatusAwaitPayment, TxSummaryStatusCancelled, TxSummaryStatusCompleted, TxSummaryStatusInitial, TxSummaryStatusInProgress, TxSummaryStatusPickup, TxSummaryStatusRefunding } from '../constants/summary_tx'


const TxSummaryChip = ({status=TxSummaryStatusInitial}={status}) => {
	const st = SummaryTxStatuses[status]
  switch (status) {
    case TxSummaryStatusInitial:
      return (
        <View className={`px-2 py-1 rounded-full border-[1px] ${st.borderColor} ${st.bgColor}`}>
          <Text className={`${st.textColor}`}>{st.label}</Text>
        </View>
      )
    case TxSummaryStatusAwaitPayment:
      return (
        <View className={`px-2 py-1 rounded-full border-[1px] ${st.borderColor} ${st.bgColor}`}>
          <Text className={`${st.textColor}`}>{st.label}</Text>
        </View>
      )
    case TxSummaryStatusInitial:
      return (
        <View className={`px-2 py-1 rounded-full border-[1px] ${st.borderColor} ${st.bgColor}`}>
          <Text className={`${st.textColor}`}>{st.label}</Text>
        </View>
      )
    case TxSummaryStatusInProgress:
      return (
        <View className={`px-2 py-1 rounded-full border-[1px] ${st.borderColor} ${st.bgColor}`}>
          <Text className={`${st.textColor}`}>{st.label}</Text>
        </View>
      )
    case TxSummaryStatusCompleted:
      return (
        <View className={`px-2 py-1 rounded-full border-[1px] border-[#005a32] bg-[#c7e9c0]`}>
          <Text className={`text-[#005a32]`}>{st.label}</Text>
        </View>
      )
    case TxSummaryStatusCancelled:
      return (
        <View className={`px-2 py-1 rounded-full border-[1px] ${st.borderColor} ${st.bgColor}`}>
          <Text className={`${st.textColor}`}>{st.label}</Text>
        </View>
      )
    case TxSummaryStatusPickup:
      return (
        <View className={`px-2 py-1 rounded-full border-[1px] ${st.borderColor} ${st.bgColor}`}>
          <Text className={`${st.textColor}`}>{st.label}</Text>
        </View>
      )
    case TxSummaryStatusRefunding:
      return (
        <View className={`px-2 py-1 rounded-full border-[1px] ${st.borderColor} ${st.bgColor}`}>
          <Text className={`${st.textColor}`}>{st.label}</Text>
        </View>
      )
    default:
      <></>
  }
}

export default TxSummaryChip