import { ScrollView, Text, View, Image } from 'react-native'
import React, { useState } from 'react'
import { Link } from 'expo-router'
import { SafeAreaView } from 'react-native-safe-area-context'

import { images } from '../../constants'
import FormField from '../../components/FormField'
import CustomButton from '../../components/CustomButton'

const SignIn = () => {
  const [form, setForm] = useState({
    email: '',
    password: ''
  })

  const [isSubmitting, setIsSubmitting] = useState(false)

  const submit = () => {
    setIsSubmitting(true)
    setTimeout(() => {
      setIsSubmitting(false)
    }, 2000);
  }

  return (
    //TODO: investigate why StatusBar.backgroundColor not working.
    // <SafeAreaView className="h-full">
    // <ScrollView>
    <SafeAreaView className="h-full bg-[#004d40]">
      <ScrollView className="bg-slate-100">
        <View className="w-full">
          <View className="px-4 mt-5">
            <Image
              source={images.logoSmall}
              resizeMode='contain'
              className="w-[36px] h-[36px]"
            />
            <View>
              <Text className="mt-6 text-secondary text-left font-psemibold text-xl">
                Sign-up & receive $44 to try the fastest global transfer
              </Text>
              <Text className="mt-4 text-secondary text-left font-bold">
                &#10003; $0 fees and best FX rates
              </Text>
              <Text className="mt-4 text-secondary text-left font-bold">
                &#10003; Transfer to 35+ USA banks
              </Text>
              <Text className="mt-4 text-secondary text-left font-bold">
                &#10003; Cash pick-ups from 1500+ USA branches
              </Text>
            </View>
          </View>
          <View>
            <View>
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
                title="Sign In"
                handlePress={submit}
                containerStyles="mt-7"
                isLoading={isSubmitting}
              />
              <View className="justify-center items-center pt-4 flex-row gap-2">
                <Text className="text-lg text-gray-100 font-pregular">
                  Don't have account?
                </Text>
                <Link href="/sign_up" className='"text-lg font-psemibold text-secondary'>Sign Up</Link>
              </View>
            </View>
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  )
}

export default SignIn