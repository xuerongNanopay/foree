import { View, Text, ScrollView } from 'react-native'
import { SafeAreaView } from 'react-native'
import React from 'react'

const Onboarding = () => {
  return (
    <SafeAreaView className="h-full">
      <ScrollView
        className="bg-slate-100"
        automaticallyAdjustKeyboardInsets
      >
      </ScrollView>
    </SafeAreaView>
  )
}

export default Onboarding