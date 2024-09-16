import { View, Text, SafeAreaView, TouchableOpacity } from 'react-native'
import React from 'react'

const TransactionTab = () => {
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
          {/* Status Pagenation */}
          <View className="flex flex-row items-center">
            <View className="flex-1">
              <Text>111</Text>
            </View>
            <View className="flex flex-row items-center">
              <Text className="mr-2">1-50 of 22,248</Text>
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
          {/* Status */}
          <View>
            
          </View>
        </View>
      </View>
    </SafeAreaView>
  )
}

export default TransactionTab