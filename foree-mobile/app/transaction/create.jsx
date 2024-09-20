import { number, object, string } from "yup"
import { View, Text, SafeAreaView } from 'react-native'
import React, { useCallback, useEffect, useMemo, useState } from 'react'

import { accountService, transactionService } from '../../service'
import MultiStepForm from '../../components/MultiStepForm'
import { router, useFocusEffect, useLocalSearchParams } from 'expo-router'
import ModalSelect from '../../components/ModalSelect'
import { ContactTransferCashPickup } from '../../constants/contacts'
import { formatContactMethod, formatName } from '../../util/foree_util'
import CurrencyInputField from '../../components/CurrencyInputField'
import Currencies from '../../constants/currency'
import { TransactionPurposes } from '../../constants/transactions'

const TransactionCreate = () => {
  const {contactId} = useLocalSearchParams()
  const [rate, setRate] = useState(0.0)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isEditable, setIsEditable]  = useState(false)
  const [contacts, setContacts] = useState([])
  const [sourceAccounts , setSourceAccounts] = useState([])
  const [errors, setErrors] = useState({})
  const [dailyLimit, setDailyLimit] = useState(null)
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
  const [quote, setQuote] = useState(null)
  const [txOutage, setTxOutage] = useState(null)
  const [rewards, setRewards] = useState([])

  const quoteTransactionScheme = useMemo(() => {
    const sourceAmoutScheme = !!dailyLimit ? 
      number().required("required").min(20, "Minimum $20.00 CAD").max(dailyLimit.maxAmount-dailyLimit.usedAmount, `Maximum $${(dailyLimit.maxAmount-dailyLimit.usedAmount).toFixed(2)} CAD`) :
      number().required("required").min(20, "Minimum $20.00 CAD").max(1000, "Maximum $1000.00 CAD")
    return object({
      cinAccId: number().required("required"),
      coutAccId: number().integer().required("required").min(1, "required"),
      srcAmount: sourceAmoutScheme,
      transactionPurpose: string().required("required")
    })
  }, [dailyLimit])

  useEffect(() => {
    async function validate() {
      try {
        await quoteTransactionScheme.validate(form, {abortEarly: false})
        setErrors({})
      } catch (err) {
        let e = {}
        for ( let i of err.inner ) {
          e[i.path] =  e[i.path] ?? i.errors[0]
        }
        setErrors(e)
        console.log(e)
      }
    }
    validate()
  }, [form])

  // useEffect(() => {
  //   console.log("aaa", contacts.length, sourceAccounts.length)
  //   if ( !!contacts && contacts.length > 0 && !!sourceAccounts && sourceAccounts.length >0 ) {
  //     setIsEditable(true)
  //   } else {
  //     setIsEditable(false)
  //   }
  // }, [contacts, sourceAccounts])

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
          const destAccs = resp.data.data
          setContacts([...destAccs])
          if ( !!(destAccs.find(x => x.id == parseInt(contactId))) ) {
            setForm((form) => ({
              ...form,
              coutAccId: parseInt(contactId)
            }))
          } else if ( contacts.length === 1 ) {
            setForm((form) => ({
              ...form,
              coutAccId: contacts[0].id
            }))
          }
          setIsEditable(true)
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
          const interacs = resp.data.data
          setSourceAccounts(interacs)
          if ( interacs.length === 1 ) {
            setForm((form) =>({
              ...form,
              cinAccId: interacs[0].id
            }))
          }
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
          const rate = resp.data.data
          setRate(rate)
        }
      } catch (e) {
        console.error("transaction_create--getRate", e)
      }
    }

    const getDailyLimit = async(signal) => {
      try {
        const resp = await transactionService.getDailyLimit({signal})
        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("transaction_create--getDailyLimit", resp.status, resp.data)
        } else {
          const dailyLimit = resp.data.data
          setDailyLimit(dailyLimit)
        }
      } catch (e) {
        console.error("transaction_create--getDailyLimit", e)
      }
    }

    const getRewards = async(signal) => {
      try {
        const resp = await transactionService.getRewards({signal})
        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("transaction_create--getRewards", resp.status, resp.data)
        } else {
          const r = resp.data.data
          console.log('eeeee', r)
          setRewards(r)
        }
      } catch (e) {
        console.error("transaction_create--getRewards", e)
      }
    }

    getAllContacts(controller.signal)
    getSourceAccounts(controller.signal)
    getRate(controller.signal)
    getDailyLimit(controller.signal)
    getRewards(controller.signal)
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

  const quoteTransaction = async () => {
    try {
      console.log(form)
      const resp = await transactionService.quote(form)
      if ( resp.status / 100 !== 2 ) {
        console.warn("quote transaction", resp.status, resp.data)
        return false
      }
      setQuote(resp.data.data)
      console.log(resp.data.data)
      return true
    } catch (e) {
      console.log(e)
      return false
    }
  }

  const InteracListItem = useCallback((interac) => {
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

  const RewardListItem = useCallback((reward) => {
    if ( !reward ) return <></>
    return (
      <View>
        <Text>{reward.description}</Text>
      </View>
    )
  }, [form.rewardIds])

  const TransactionCreateTitle = useCallback(() => {
    return (
      <View>
        <Text className="text-lg font-pbold text-center">Transaction Details</Text>
      </View>
    )
  }, [])

  //TODO: apply use callback
  const TransactionCreate = useCallback(() => (
    <View>
      <Text className="font-pregular text-center mb-4">
        Enter the details for your transactions.
      </Text>
      <ModalSelect
        title="From"
        modalTitle="Select Source"
        errorMessage={errors['cinAccId']}
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
          setForm((form) => ({
            ...form,
            cinAccId: o
          }))
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
          setForm((form) => ({
            ...form,
            coutAccId: o
          }))
        }}
        placeholder="Select Contact"
      />
      {/*  */}
      <CurrencyInputField
        title="You Send"
        defaultValue={form.srcAmount}
        containerStyles="mt-2"
        errorMessage={errors['srcAmount']}
        placeholder={!dailyLimit ? "type amount..." : `available: ${(dailyLimit.maxAmount-dailyLimit.usedAmount).toFixed(2)}`}
        onCurrencyChange={((e) => {
          setForm((form) => ({
            ...form,
            srcAmount: e.amount,
            destAmount: Math.floor(e.amount*rate.destAmount*100) / 100
          }))
        })}
        supportCurrencies={[Currencies["CAD"]]}
      />
      <CurrencyInputField
        title="Recipient Receives"
        containerStyles="mt-2"
        producer={false}
        value={form.destAmount}
        editable={false}
        supportCurrencies={[Currencies["PKR"]]}
      />
      <View className="mt-2">
        <Text className="font-semibold text-green-800">Current Rate: <Text className="text-green-600">{rate?.description}</Text></Text>
      </View>
      {
        !!rewards && rewards.length > 0 ?
          <ModalSelect
          title="Apply Promotion"
          modalTitle="apply promotions"
          containerStyles="mt-2"
          value={() => {
            let totalReward = 0
            let totalRewardCurrency = ''
            form.rewardIds.forEach((id) => {
              let reward = rewards.find(r => r.id === id)
              if ( !!reward ) {
                totalReward += reward.amount
                totalRewardCurrency = reward.currency
              }
            })
            if ( totalReward === 0 ) return ''
            return new Intl.NumberFormat("en", {minimumFractionDigits: 2}).format(totalReward) + (!!totalRewardCurrency ? ` ${totalRewardCurrency}` : '')
          }}
          keyExtractor="id"
          valueExtractor="id"
          list={rewards}
          listView={RewardListItem}
          uselistSeperator={false}
          isEditable={isEditable}
          inputStyles="text-right"
          onPress={(w) => {
            setForm((form) => {
              if (!!form.rewardIds.find(x => x === w.id)) {
                return {
                  ...form,
                  rewardIds: [...form.rewardIds.filter(x => x !== w.id)]
                }
              } else {
                return {
                  ...form,
                  rewardIds: [...form.rewardIds, w.id]
                }
              }
            })
            console.log(form)
          }}
          placeholder="...choose"
        />
        :<></>
      }
    </View>
  ), [
    rate,
    rewards,
    dailyLimit,
    isEditable,
    contacts, 
    sourceAccounts, 
    form.destAmount, 
    form.cinAccId, 
    form.coutAccId, 
    errors['srcAmount'],
    errors['cinAccId'],
    errors['coutAccId']
  ])
  
  const TransactionPurposeTitle = useCallback(() => (
    <View>
      <Text className="text-lg font-pbold text-center">Transaction Purpose</Text>
    </View>
  ), [])

  const TransactionPurpose = useCallback(() => (
    <View>
      <Text className="font-pregular text-center mb-4">
        Enter the details for your transactions.
      </Text>
      <ModalSelect
        title="Transaction Purpose"
        modalTitle="select a purpose"
        errorMessage={errors['transactionPurpose']}
        containerStyles="mt-2"
        value={form.transactionPurpose}
        keyExtractor="name"
        showExtractor="name"
        valueExtractor="name"
        listView={(purpose) => (
          <Text className="font-pregular py-3 text-xl">
            {purpose["name"]}
          </Text>
        )}
        list={Object.values(TransactionPurposes)}
        onPress={(o) => {
          setForm((form) =>({
            ...form,
            transactionPurpose: o
          }))
        }}
        placeholder="Choose a transaction purpose"
      />
    </View>
  ), [errors['transactionPurpose'], form.transactionPurpose])

  const ReviewTitle = useCallback(() => (
    <View>
      <Text className="text-lg font-pbold text-center">Review</Text>
    </View>
  ), [])

  const Review = useCallback(() => {
    const txSum = quote.txSum
    if ( !txSum ) return <></>
    return (
      <View>
        <Text className="font-pregular text-center mb-4">
          Review your transactions.
        </Text>
        <View>
          <View className="border-b-[1px] border-slate-400 pb-2">
            <Text className="font-semibold">You Send</Text>
            <Text className="font-bold mb-1 text-lg">${new Intl.NumberFormat("en", {minimumFractionDigits: 2}).format(txSum.srcAmount)} {txSum.srcCurrency}</Text>
            <Text className="font-semibold text-slate-500">From</Text>
            <Text className="font-bold mb-1 text-lg">{formatName(txSum.srcAccount)}</Text>
            <Text className="font-semibold text-slate-500">Interac E-Transfer</Text>
            <Text className="font-bold mb-1 text-lg">{txSum.srcAccount.email}</Text>
          </View>
          <View className="mt-2 border-b-[1px] border-slate-400 pb-2">
            <Text className="font-semibold">Recipient Receives</Text>
            <Text className="font-bold mb-1 text-lg">${new Intl.NumberFormat("en", {minimumFractionDigits: 2}).format(txSum.destAmount)} {txSum.destCurrency}</Text>
            <Text className="font-semibold text-slate-500">To</Text>
            <Text className="font-bold mb-1 text-lg">{formatName(txSum.destAccount)}</Text>
            <Text className="font-semibold">Destination Account</Text>
            <Text className="font-bold mb-1 text-lg">
              {
                txSum.destAccount.type == ContactTransferCashPickup ?
                <Text className="font-bold mb-1 text-lg">Cash Pickup</Text> : 
                <Text className="font-bold mb-1 text-lg">
                  {!! txSum.destAccount.institutionName ?  txSum.destAccount.institutionName.slice(0, 16) + ( txSum.destAccount.institutionName.length > 16 ? "..." : "") : ""}
                  <Text className="italic">
                    ({!!txSum.destAccount.accountNumber ? txSum.destAccount.accountNumber.slice(0, 13) + (txSum.destAccount.accountNumber.length > 13 ? "..." : "") : ""})
                  </Text>
                </Text>
              }
            </Text>
          </View>
          <View className="mt-2">
              <Text className="font-semibold mb-2 text-lg">Details</Text>
              <View className="flex flex-row justify-between items-center mb-1">
                <Text className="font-semibold">Exchange Rate</Text>
                <Text className="font-bold" >{txSum.rate}</Text>
              </View>
              <View className="flex flex-row justify-between items-center mb-1">
                <Text className="font-semibold">Fees</Text>
                <Text className="font-bold" >${new Intl.NumberFormat("en", {minimumFractionDigits: 2}).format(txSum.feeAmount)}{!!txSum.feeCurrency ? ` ${txSum.feeCurrency}` : ''}</Text>
              </View>
              <View className="flex flex-row justify-between items-center mb-1">
                <Text className="font-semibold">Promotions</Text>
                <Text className="font-bold" >${new Intl.NumberFormat("en", {minimumFractionDigits: 2}).format(txSum.rewardAmount)}{!!txSum.rewardCurrency ? ` ${txSum.rewardCurrency}` : ''}</Text>
              </View>
              <View className="flex flex-row justify-between items-center mb-1">
                <Text className="font-semibold text-green-800">Amount Debited</Text>
                <Text className="font-bold text-green-800" >${new Intl.NumberFormat("en", {minimumFractionDigits: 2}).format(txSum.totalAmount)} {txSum.totalCurrency}</Text>
              </View>
          </View>
        </View>
      </View>
    )
  }, [quote])

  const CreateTransactionFlow = [
    {
      titleView: TransactionCreateTitle,
      formView: TransactionCreate,
      canGoNext: () => {
        return !errors.cinAccId &&
                !errors.coutAccId &&
                !errors.srcAmount
      }
    },
    {
      titleView: TransactionPurposeTitle,
      formView: TransactionPurpose,
      canGoNext: () => {
        return !errors.transactionPurpose
      },
      goNext: async () => {
        return await quoteTransaction()
      }
    },
    {
      titleView: ReviewTitle,
      formView: Review,
      canGoNext: () => {
        return true
      },
    },
  ]

  return (
    <SafeAreaView className="bg-slate-100">
      <MultiStepForm
        steps={() => CreateTransactionFlow}
        onSumbit={submit}
        containerStyle=""
        submitDisabled={isSubmitting}
        submitTintTitle="Send"
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