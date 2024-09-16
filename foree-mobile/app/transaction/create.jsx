import { View, Text, SafeAreaView } from 'react-native'
import React, { useCallback, useEffect, useState } from 'react'

import { accountService } from '../../service'
import MultiStepForm from '../../components/MultiStepForm'
import { useFocusEffect } from 'expo-router'

const TransactionCreate = () => {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [contacts, setContacts] = useState([])
  const [sourceAccounts , setSourceAccounts] = useState([])
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

  useFocusEffect(useCallback(() => {
    const  controller = new AbortController()
    const getAllContacts = async (signal) => {
      try {
        const resp = await accountService.getAllContactAccounts({signal})

        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("get all active contacts", resp.status, resp.data)
        } else {
          //How do this: because there is cache in getAllContactAccounts
          //TODO: redesign the cache?
          setContacts([...resp.data.data])
        }
      } catch (e) {
        console.error(e)
      }
    }

    const getSourceAccounts = async(signal) => {
      try {
        const resp = await accountService.getInteracAccounts({signal})

        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("get all active contacts", resp.status, resp.data)
        } else {
          console.log(resp.data.data)
          setSourceAccounts(resp.data.data)
        }
      } catch (e) {
        console.error(e)
      }
    }
    getAllContacts(controller.signal)
    getSourceAccounts(controller.signal)
    return () => {
      controller.abort()
    }
  }, []))

  const submit = async () => {
    setIsSubmitting(true)
    try {
      console.log("TODO: query transaction")
    } catch (e) {

    } finally {
      setIsSubmitting(false)
    }
  
  }

  const TransactionCreateTitle = () => (
    <View>
      <Text className="text-lg font-pbold text-center">Transaction Details</Text>
    </View>
  )

  const TransactionCreate = () => (
    <View>
      <Text className="font-pregular text-center mb-4">
        Enter the details for your transactions.
      </Text>

    </View>
  )

  const CreateTransactionFlow = [
    {
      titleView: TransactionCreateTitle,
      formView: TransactionCreate,
      canGoNext: () => {
        return true
      }
    },
  ]

  return (
    <SafeAreaView className="bg-slate-100">
      <MultiStepForm
        steps={() => CreateTransactionFlow}
        onSumbit={submit}
        containerStyle=""
        submitDisabled={isSubmitting}
      />
    </SafeAreaView>
  )
}

export default TransactionCreate