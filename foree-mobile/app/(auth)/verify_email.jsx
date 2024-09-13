import { View, Text, ScrollView, Button, Alert } from 'react-native'
import React, { useState, useEffect } from 'react'
import { router } from 'expo-router'
import { SafeAreaView } from 'react-native'
import { authPayload, authService } from '../../service'

import CustomButton from '../../components/CustomButton'
import FormField from '../../components/FormField'
import string_util from '../../util/string_util'


const VerifyEmail = () => {
  const [errors, setErrors] = useState({})
  const [isError, setIsError] = useState(true);

  const [form, setForm] = useState({
    code: ''
  })

  useEffect(() => {
    async function validate() {
      try {
        await authPayload.VerifyEmailScheme.validate(form, {abortEarly: false})
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
      const resp = await authService.verifyEmail(string_util.trimStringInObject(form))
      if ( resp.status / 100 !== 2 ) {
        console.log("verify_email", resp.status, resp.data)
        return
      }
      router.replace("/onboarding")
    } catch (err) {
      console.error(err)
    } finally {
      setIsSubmitting(false)
    }
  }

  const resendSubmit = async () => {
    setIsSubmitting(true)
    try {
      const resp = await authService.resendCode()
      if ( resp.status / 100 !== 2 ) {
        console.log("recend_code", resp.status, resp.data)
        return
      }
      Alert.alert("New code sent", "please check your email", [
        {text: 'OK', onPress: () => {}},
      ])
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
          <Text className="text-lg font-pbold text-center m-4">Let's Verify Your Email Address</Text>
          <Text className="font-pregular text-center mt-4">
            We have sent a verification code to your email. Please enter the code below to confirm that this account belongs to you.
          </Text>

          <View className="px-2 mt-4">
            <FormField
              value={form.code}
              handleChangeText={(e) => setForm({
                ...form,
                code:e
              })}
              errorMessage={errors['code']}
              inputStyles="text-center"
              variant="flat"
              containerStyles="mt-1"
              keyboardType="numeric"

            />
            <Text className="text-sm text-slate-600">Please enter the 6-digit code sent to your email</Text>
            <CustomButton
              title="Submit"
              handlePress={submit}
              containerStyles="mt-7"
              disabled={isSubmitting || isError}
            />
            <View className="flex-row items-center justify-center mt-4">
              <Text>Do not receive code?</Text>
              <Button
                onPress={resendSubmit}
                title="Resend"
                color="#26a69a"
                disabled={isSubmitting}
              />
            </View>
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  )
}

export default VerifyEmail