import { View, Text, ScrollView, SafeAreaView, TouchableOpacity, StyleSheet } from 'react-native'
import * as Linking from 'expo-linking';
import { useLocalSearchParams } from 'expo-router'
import React, { useState, useEffect, useCallback, useMemo } from 'react'
import { transactionService } from '../../service'
import { SummaryTxStatuses, TxSummaryStatusAwaitPayment, TxSummaryStatusCancelled, TxSummaryStatusCompleted, TxSummaryStatusInitial, TxSummaryStatusInProgress } from '../../constants/summary_tx'
import { currencyFormatter } from '../../util/foree_util';


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
            <Text className="font-semibold text-slate-800 text-xl mb-2">{currencyFormatter(sumTx.totalAmount, sumTx.totalCurrency)}</Text>
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

//Why do this: it looks tailwind has issue with dynamic css.
const StatusView = (tx) => {
  const sumTx = SummaryTxStatuses[tx?.status]
  switch (tx.status) {
    case TxSummaryStatusInitial:
      return (
        <View className={`border-[1px] p-2 rounded-md border-purple-800 bg-purple-200`}>
          <View
            className={`border-b-[1px] border-purple-400`}
          >
            <Text className={`font-semibold text-lg text-purple-800`}>
              {sumTx.label}
            </Text>
          </View>
          <View className="mt-2">
            <Text className={`text-purple-800`}>{sumTx.description}</Text>
          </View>
        </View>
      )
    case TxSummaryStatusAwaitPayment:
      return (
        <View className={`border-[1px] p-2 rounded-md border-yellow-800 bg-yellow-200`}>
          <View
            className={`border-b-[1px] border-yellow-400`}
          >
            <Text className={`font-semibold text-lg text-yellow-800`}>
              {sumTx.label}
            </Text>
          </View>
          <View className="mt-2">
            <Text className={`text-yellow-800`}>{sumTx.description}</Text>
          </View>
          {
            !tx.paymentUrl ? <></>:
            <View className="flex flex-row justify-end">
              <TouchableOpacity
                onPress={() => {
                  Linking.openURL(tx.paymentUrl)
                }}
                className={`mr-3 border-2 ${sumTx.borderColor} py-1 px-2 rounded-lg`}
              >
                <Text className={`font-psemibold ${sumTx.textColor}`}>Pay Now</Text>
              </TouchableOpacity>
            </View>
          }
        </View>
      )
    case TxSummaryStatusInitial:
      return (
        <View className={`border-[1px] p-2 rounded-md border-purple-800 bg-purple-200`}>
          <View
            className={`border-b-[1px] border-purple-300`}
          >
            <Text className={`font-semibold text-lg text-purple-800`}>
              {sumTx.label}
            </Text>
          </View>
          <View className="mt-2">
            <Text className={`text-purple-800`}>{sumTx.description}</Text>
          </View>
        </View>
      )
    case TxSummaryStatusInProgress:
      return (
        <View className={`border-[1px] p-2 rounded-md border-purple-800 bg-purple-200`}>
          <View
            className={`border-b-[1px] border-purple-300`}
          >
            <Text className={`font-semibold text-lg text-purple-800`}>
              {sumTx.label}
            </Text>
          </View>
          <View className="mt-2">
            <Text className={`text-purple-800`}>{sumTx.description}</Text>
          </View>
        </View>
      )
    case TxSummaryStatusCompleted:
      return (
        <View className={`border-[1px] p-2 rounded-md border-green-800 bg-green-200`}>
          <View
            className={`border-b-[1px] border-green-300`}
          >
            <Text className={`font-semibold text-lg text-green-800`}>
              {sumTx.label}
            </Text>
          </View>
          <View className="mt-2">
            <Text className={`text-green-800`}>{sumTx.description}</Text>
          </View>
        </View>
      )
    case TxSummaryStatusCancelled:
      return (
        <View className={`border-[1px] p-2 rounded-md border-red-800 bg-red-200`}>
          <View
            className={`border-b-[1px] border-red-300`}
          >
            <Text className={`font-semibold text-lg text-red-800`}>
              {sumTx.label}
            </Text>
          </View>
          <View className="mt-2">
            <Text className={`text-red-800`}>{sumTx.description}</Text>
          </View>
        </View>
      )
    case TxSummaryStatusPickup:
      return (
        <View className={`border-[1px] p-2 rounded-md border-red-800 bg-red-200`}>
          <View
            className={`border-b-[1px] border-red-300`}
          >
            <Text className={`font-semibold text-lg text-red-800`}>
              {sumTx.label}
            </Text>
          </View>
          <View className="mt-2">
            <Text className={`text-red-800`}>{sumTx.description}</Text>
          </View>
        </View>
      )
    case TxSummaryStatusRefunding:
      return (
        <View className={`border-[1px] p-2 rounded-md border-purple-800 bg-purple-200`}>
          <View
            className={`border-b-[1px] border-purple-300`}
          >
            <Text className={`font-semibold text-lg text-purple-800`}>
              {sumTx.label}
            </Text>
          </View>
          <View className="mt-2">
            <Text className={`text-purple-800`}>{sumTx.description}</Text>
          </View>
        </View>
      )
    default:
      return <></>
  }
}

export default TransactionDetail