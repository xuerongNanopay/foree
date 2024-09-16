import { View, Text, SafeAreaView } from 'react-native'
import React, { useCallback, useEffect, useState } from 'react'

import { accountService } from '../../service'
import MultiStepForm from '../../components/MultiStepForm'
import { router, useFocusEffect } from 'expo-router'
import ModalSelect from '../../components/ModalSelect'
import { ContactTransferCashPickup } from '../../constants/contacts'
import { formatContactName } from '../../util/contact_util'

const TransactionCreate = () => {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isEditable, setIsEditable]  = useState(false)
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

  useEffect(() => {
    if ( !!contacts && contacts.length > 0 && !!sourceAccounts && sourceAccounts.length >0 ) {
      setIsEditable(true)
    } else {
      setIsEditable(false)
    }
  }, [contacts, sourceAccounts])

  useFocusEffect(useCallback(() => {
    console.log("vvvvvvv")
    const  controller = new AbortController()
    const getAllContacts = async (signal) => {
      try {
        const resp = await accountService.getAllContactAccounts({signal})

        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("transaction_create--getAllContacts", resp.status, resp.data)
        } else {
          //How do this: because there is cache in getAllContactAccounts
          //TODO: redesign the cache?
          setContacts([...resp.data.data])
        }
      } catch (e) {
        console.error("transaction_create--getAllContacts", e)
      }
    }

    const getSourceAccounts = async(signal) => {
      try {
        const resp = await accountService.getInteracAccounts({signal})

        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("transaction_create--getAllInteracs", resp.status, resp.data)
        } else {
          console.log(resp.data.data)
          setSourceAccounts(resp.data.data)
        }
      } catch (e) {
        console.error("transaction_create--getAllInteracs", e)
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
      <ModalSelect
        title="Send to"
        modalTitle="Select Contact"
        errorMessage={errors['coutAccId']}
        containerStyles="mt-2"
        allowSearch={true}
        allowAdd={true}
        addHandler={() => {
          router.push('/contact/create')
        }}
        value={form.coutAccId}
        variant='flat'
        searchKey={({firstName, middelName, lastName}) => {
          return `${firstName ?? ""}${middelName ?? ""}${lastName ?? ""}`
        }}
        keyExtractor="id"
        showExtractor="firstName"
        valueExtractor="id"
        list={contacts}
        listView={ContactListItem}
        uselistSeperator={false}
        isEditable={isEditable}
        onPress={(o) => {
          console.log(o)
          setForm({
            ...form,
            coutAccId: o
          })
        }}
        placeholder="Select Contact"
      />
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

//TODO: refactor -- Duplicate code from contact_tab 
const ContactListItem = (
  contact
) => {
  if ( ! contact ) return <></>

  return (
    <View className="mb-2 p-2 rounded-lg bg-[#ccded6]">
      <Text className="font-bold">{formatContactName(contact)}</Text>
      <FormatContactTransferInfo {...contact}/>
      <FormatContactTransferRecentActivity {...contact}/>
    </View>
  )
}


const FormatContactTransferInfo = ({transferMethod, bankName, accountNoOrIBAN}) => {
  if ( transferMethod === ContactTransferCashPickup ) 
    return <Text className="font-semibold text-slate-700">Cash Pickup</Text>
  return (
    <Text className="font-semibold text-slate-700">
      {!!bankName ? bankName.slice(0, 14) + (bankName.length > 14 ? "..." : "") : ""}
      <Text className="italic">
        ({!!accountNoOrIBAN ? accountNoOrIBAN.slice(0, 7) + (accountNoOrIBAN.length > 7 ? "..." : "") : ""})
      </Text>
    </Text>
  )
}

const FormatContactTransferRecentActivity = ({latestActiveAt}) => {
  if ( !latestActiveAt ) return <Text className="text-slate-600 italic">Last sent: -</Text>
  return <Text className="text-slate-600 italic">Last sent: {
    new Date(latestActiveAt).toLocaleString()
  }</Text>
}

export default TransactionCreate