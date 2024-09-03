import { View, Text, ScrollView, Image } from 'react-native'
import React, { useState } from 'react'
import { SafeAreaView } from 'react-native-safe-area-context'

import { images } from '../../constants'
import CustomButton from '../../components/CustomButton'
import FormField from '../../components/FormField'

const ForgetPassword = () => {
  const [form, setForm] = useState({
    email: ''
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
          <Text className="text-lg font-psemibold text-center m-4">Forget Password?</Text>
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
              otherStyles="mt-1"
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