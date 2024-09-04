import { View, Text, ScrollView } from 'react-native'
import React, { useState } from 'react'
import { router, useLocalSearchParams } from 'expo-router'
import { SafeAreaView } from 'react-native'
import FormField from '../../components/FormField'
import CustomButton from '../../components/CustomButton'

const UpdatePassword = () => {
  const { token } = useLocalSearchParams()

  const [form, setForm] = useState({
    token,
    password: '',
  })

  const [isSubmitting, setIsSubmitting] = useState(false)

  const submit = () => {
    setIsSubmitting(true)
    setTimeout(() => {
      setIsSubmitting(false)
      router.replace("/login")
    }, 2000);
  }

  return (
    <SafeAreaView className="h-full">
      <ScrollView
        className="bg-slate-100"
        automaticallyAdjustKeyboardInsets
      >
        <View className="w-full mt-4 px-2">
          <Text className="text-lg font-pbold text-center m-4">Renew Your Password</Text>
          <Text className="font-pregular text-center m-4">
            Please provide new password for login.
          </Text>
          <FormField
            title="Password"
            value={form.email}
            handleChangeText={(e) => setForm({
              ...form,
              password:e
            })}
            otherStyles="mt-7"
            keyboardType="email-address"
          />
          <CustomButton
            title="Update"
            handlePress={submit}
            containerStyles="mt-7"
            isLoading={isSubmitting}
          />
        </View>
      </ScrollView>
    </SafeAreaView>
  )
}

export default UpdatePassword