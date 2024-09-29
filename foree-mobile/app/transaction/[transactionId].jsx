import { View, Text } from 'react-native'
import { useLocalSearchParams } from 'expo-router'
import React, { useState, useEffect } from 'react'
import { transactionService } from '../../service'

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
    <View>
    <Text>TransactionDetail</Text>
    </View>
  )
}

export default TransactionDetail