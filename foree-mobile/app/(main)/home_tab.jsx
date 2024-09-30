import { View, Text, SafeAreaView, FlatList, ScrollView, Image, Touchable, TouchableOpacity } from 'react-native'
import { Link, router, useFocusEffect } from 'expo-router'
import React, { useState, useCallback } from 'react'

import { icons } from '../../constants'
import { useGlobalContext } from '../../context/GlobalProvider'
import { transactionService } from '../../service'
import Currency from '../../constants/currency'
import string_util from '../../util/string_util'
import TxSummaryChip from '../../components/TxSummaryChip'

const HomeTab = () => {
  const { user } = useGlobalContext()
  const [ cpRate , setCPRate ] = useState({
    srcAmount: 0,
    srcCurrency: "CAD",
    destAmount: 0,
    destCurrency: "PKR",
  })
  const [ latestTxs, setLastestTxs] = useState([])
  
  useFocusEffect(useCallback(() => {
    const controller = new AbortController()
    const getRate = async() => {
      try {
        const resp = await transactionService.getCADToPRKRate({signal: controller.signal})
        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("get rate", resp.status, resp.data)
        } else {
          setCPRate({
            ...cpRate,
            ...resp.data.data
          })
        }
      } catch (e) {
        console.error("get rate", e, e.response, e.response?.status, e.response?.data)
      }
    }
    const getLastestTransactions = async() => {
      try {
        const resp = await transactionService.getTransactions({offset:0, limit:5}, {signal: controller.signal})
        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("get transactions", resp.status, resp.data)
        } else {
          setLastestTxs(resp.data.data)
        }
      } catch (e) {
        console.error("get transactions", e, e.response, e.response?.status, e.response?.data)
      }
    }
    getRate()
    getLastestTransactions()
    return () => {
      controller.abort()
    }
  }, []))

  return (
    <SafeAreaView className="h-full flex flex-row items-center mb-28">
      <View className="px-4 pt-4">
        <View className="mb-2 pb-2 border-b-[1px] border-slate-400 flex-row items-center justify-between">
          <View className="">
            <Text className="font-pregular text-xl">Welcome Back</Text>
            <Text className="font-pbold text-2xl text-slate-700">{user?.firstName} {user?.lastName}</Text>
          </View>
          <TouchableOpacity
            onPress={() => {console.log("TODO: notifications")}}
          >
            <Image
              source={icons.bell}
              resizeMode='contain'
              className="w-[30px] h-[30px] mr-2"
            />
          </TouchableOpacity>
        </View>
        <ScrollView 
          showsVerticalScrollIndicator={false}
        >
          <View className="bg-[#ccded6] rounded-2xl p-4 my-4">
            <View className="flex-1">
              <Text className="font-pbold text-lg">Current Rate</Text>
              {/* <Text className="font-psemibold text-lg">🇨🇦 $1.00 CAD = 🇵🇰 $208.00 PKR</Text> */}
              <Text className="font-psemibold text-lg">{`${Currency[cpRate.srcCurrency]?.["unicodeIcon"]} $${cpRate.srcAmount.toFixed(2)} ${cpRate.srcCurrency} = ${Currency[cpRate.destCurrency]?.["unicodeIcon"]} $${cpRate.destAmount.toFixed(2)} ${cpRate.destCurrency}`}</Text>
            </View>
            <View>
              <View className="mt-4 p-2 rounded-xl bg-[#1A6B54]">
                <Link href="/transaction/create">
                  <Text className="text-lg text-center font-semibold text-white">Send Money</Text>
                </Link>
              </View>
            </View>
          </View>
          <View className="bg-[#ccded6] rounded-2xl p-4 my-4">
            <Text className="font-pbold mb-2">Welcome to Foree Remittance, stress free money transfers to ....... in exclusive partnership with ...</Text>
            <Text className="font-psemibold mb-2">Foree brings more value & exciting rewards for new & existing users</Text>
            <View className="pl-2 flex flex-row font-pregular">
              <Text>{"\u2022"}</Text>
              <Text className="pl-2">Every new Sifn-Up gets a $20 credit for a limited time</Text>
            </View>
            <View className="pl-2 flex flex-row font-pregular">
              <Text>{"\u2022"}</Text>
              <Text className="pl-2">Refer a friend or family - they get $20 credit upon sign-up, using your referral link and your get $20 credit when they complete first transaction!</Text>
            </View>
            <Text className="font-psemibold mt-2">Refer today & start earning the rewards</Text>
          </View>
          {
            !latestTxs || latestTxs.length === 0 ? <></> :
            <View className="bg-[#ccded6] rounded-2xl py-4 my-4">
            <View className="px-4 pb-2 border-b-[1px] border-[#b6d4c7]">
              <Text className="font-pbold text-lg">Recent Activities</Text>
            </View>
            {
              latestTxs.map((tx, idx) => {
                return(
                  <TouchableOpacity 
                    key={tx.id}
                    onPress={() => router.push(`/transaction/${tx.id}`)}
                    className={`py-2 ${idx !== latestTxs.length-1 ? "border-b-[1px] border-[#b6d4c7]" : ""}`}
                  >
                    <View className="px-3 mb-1 flex-row items-center justify-between">
                      <Text className="font-semibold">{!!tx.destAccSummary ? string_util.formatStringWithLimit(tx.destAccSummary, 14) : ""}</Text>
                      <Text className="font-semibold text-slate-600">${new Intl.NumberFormat("en", {minimumFractionDigits: 2}).format(tx.destAmount)}{!!tx.destCurrency ? ` ${tx.destCurrency}` : ''}</Text>
                    </View>
                    <View className="px-3 mb-1 flex-row items-center justify-between">
                      <Text className="font-semibold">Total Amount</Text>
                      <Text className="font-semibold text-slate-600">${new Intl.NumberFormat("en", {minimumFractionDigits: 2}).format(tx.totalAmount)}{!!tx.totalCurrency ? ` ${tx.totalCurrency}` : ''}</Text>
                    </View>
                    <View className="px-3 flex-row items-center justify-between">
                      <Text className="italic text-slate-600">{tx.nbpReference}</Text>
                      <TxSummaryChip
                        status={tx.status}
                      />
                    </View>
                  </TouchableOpacity>
                )
              })
            }
            <View className="px-4 border-t-[1px] border-[#b6d4c7]">
              <Link href="/transaction" className="pt-2">
                <Text className="font-pregular text-center">See more...</Text>
              </Link>
            </View>
          </View>
          }
        </ScrollView>
      </View>
    </SafeAreaView>
  )
}

export default HomeTab