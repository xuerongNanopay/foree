import { View, Text, SafeAreaView, TouchableOpacity, Share } from 'react-native'
import React from 'react'

const Invitation = () => {
  const onShare = async() => {
    console.log("TODO: share")
    try {
      const result = await Share.share({
        message:
          `Here's a link to try Foree Remittance, the fastest way to send money to Pakistan!\n${'url'}`,
      });
      console.log(result)
    } catch (e) {
      console.error("invitation share", e)
    }
  }

  return (
    <SafeAreaView className="h-full">
      <View className="px-2">
        <Text className="mt-6 font-pbold text-lg text-center text-[#005a32]">Share Foree Remittance With a Friend!</Text>
        <Text className="mt-4 text-center text-slate-800">You both get payment credits when you fefer a friend to Foree Remittance.</Text>
        <View className="mt-6">
          <Text className="font-pbold text-lg text-[#005a32] text-center">You get $20.00</Text>
          <Text className="mt-2 text-slate-600 text-center">When your referral completes their first transaction</Text>
        </View>
        <View className="mt-6">
          <Text className="font-pbold text-lg text-[#005a32] text-center">They get $20.00</Text>
          <Text className="mt-2 text-slate-600 text-center">On their first transaction when they sign-up using your link</Text>
        </View>

        <TouchableOpacity
          className="mt-4 py-2 bg-[#005a32] rounded-lg"
          onPress={onShare}
        >
          <Text className="text-white text-center font-pbold text-lg">Share</Text>
        </TouchableOpacity>
      </View>
    </SafeAreaView>
  )
}

export default Invitation