import { View, Text, ScrollView, Alert } from 'react-native'
import React, { useState, useEffect } from 'react'
import { router } from 'expo-router'
import { SafeAreaView } from 'react-native'

import CustomButton from '../../components/CustomButton'
import FormField from '../../components/FormField'
import { authService, authPayload } from '../../service'

const ForgetPassword = () => {
  const [errors, setErrors] = useState({})
  const [isError, setIsError] = useState(true)
  const [form, setForm] = useState({
    email: ''
  })

  useEffect(() => {
    async function validate() {
      try {
        await authPayload.ForgetPasswdScheme.validate(form, {abortEarly: false})
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
      const resp = await authService.forgetPassword(form)
      if ( resp.status / 100 !== 2 ) {
        console.info("forget_password", resp.status, resp.data)
        return
      }
      router.replace({ pathname: `/forget_password_verify`, params: form });
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
                email:e.toLowerCase()
              })}
              containerStyles="mt-1"
              keyboardType="email-address"
              errorMessage={errors['email']}
            />
            <CustomButton
              title="Retrieve"
              handlePress={submit}
              containerStyles="mt-7"
              disabled={isSubmitting || isError}
            />
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  )
}

export default ForgetPassword