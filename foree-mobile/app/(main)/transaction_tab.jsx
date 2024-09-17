import { View, Text, SafeAreaView, TouchableOpacity, FlatList, Image } from 'react-native'
import React, { useCallback, useState } from 'react'
import { router, useFocusEffect } from 'expo-router'

import { icons } from '../../constants'

const TransactionTab = () => {
  const [selectedStatus, setSelectedStatus] = useState('All')

  useFocusEffect(useCallback(() => {
    setSelectedStatus('All')
    const  controller = new AbortController()
    return () => {
      controller.abort()
    }
  }, []))

  const statusChipItem = useCallback(({item}) => {
    const bgColor = selectedStatus === item.id ? `${item.selectBgColor}` : ""
    return (
      <TouchableOpacity 
        onPress={() => setSelectedStatus(item.id)}
        className={`p-2 border-[1px] ${item.borderColor} rounded-2xl mr-2 ${bgColor}`}
      >
        <Text className={`${item.textColor}`}>{item.id}</Text>
      </TouchableOpacity>
    )
  },[selectedStatus])

  return (
    <SafeAreaView>
      <View className="flex h-full px-4 pt-4">
        <View className="pb-2 mb-4 border-b-[1px] border-slate-300">
          {/* Title */}
          <View className="flex flex-row items-center pb-2 mb-2 border-b-[1px] border-slate-300">
            <Text className="flex-1 font-pbold text-2xl">Transactions</Text>
            <TouchableOpacity
              onPress={()=> {router.push("/transaction/create")}}
              activeOpacity={0.7}
              className="bg-[#1A6B54] py-2 px-4 rounded-full"
              disabled={false}
            >
              <Text className="font-pextrabold text-white">Send</Text>
            </TouchableOpacity>
          </View>
          {/* Status */}
          <View className="flex flex-row items-center">
            <TouchableOpacity
              onPress={() => {console.log("TODO: transaction refresh")}}
              className="border-[1px] border-slate-400 rounded-lg p-1"
            >
              <Image
                source={icons.renewable}
                className="w-[27px] h-[27px]"
                resizeMode='contain'
              />
            </TouchableOpacity>
            <FlatList
              className="flex-1 mx-2"
              data={transactionStatuses}
              renderItem={statusChipItem}
              keyExtractor={item => item.id}
              showsHorizontalScrollIndicator={false}
              horizontal={true}
            />
          </View>
          {/* Status Pagenation */}
          <View className="flex flex-row items-center">
            <View className="flex-1">
              <View className>
                <Text className="text-green-800 font-psemibold">{selectedStatus}</Text>
              </View>
            </View>
            <View className="flex flex-row items-center">
              <Text className="mr-2">1-50</Text>
              <TouchableOpacity
                onPress={()=> {console.log("TODO: transaction left")}}
                activeOpacity={0.7}
                disabled={false}
                className="mr-2"
              >
                <Text className="text-2xl">◀️</Text>
              </TouchableOpacity>
              <TouchableOpacity
                onPress={()=> {console.log("TODO: transaction right")}}
                activeOpacity={0.7}
                disabled={false}
              >
                <Text className="text-2xl">▶️</Text>
              </TouchableOpacity>
            </View>
          </View>
        </View>
      </View>
    </SafeAreaView>
  )
}

const transactionStatuses = [
  {
    id:"All",
    borderColor: "border-green-800",
    textColor: "text-green-800",
    selectBgColor: "bg-[#9cd1b9]",
  },
  { 
    id:"Initial",
    borderColor: "border-green-800",
    textColor: "text-green-800",
    selectBgColor: "bg-[#9cd1b9]",
  },
  { 
    id:"Await Payment",
    borderColor: "border-green-800",
    textColor: "text-green-800",
    selectBgColor: "bg-[#9cd1b9]",
  },
  { 
    id:"In Progress",
    borderColor: "border-green-800",
    textColor: "text-green-800",
    selectBgColor: "bg-[#9cd1b9]",
  },
  { 
    id:"Completed",
    borderColor: "border-green-800",
    textColor: "text-green-800",
    selectBgColor: "bg-[#9cd1b9]",
  },
  { 
    id:"Cancelled",
    borderColor: "border-green-800",
    textColor: "text-green-800",
    selectBgColor: "bg-[#9cd1b9]",
  },
  { 
    id:"Ready To Pickup",
    borderColor: "border-green-800",
    textColor: "text-green-800",
    selectBgColor: "bg-[#9cd1b9]",
  },
  { 
    id:"Refunding",
    borderColor: "border-green-800",
    textColor: "text-green-800",
    selectBgColor: "bg-[#9cd1b9]",
  },
  { 
    id:"Refunded",
    borderColor: "border-green-800",
    textColor: "text-green-800",
    selectBgColor: "bg-[#9cd1b9]",
  }
]

export default TransactionTab