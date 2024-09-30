import { View, Text, ScrollView } from 'react-native'
import { useLocalSearchParams } from 'expo-router'
import React, { useState, useEffect } from 'react'
import { transactionService } from '../../service'
import TxSummaryChip from '../../components/TxSummaryChip'

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

  return (
    <>
      {
        !sumTx ? <></>:
        <View className="px-2">
          <View className="mt-6 flex items-center">
            <Text className="text-xl text-slate-600 mb-4">Total Amount</Text>
            <Text className="font-semibold text-slate-800 text-xl mb-4">${new Intl.NumberFormat("en", {minimumFractionDigits: 2}).format(sumTx.totalAmount)}{!!sumTx.totalCurrency ? ` ${sumTx.totalCurrency}` : ''}</Text>
            <View>
              <TxSummaryChip status={sumTx.status}/>
            </View>
          </View>
          <ScrollView>

          </ScrollView>
        </View>
      }
    </>
  )
}

export default TransactionDetail