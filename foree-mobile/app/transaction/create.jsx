import { View, Text, SafeAreaView } from 'react-native'
import React, { useCallback, useEffect, useState } from 'react'

import { accountService, transactionService } from '../../service'
import MultiStepForm from '../../components/MultiStepForm'
import { router, useFocusEffect } from 'expo-router'
import ModalSelect from '../../components/ModalSelect'
import { ContactTransferCashPickup } from '../../constants/contacts'
import { formatContactMethod, formatName } from '../../util/foree_util'
import CurrencyInputField from '../../components/CurrencyInputField'
import Currencies from '../../constants/currency'

const TransactionCreate = () => {
  const [rate, setRate] = useState(0.0)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isEditable, setIsEditable]  = useState(false)
  const [contacts, setContacts] = useState([])
  const [sourceAccounts , setSourceAccounts] = useState([])
  const [errors, setErrors] = useState({})
  const [form, setForm] = useState({
    cinAccId: 0,
    coutAccId: 0,
    srcAmount: 0.0,
    destAmount: 0.0,
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
          setSourceAccounts(resp.data.data)
        }
      } catch (e) {
        console.error("transaction_create--getAllInteracs", e)
      }
    }

    const getRate = async(signal) => {
      try {
        const resp = await transactionService.getCADToPRKRate({signal})

        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("transaction_create--getRate", resp.status, resp.data)
        } else {
          setRate(resp.data.data)
        }
      } catch (e) {
        console.error("transaction_create--getRate", e)
      }
    }

    getAllContacts(controller.signal)
    getSourceAccounts(controller.signal)
    getRate(controller.signal)
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

  const InteracListItem = useCallback((interac) =>{
    if ( !interac ) return <></>
    return (
      <View className={`mt-4 p-2 py-4 rounded-lg bg-[#ccded6] ${interac.id == form.cinAccId ? "bg-[#9cd1b9]" : ""}`}>
        <Text className="font-bold">{formatName(interac)}</Text>
        <Text className="font-semibold text-slate-700">
          {"Interac"}
          <Text className="italic">
            ({interac.email})
          </Text>
        </Text>
      </View>
    )
  }, [form.cinAccId])

  const ContactListItem = useCallback((
    contact
  ) => {
    if ( ! contact ) return <></>
  
    return (
      <View className={`mb-2 p-2 rounded-lg bg-[#ccded6] ${contact.id == form.coutAccId ? "bg-[#9cd1b9]" : ""}`}>
        <Text className="font-bold">{formatName(contact)}</Text>
        <FormatContactTransferInfo {...contact}/>
        <FormatContactTransferRecentActivity {...contact}/>
      </View>
    )
  },[form.coutAccId])

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
        title="From"
        modalTitle="Select Source"
        errorMessage={errors['coutAccId']}
        containerStyles="mt-2"
        value={() => {
          if ( !form.cinAccId ) return ""
          const sourcAcc = sourceAccounts.find(c => c.id === form.cinAccId)
          if ( !sourcAcc ) return ""
          return formatName(sourcAcc) + "\nInterac(" + sourcAcc.email + ")"
        }}
        inputContainerStyles="h-16"
        multiline={true}
        numberOfLines={2}
        keyExtractor="id"
        valueExtractor="id"
        list={sourceAccounts}
        listView={InteracListItem}
        uselistSeperator={false}
        isEditable={isEditable}
        onPress={(o) => {
          console.log(o)
          setForm({
            ...form,
            cinAccId: o
          })
        }}
        placeholder="Select Source"
      />
      <ModalSelect
        title="Send to"
        modalTitle="Select Contact"
        errorMessage={errors['coutAccId']}
        containerStyles="mt-2"
        allowSearch={true}
        searchTitle="search name..."
        allowAdd={true}
        addHandler={() => {
          router.push('/contact/create')
        }}
        value={() => {
          if ( !form.coutAccId ) return ""
          const contact = contacts.find(c => c.id === form.coutAccId)
          if ( !contact ) return ""
          return formatName(contact) + "\n" + formatContactMethod(contact)
        }}
        inputContainerStyles="h-16"
        multiline={true}
        numberOfLines={2}
        searchKey={({firstName, middelName, lastName}) => {
          return `${firstName ?? ""}${middelName ?? ""}${lastName ?? ""}`
        }}
        keyExtractor="id"
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
      {/*  */}
      <CurrencyInputField
        title="You Send"
        containerStyles="mt-2"
        placeholder="type amount..."
        onCurrencyChange={((e) => {
          setForm({
            ...form,
            srcAmount: e.amount,
            destAmount: Math.floor(e.amount*rate.destAmount*100) / 100
          })
        })}
        supportCurrencies={[Currencies["CAD"]]}
      />
      <CurrencyInputField
        title="Recipient Receives"
        containerStyles="mt-2"
        value={form.destAmount}
        editable={false}
        supportCurrencies={[Currencies["PKR"]]}
      />
      <View className="mt-2">
        <Text className="font-semibold text-green-800">Current Rate: <Text className="text-green-600">{rate?.description}</Text></Text>
      </View>
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