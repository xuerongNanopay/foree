import { ScrollView, Text, View, Image } from 'react-native'
import React, { useState, useEffect } from 'react'
import { Link, router } from 'expo-router'
import { SafeAreaView } from 'react-native'

import { images } from '../../constants'
import FormField from '../../components/FormField'
import CustomButton from '../../components/CustomButton'
import { authPayload, authService } from '../../service'
import string_util from '../../util/string_util'
import * as Linking from 'expo-linking';

const SignUp = () => {
  const [errors, setErrors] = useState({})
  const [isError, setIsError] = useState(true);
  const url = Linking.useURL();

  const [form, setForm] = useState({
    email: '',
    password: '',
  })

  useEffect(() => {
    if (url) {
      const { hostname, path, queryParams } = Linking.parse(url);
      console.log(url)
      console.log(queryParams)
      if ( !!queryParams.referrerReference ) {
        setForm((form)=>({
          ...form,
          referrerReference: queryParams.referrerReference
        }))
      }
    }
  }, [url])

  useEffect(() => {
    async function validate() {
      try {
        await authPayload.SignUpScheme.validate(form, {abortEarly: false})
        setIsError(false)
        setErrors({})
      } catch (err) {
        let e = {}
        for ( let i of err.inner ) {
          e[i.path] =  e[i.path] ?? i.errors[0]
        }
        setErrors(e)
        setIsError(true)
      }
    }
    validate()
  }, [form])

  const [isSubmitting, setIsSubmitting] = useState(false)

  const submit = async () => {
    setIsSubmitting(true)
    try {
      console.log(form)
      const resp = await authService.signUp(string_util.trimStringInObject(form))
      if ( resp.status / 100 !== 2 ) {
        console.log("sign_up", resp.status, resp.data)
        return
      }
      router.replace("/verify_email")
    } catch (err) {
      console.error(err)
    } finally {
      setIsSubmitting(false)
    }
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
            <View className="rounded-lg border-2 border-secondary-100">
              <Link 
                href="/login" 
                className="text-lg text-secondary-100 font-psemibold p-1"
              >Sign In</Link>
            </View>
          </View>
          <Text className="text-2xl text-slate-700 mt-10 font-pbold">
            Create an account
          </Text>
          <FormField
            title="Email"
            value={form.email}
            handleChangeText={(e) => setForm({
              ...form,
              email:e.toLowerCase()
            })}
            containerStyles="mt-7"
            errorMessage={errors['email']}
            keyboardType="email-address"
          />
          <FormField
            title="Password"
            value={form.password}
            handleChangeText={(e) => setForm({
              ...form,
              password:e
            })}
            errorMessage={errors['password']}
            containerStyles="mt-7"
          />
          <CustomButton
            title="Sign Up"
            handlePress={submit}
            containerStyles="mt-7"
            disabled={isSubmitting || isError}
          />
        </View>
      </ScrollView>
    </SafeAreaView>
  )
}

export default SignUp