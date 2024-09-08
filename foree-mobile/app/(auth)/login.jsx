import { ScrollView, Text, View, Image } from 'react-native'
import React, { useState, useEffect } from 'react'
import { Link, router, useNavigation } from 'expo-router'
import { SafeAreaView } from 'react-native'

import { images } from '../../constants'
import FormField from '../../components/FormField'
import CustomButton from '../../components/CustomButton'

const Login = () => {
  const navigation = useNavigation()
  useEffect(() => {
    console.log("TODO: clean token")
  }, [navigation])

  const [form, setForm] = useState({
    email: '',
    password: ''
  })

  const [isSubmitting, setIsSubmitting] = useState(false)

  const submit = () => {
    setIsSubmitting(true)
    setTimeout(() => {
      setIsSubmitting(false)
      router.push('/verify_email')
    }, 1000);
  }

  return (
    <SafeAreaView className="h-full">
      <ScrollView
        className="bg-slate-100"
        automaticallyAdjustKeyboardInsets
      >
        <View className="w-full">
          <View className="px-4 mt-5">
            <View className="flex-row items-center justify-between">
              <Image
                source={images.logoSmall}
                resizeMode='contain'
                className="w-[36px] h-[36px]"
              />
              <View className="rounded-lg border-2 border-secondary-100">
                <Link
                  href="/sign_up" 
                  className="text-lg text-secondary-100 font-psemibold p-1"
                  disabled={isSubmitting}
                >Sign Up</Link>
              </View>
            </View>
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
          <View className="bg-slate-100 px-2 mt-4">
            <Text className="font-pbold text-xl text-slate-700">Wellcome back</Text>
            <FormField
              title="Email"
              value={form.email}
              handleChangeText={(e) => setForm({
                ...form,
                email:e
              })}
              containerStyles="mt-4"
              keyboardType="email-address"
            />
            <FormField
              title="Password"
              value={form.password}
              handleChangeText={(e) => setForm({
                ...form,
                password:e
              })}
              containerStyles="mt-7"
            />
            <Link 
                href="/forget_password" 
                className="text-slate-500 p-2"
              >Forget Password?</Link>
            <CustomButton
              title="Sign In"
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

export default Login