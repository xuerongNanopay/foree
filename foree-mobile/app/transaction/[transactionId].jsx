import { View, Text, ScrollView, SafeAreaView, TouchableOpacity, Image } from 'react-native'
import * as Linking from 'expo-linking';
import { useLocalSearchParams } from 'expo-router'
import React, { useState, useEffect } from 'react'
import { transactionService } from '../../service'
import { SummaryTxStatuses, TxSummaryStatusAwaitPayment, TxSummaryStatusCancelled, TxSummaryStatusCompleted, TxSummaryStatusInitial, TxSummaryStatusInProgress, TxSummaryStatusPickup, TxSummaryStatusRefunding } from '../../constants/summary_tx'
import { currencyFormatter, formatContactMethod, formatName } from '../../util/foree_util';
import { icons } from '../../constants';


const TransactionDetail = () => {
  const { transactionId } = useLocalSearchParams()
  const [sumTx, setSumTx] = useState(null)

  useEffect(() => {
    const controller = new AbortController()
    const getTransactionDetail = async () => {
      try {
        const resp = await transactionService.getTransaction(transactionId, {signal: controller.signal})
        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("get transaction detail", resp.status, resp.data)
          router.replace('/transaction')
        } else {
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

  return (
    <>
      {
        !sumTx ? <></>:
        <SafeAreaView>
          <View className="h-full relative">
            <View className="mt-6 flex items-center">
              <Text className="text-xl text-slate-600 mb-4">Total Amount</Text>
              <Text className="font-semibold text-slate-800 text-xl mb-2">{currencyFormatter(sumTx.totalAmount, sumTx.totalCurrency)}</Text>
            </View>
            {
              //see: https://stackoverflow.com/questions/36938742/touchablehighlight-not-clickable-if-position-absolute
              1 === 1 ? <></> :
              <View
                className="absolute p-1 right-3 top-1"
              >
                <TouchableOpacity
                  activeOpacity={0.7}
                  className="p-2"
                  onPress={() => console.log("more todo1")}
                >
                  <Image
                    source={icons.caretDownDark}
                    className="w-[17px] h-[17px]"
                    resizeMode='contain'
                  />
                </TouchableOpacity>
              </View>
            }
            <ScrollView
              className="px-2 py-4"
              showsVerticalScrollIndicator={false}
            >
              {StatusView(sumTx)}
              <View
                className="mt-4 border-[1px] rounded-md px-2 py-2 border-slate-400"
              >
                <View className="pb-1 border-b-[1px] border-slate-300">
                  <Text className="text-slate-500 mb-1">Created</Text>
                  <Text className="font-psemibold text-slate-600">{new Date(sumTx.createAt).toLocaleString()}</Text>
                </View>
                <View className="mt-2 pb-1">
                  <Text className="text-slate-500 mb-1">Reference #</Text>
                  <Text className="font-psemibold text-slate-600">{sumTx.nbpReference}</Text>
                </View>
                {
                  !sumTx.destAccount ? <></>:
                  <View className="pt-2 border-t-[1px] border-slate-300">
                    <Text className="text-slate-500 mb-1">Remitee</Text>
                    <Text className="font-psemibold text-slate-600">{formatName(sumTx.destAccount)}</Text>
                    <Text className="text-slate-500 mt-1 mb-1">Receive Amount</Text>
                    <Text className="font-psemibold text-slate-600">{currencyFormatter(sumTx.destAmount, sumTx.destCurrency)}</Text>
                    <Text className="text-slate-500 mt-1 mb-1">Receive Method</Text>
                    <Text className="font-psemibold text-slate-600">{formatContactMethod(sumTx.destAccount, max=20)}</Text>
                  </View>
                }
              </View>
              <View
                className="mt-5 mb-10"
              >
                <Text className="font-psemibold text-lg">Details</Text>
                <View className="mt-3 pb-2 flex-row justify-between items-center border-b-[1px] border-slate-300">
                  <Text className="text-slate-500">Exchange Rate</Text>
                  <Text className="font-psemibold text-slate-600">{sumTx.rate}</Text>
                </View>
                <View className="mt-3 pb-2 flex-row justify-between items-center border-b-[1px] border-slate-300">
                  <Text className="text-slate-500">Fees</Text>
                  <Text className="font-psemibold text-slate-600">{currencyFormatter(sumTx.feeAmount, sumTx.feeCurrency)}</Text>
                </View>
                <View className="mt-3 pb-2 flex-row justify-between items-center border-b-[1px] border-slate-300">
                  <Text className="text-slate-500">Rewards</Text>
                  <Text className="font-psemibold text-slate-600">{currencyFormatter(sumTx.rewardAmount, sumTx.rewardCurrency)}</Text>
                </View>
                <View className="mt-3 pb-2 flex-row justify-between items-center border-b-[1px] border-slate-300">
                  <Text className="text-slate-500">Total Amount</Text>
                  <Text className="font-psemibold text-slate-600">{currencyFormatter(sumTx.totalAmount, sumTx.totalCurrency)}</Text>
                </View>
                {
                  1 !== 1 ? <></> :
                  <TouchableOpacity
                    onPress={() => {
                      console.log("TODO: cancel transaction")
                    }}
                    className = "mt-3 p-2 w-24 border border-red-800 rounded-lg flex"
                  >
                  <Text className="font-psemibold text-center text-red-800">Cancel</Text>
                </TouchableOpacity>
                }
              </View>
          </ScrollView>
          </View>
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
        <View className={`border-[1px] p-2 rounded-md border-yellow-800 bg-yellow-200`}>
          <View
            className={`border-b-[1px] border-yellow-300`}
          >
            <Text className={`font-semibold text-lg text-yellow-800`}>
              {sumTx.label}
            </Text>
          </View>
          <View className="mt-2">
            <Text className={`text-yellow-800`}>{sumTx.description}</Text>
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