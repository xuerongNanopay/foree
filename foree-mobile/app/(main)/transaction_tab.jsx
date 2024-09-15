import { View, Text, SafeAreaView, TouchableOpacity } from 'react-native'
import React from 'react'

const TransactionTab = () => {
  return (
    <SafeAreaView>
      <View className="flex h-full px-4 pt-4">
        <View className="pb-4 mb-4 border-b-[1px] border-slate-300">
          <View className="flex flex-row items-center">
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
          <View>

          </View>
        </View>
      </View>
    </SafeAreaView>
  )
}

export default TransactionTab