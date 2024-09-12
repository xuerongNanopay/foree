import { View, Text, SafeAreaView, FlatList, ScrollView } from 'react-native'
import { Link } from 'expo-router'
import React, { useEffect, useState } from 'react'
import { useGlobalContext } from '../../context/GlobalProvider'
import { transactionService } from '../../service'

import Currency from '../../constants/currency'

const Home = () => {
  const { user } = useGlobalContext()
  const [ cpRate , setCPRate ] = useState({
    srcAmount: 0,
    srcCurrency: "CAD",
    destAmount: 0,
    destCurrency: "PKR",
  })
  
  useEffect(() => {
    const getRate = async() => {
      try {
        resp = await transactionService.getCADToPRKRate()
        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("get rate", resp.status, resp.data)
        } else {
          console.log(resp.data.data)
          setCPRate({
            ...cpRate,
            ...resp.data.data
          })
        }
      } catch (e) {
        console.error("get rate", e.response?.status, e.response?.data)
      }
  
    }
    getRate()
  }, [])
  return (
    <SafeAreaView className="h-full flex flex-row items-center mb-16">
      <View className="px-4 pt-4">
        <View className="mb-4">
          <Text className="font-pregular text-xl">Welcome Back</Text>
          <Text className="font-pbold text-2xl text-slate-700">{user?.firstName} {user?.lastName}</Text>
        </View>
        <ScrollView 
          showsVerticalScrollIndicator={false}
        >
          <View className="bg-[#d0f2e4] rounded-2xl p-4 my-4">
            <View className="flex-1">
              <Text className="font-pbold text-lg">Current Rate</Text>
              {/* <Text className="font-psemibold text-lg">🇨🇦 $1.00 CAD = 🇵🇰 $208.00 PKR</Text> */}
              <Text className="font-psemibold text-lg">{`${Currency[cpRate.srcCurrency]?.["unicodeIcon"]} $${cpRate.srcAmount.toFixed(2)} ${cpRate.srcCurrency} = ${Currency[cpRate.destCurrency]?.["unicodeIcon"]} $${cpRate.destAmount.toFixed(2)} ${cpRate.destCurrency}`}</Text>
            </View>
            <View>
              <View className="mt-4 p-2 rounded-xl bg-[#1A6B54]">
                <Link href="/profile">
                  <Text className="text-lg text-center font-semibold text-white">Send Money</Text>
                </Link>
              </View>
            </View>
          </View>
          <View className="bg-[#d0f2e4] rounded-2xl p-4 my-4">
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
          <View className="bg-[#d0f2e4] rounded-2xl py-4 my-4">
            <View className="px-4 pb-2 border-b-[1px] border-[#b6d4c7]">
              <Text className="font-pbold text-lg">Recent Activities</Text>
            </View>
            {/* <FlatList
              data={[{id:1}, {id:2}]}
              keyExtractor={(item) => item.id}
              renderItem={({item}) => (
                <Text>{item.id}</Text>
              )}
            /> */}
            {/* <Text>1111</Text>
            <Text>1111</Text>
            <Text>1111</Text>
            <Text>1111</Text>
            <Text>1111</Text>
            <Text>1111</Text>
            <Text>1111</Text>
            <Text>1111</Text>
            <Text>1111</Text>
            <Text>1111</Text>
            <Text>1111</Text>
            <Text>1111</Text>
            <Text>1111</Text>
            <Text>1111</Text>
            <Text>1111</Text> */}
            <View className="px-4 border-t-[1px] border-[#b6d4c7]">
              <Link href="/profile" className="pt-2">
                <Text className="font-pregular text-center">See more...</Text>
              </Link>
            </View>
          </View>
        </ScrollView>
      </View>
    </SafeAreaView>
  )
}

export default Home