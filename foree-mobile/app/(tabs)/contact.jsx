import { View, Text } from 'react-native'
import React from 'react'
import { SafeAreaView } from 'react-native-safe-area-context'

const Contact = () => {
  return (
    <SafeAreaView className="border-2 border-red-600">
      <View className="px-4 pt-4">
        <View className="pb-2 border-b-[1px] border-slate-400">
          <Text className="font-pbold text-2xl">Contacts</Text>
        </View>
      </View>
    </SafeAreaView>
  )
}

export default Contact