import { View, Text, ScrollView } from 'react-native'
import React from 'react'
import { SafeAreaView } from 'react-native'

const UpdatePassword = () => {
  return (
    <SafeAreaView className="h-full">
      <ScrollView
        className="bg-slate-800"
        automaticallyAdjustKeyboardInsets
      >
        <View className="w-full mt-4">
          <Text className="text-lg font-pbold text-center m-4">Update Your Password</Text>
          <Text className="font-pregular text-center m-4">
            Please provide new password for login.
          </Text>
        </View>
      </ScrollView>
    </SafeAreaView>
  )
}

export default UpdatePassword