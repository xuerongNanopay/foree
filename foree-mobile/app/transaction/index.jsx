import { View, Text } from 'react-native'
import React from 'react'
import { Redirect } from 'expo-router'

const TransactionIdex = () => {
  return (
    <View>
      <Redirect href="/transaction_tab"/>
    </View>
  )
}

export default TransactionIdex