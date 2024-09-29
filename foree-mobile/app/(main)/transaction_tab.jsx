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
            {
              selectedStatus !== "All" ? 
              <TouchableOpacity
                onPress={()=> {setSelectedStatus("All")}}
                className="p-2 h-full flex flex-row items-center"
                disabled={false}
              >
              <Image 
                  source={icons.x}
                  className="w-3 h-3"
                  resizeMode='contain'
              />
              </TouchableOpacity> :
              <></>
            }
          </View>
          {/* Status Pagenation */}
          <View className="mt-2 flex flex-row items-center">
            <View className="flex-1">
              <View className="flex flex-row">
                <View className={`border-[1px] rounded-2xl p-2 ${transactionStatuses.find(x=>x.id===selectedStatus).borderColor} ${transactionStatuses.find(x=>x.id===selectedStatus).selectBgColor}`}>
                  <Text className={`font-psemibold ${transactionStatuses.find(x=>x.id===selectedStatus).textColor}`}>{selectedStatus}</Text>
                </View>
              </View>
            </View>
            <View className="flex flex-row items-center">
              <Text className="mr-2">1-50</Text>
              <TouchableOpacity
                onPress={()=> {console.log("TODO: transaction left")}}
                activeOpacity={0.7}
                disabled={false}
                className="mr-2 p-2"
              >
                <Image
                  source={icons.leftArrowDark}
                  className="w-[15px] h-[15px]"
                  resizeMode='contain'
                />
              </TouchableOpacity>
              <TouchableOpacity
                onPress={()=> {console.log("TODO: transaction right")}}
                activeOpacity={0.7}
                disabled={false}
                className="mr-2 p-2"
              >
                <Image
                  source={icons.rightArrowDark}
                  className="w-[15px] h-[15px]"
                  resizeMode='contain'
                />
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
    borderColor: "border-slate-800",
    textColor: "text-slate-800",
    selectBgColor: "bg-slate-200",
  },
  { 
    id:"Initial",
    borderColor: "border-purple-800",
    textColor: "text-purple-800",
    selectBgColor: "bg-purple-200",
  },
  { 
    id:"Await Payment",
    borderColor: "border-yellow-800",
    textColor: "text-yellow-800",
    selectBgColor: "bg-yellow-200",
  },
  { 
    id:"In Progress",
    borderColor: "border-purple-800",
    textColor: "text-purple-800",
    selectBgColor: "bg-purple-200",
  },
  { 
    id:"Completed",
    borderColor: "border-green-800",
    textColor: "text-green-800",
    selectBgColor: "bg-green-200",
  },
  { 
    id:"Cancelled",
    borderColor: "border-red-800",
    textColor: "text-red-800",
    selectBgColor: "bg-red-200",
  },
  { 
    id:"Ready To Pickup",
    borderColor: "border-yellow-800",
    textColor: "text-yellow-800",
    selectBgColor: "bg-yellow-200",
  },
  { 
    id:"Refunding",
    borderColor: "border-purple-800",
    textColor: "text-purple-800",
    selectBgColor: "bg-purple-200",
  }
]

export default TransactionTab