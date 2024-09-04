import { View, Text, ScrollView } from 'react-native'
import { SafeAreaView } from 'react-native'
import React, { useState } from 'react'
import MultiStepForm from '../../components/MultiStepForm'


const Onboarding = () => {

  const [form, setForm] = useState({
    email: ''
  })

  const [isSubmitting, setIsSubmitting] = useState(false)

  const OnboardingFlow = [
    {
      formView: NameField,
      canGoNext: () => {
        return true
      }
    }
  ]

  const NameField = () => (
    <View>
      <Text className="text-lg font-pbold text-center m-4">Let's Get to Know You!</Text>
      <Text className="font-pregular text-center m-4">
        Please enter your full legal name so we can begin setting up your account
      </Text>
    </View>
  )

  return (
    <SafeAreaView className="bg-slate-100">
      <MultiStepForm
        steps={() => [NameField]}
      />
      {/* <ScrollView
        className="bg-slate-100"
        automaticallyAdjustKeyboardInsets
      >
        <MultiStepForm/>
      </ScrollView> */}
    </SafeAreaView>
  )
}

export default Onboarding