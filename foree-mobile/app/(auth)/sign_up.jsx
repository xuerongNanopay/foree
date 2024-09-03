import { ScrollView, Text, View, Image } from 'react-native'
import React, { useState } from 'react'
import { Link } from 'expo-router'
import { SafeAreaView } from 'react-native-safe-area-context'

import { images } from '../../constants'
import FormField from '../../components/FormField'
import CustomButton from '../../components/CustomButton'

const SignUp = () => {
  const [form, setForm] = useState({
    email: '',
    password: '',
    rePassword: ''
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
        <View className="w-full px-2 my-6">
          <View className="flex-row items-center justify-between">
            <Image
              source={images.logoSmall}
              resizeMode='contain'
              className="w-[36px] h-[36px]"
            />
            <View className="bg-secondary-200 rounded-lg">
              <Link 
                href="/sign_in" 
                className="text-lg font-psemibold text-white p-2"
              >Sign In</Link>
            </View>
          </View>
          <Text className="text-2xl text-slate-700 text-semibold mt-10 font-psemibold">
            Create an account
          </Text>
          <FormField
            title="Email"
            value={form.email}
            handleChangeText={(e) => setForm({
              ...form,
              email:e
            })}
            otherStyles="mt-7"
            keyboardType="email-address"
          />
          <FormField
            title="Password"
            value={form.password}
            handleChangeText={(e) => setForm({
              ...form,
              password:e
            })}
            otherStyles="mt-7"
          />
          <CustomButton
            title="Sign Up"
            handlePress={submit}
            containerStyles="mt-7"
            isLoading={isSubmitting}
          />
        </View>
      </ScrollView>
    </SafeAreaView>
  )
}

export default SignUp