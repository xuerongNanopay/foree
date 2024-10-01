import { View, Text, ScrollView, SafeAreaView, TouchableOpacity } from 'react-native'
import * as Linking from 'expo-linking';
import { useLocalSearchParams } from 'expo-router'
import React, { useState, useEffect, useCallback, useMemo } from 'react'
import { transactionService } from '../../service'
import { SummaryTxStatuses, TxSummaryStatusAwaitPayment, TxSummaryStatusInitial } from '../../constants/summary_tx'


const TransactionDetail = () => {
  const { transactionId } = useLocalSearchParams()
  const [sumTx, setSumTx] = useState(null)

  useEffect(() => {
    const controller = new AbortController()
    const getTransactionDetail = async () => {
      try {
        console.log(transactionId)
        const resp = await transactionService.getTransaction(transactionId, {signal: controller.signal})
        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("get transaction detail", resp.status, resp.data)
          router.replace('/transaction')
        } else {
          console.log(resp.data.data)
          setSumTx(resp.data.data)
        }
      } catch (e) {
        console.error(e)
        router.replace('/transaction')
      }
    }
    getTransactionDetail()
    return () => {
      controller.abort()
    }
  }, [])

  const sumTxStatus = useMemo(() => {
    return SummaryTxStatuses[!!sumTx ? sumTx.status : TxSummaryStatusInitial]
  }, [sumTx])
  return (
    <>
      {
        !sumTx ? <></>:
        <SafeAreaView
          className="h-full flex flex-col"
        >
          <View className="mt-6 flex items-center">
            <Text className="text-xl text-slate-600 mb-4">Total Amount</Text>
            <Text className="font-semibold text-slate-800 text-xl mb-2">${new Intl.NumberFormat("en", {minimumFractionDigits: 2}).format(sumTx.totalAmount)}{!!sumTx.totalCurrency ? ` ${sumTx.totalCurrency}` : ''}</Text>
          </View>
          <ScrollView
            className="px-2 pt-4"
          >
            {StatusView(sumTx)}
          </ScrollView>
        </SafeAreaView>
      }
    </>
  )
}

const StatusView = (tx) => {
  const sumTx = SummaryTxStatuses[tx.status]
  return (
    <View className={`border-[1px] p-2 rounded-md ${sumTx.borderColor} ${sumTx.bgColor}`}>
      <View
        className={`border-b-[1px] ${sumTx.borderColor}`}
      >
        <Text className={`font-semibold text-lg ${sumTx.textColor}`}>
          {sumTx.label}
        </Text>
      </View>
      <View className="mt-2">
        <Text className={`${sumTx.textColor}`}>{sumTx.description}</Text>
      </View>
      {
        tx.status !== TxSummaryStatusAwaitPayment ? <></>:
        <View className="flex flex-row justify-end">
          <TouchableOpacity
            onPress={() => {
              Linking.openURL("http://www.google.ca")
            }}
            className={`mr-3 border-2 ${sumTx.borderColor} py-1 px-2 rounded-lg`}
          >
            <Text className={`font-psemibold ${sumTx.textColor}`}>Pay Now</Text>
          </TouchableOpacity>
        </View>
      }
    </View>
  )
}

export default TransactionDetail