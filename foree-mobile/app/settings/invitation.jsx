import { View, Text, SafeAreaView, TouchableOpacity, Share } from 'react-native'
import React, { useCallback, useState } from 'react'
import { useFocusEffect } from 'expo-router'
import { authService } from '../../service'

const Invitation = () => {
  const [userReference, setUserReference] = useState('')
  const [linkScheme, setLinkScheme] = useState('exp://127.0.0.1:8081/--/sign_up')
  useFocusEffect(useCallback(() => {
    const controller = new AbortController()
    const getUserExtra = async() => {
      try {
        const resp = await authService.getUserExtra({signal: controller.signal})
        if ( resp.status / 100 != 2 && !resp?.data?.data ) {
          console.error("get userExtra", resp.status, resp.data)
          return
        }
        const userExtra = resp.data.data
        setUserReference(userExtra.userReference)
        console.log(userExtra)
      } catch (e) {
        console.error("get userDetail", e, e.response, e.response?.status, e.response?.data)
      }
    }
    getUserExtra()
  }, []))

  const onShare = async() => {
    try {
      const invitationLink = `${linkScheme}?referrerReference=${userReference}`
      console.log(invitationLink)
      const result = await Share.share({
        message:
          `Here's a link to try Foree Remittance, the fastest way to send money to Pakistan!\n${invitationLink}`,
      });
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