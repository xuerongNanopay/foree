import { View, Text } from 'react-native'
import React from 'react'
import { SummaryTxStatuses, TxSummaryStatusInitial } from '../constants/summary_tx'


const TxSummaryChip = ({status=TxSummaryStatusInitial}={status}) => {
	const style = SummaryTxStatuses[status]
	return (
		<View className={`px-2 py-1 rounded-full border-[1px] ${style.borderColor} ${style.bgColor}`}>
			<Text className={`${style.textColor}`}>{status}</Text>
		</View>
	)
}

export default TxSummaryChip