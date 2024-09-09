import { View, Text, ScrollView } from 'react-native'
import React, { useState } from 'react'
import { router } from 'expo-router'
import { SafeAreaView } from 'react-native'

import CustomButton from '../../components/CustomButton'
import FormField from '../../components/FormField'
import { authService } from '../../service'

const ForgetPassword = () => {

  const [form, setForm] = useState({
    email: ''
  })

  const [isSubmitting, setIsSubmitting] = useState(false)

  const submit = async () => {
    setIsSubmitting(true)
    const resp = await authService.forgetPassword(form)
    console.log(resp)
    setIsSubmitting(false)
  }

  return (
    <SafeAreaView className="h-full">
      <ScrollView
        className="bg-slate-100"
        automaticallyAdjustKeyboardInsets
      >
        <View className="w-full mt-4">
          <Text className="text-lg font-pbold text-center m-4">Forget Password?</Text>
          <Text className="font-pregular text-center m-4">
            Enter the mail you used to create your account in order to reset your password.
          </Text>

          <View className="px-2">
            <FormField
              title="Email"
              value={form.email}
              handleChangeText={(e) => setForm({
                ...form,
                email:e
              })}
              containerStyles="mt-1"
              keyboardType="email-address"
            />
            <CustomButton
              title="Submit"
              handlePress={submit}
              containerStyles="mt-7"
              isLoading={isSubmitting}
            />
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  )
}

export default ForgetPassword