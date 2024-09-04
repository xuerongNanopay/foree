import { View, Text, ScrollView, Image, Button } from 'react-native'
import React, { useState } from 'react'
import { SafeAreaView } from 'react-native-safe-area-context'

import { images } from '../../constants'
import CustomButton from '../../components/CustomButton'
import FormField from '../../components/FormField'


const VerifyEmail = () => {
  const [form, setForm] = useState({
    verifyCode: ''
  })

  const [isSubmitting, setIsSubmitting] = useState(false)

  const submit = () => {
    setIsSubmitting(true)
    setTimeout(() => {
      setIsSubmitting(false)
    }, 2000);
  }

  return (
    <SafeAreaView className="h-full">
      <ScrollView
        className="bg-slate-100"
        automaticallyAdjustKeyboardInsets
      >
        <View className="w-full mt-4">
          <View className="flex-row items-center justify-center">
            <Image
              source={images.logoSmall}
              resizeMode='contain'
              className="w-[36px] h-[36px]"
            />
          </View>
          <Text className="text-lg font-pbold text-center m-4">Let's Verify Your Email Address</Text>
          <Text className="font-pregular text-center mt-4">
            We have sent a verification code to your email. Please enter the code below to confirm that this account belongs to you.
          </Text>

          <View className="px-2 mt-4">
            <FormField
              value={form.verifyCode}
              handleChangeText={(e) => setForm({
                ...form,
                verifyCode:e
              })}
              inputStyles="text-center"
              variant="flat"
              otherStyles="mt-1"
              keyboardType="numeric"

            />
            <Text className="text-sm text-slate-600">Please enter the 6-digit code sent to your email</Text>
            <CustomButton
              title="Submit"
              handlePress={submit}
              containerStyles="mt-7"
              isLoading={isSubmitting}
            />
            <View className="flex-row items-center justify-center mt-4">
              <Text>Do not receive code?</Text>
              <Button
                onPress={submit}
                title="Resend"
                color="#26a69a"
                isSubmitting={isSubmitting}
              />
            </View>
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  )
}

export default VerifyEmail