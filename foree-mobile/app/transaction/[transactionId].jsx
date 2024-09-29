import { View, Text } from 'react-native'
import { useLocalSearchParams } from 'expo-router'
import React from 'react'

const TransactionDetail = () => {
  const {contactId} = useLocalSearchParams()
  return (
    <View>
    <Text>TransactionDetail</Text>
    </View>
  )
}

export default TransactionDetail