import { View, Text, SafeAreaView } from 'react-native'
import React, { useState } from 'react'

const TransactionCreate = () => {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [errors, setErrors] = useState({})
  const [form, setForm] = useState({
    cinAccId: 0,
    coutAccId: 0,
    srcAmount: 0,
    srcCurrency: 'CAD',
    destCurrency: 'PKR',
    rewardIds: [],
    transactionPurpose: ''
  })

  return (
    <SafeAreaView>
    <Text>TransactionCreate</Text>
    </SafeAreaView>
  )
}

export default TransactionCreate