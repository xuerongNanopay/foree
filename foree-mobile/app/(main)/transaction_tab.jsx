import { View, Text, SafeAreaView, TouchableOpacity, FlatList } from 'react-native'
import React, { useState } from 'react'

const transactionStatuses = [
  {
    id:"All"
  },
  { 
    id:"Initial"
  },
  { 
    id:"Await Payment"
  },
  { 
    id:"In Progress"
  },
  { 
    id:"Completed"
  },
  { 
    id:"Cancelled"
  },
  { 
    id:"Ready To Pickup"
  },
  { 
    id:"Refunding"
  },
  { 
    id:"Refunded"
  }
]

const TransactionTab = () => {
  const [selectedStatus, setSelectedStatus] = useState('All')

  const statusChipItem = ({item}) => {
    const bgColor = selectedStatus === item.id ? "bg-green-200" : ""
    return (
      <TouchableOpacity 
        onPress={() => setSelectedStatus(item.id)}
        className={`p-2 border-[1px] border-green-800 rounded-2xl mr-2 ${bgColor}`}
      >
        <Text className="text-green-800">{item.id}</Text>
      </TouchableOpacity>
    )
  }

  return (
    <SafeAreaView>
      <View className="flex h-full px-4 pt-4">
        <View className="pb-2 mb-4 border-b-[1px] border-slate-300">
          {/* Title */}
          <View className="flex flex-row items-center pb-2 mb-2 border-b-[1px] border-slate-300">
            <Text className="flex-1 font-pbold text-2xl">Transactions</Text>
            <TouchableOpacity
              onPress={()=> {router.push("/contact/create")}}
              activeOpacity={0.7}
              className="bg-[#1A6B54] py-2 px-4 rounded-full"
              disabled={false}
            >
              <Text className="font-pextrabold text-white">Send</Text>
            </TouchableOpacity>
          </View>
          {/* Status */}
          <FlatList
            data={transactionStatuses}
            renderItem={statusChipItem}
            keyExtractor={item => item.id}
            showsHorizontalScrollIndicator={false}
            horizontal={true}
            className="mb-1"
          />
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

export default TransactionTab