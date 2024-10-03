import { View, Text, ScrollView } from 'react-native'
import React, { useState, useEffect } from 'react'
import { router, useLocalSearchParams } from 'expo-router'
import { SafeAreaView } from 'react-native'
import FormField from '../../components/FormField'
import CustomButton from '../../components/CustomButton'
import { authService, authPayload } from '../../service'
import string_util from '../../util/string_util'

const ForgetPasswordUpdate = () => {
  const { email, retrieveCode } = useLocalSearchParams()
  const [isError, setIsError] = useState(true)
  const [errors, setErrors] = useState({})
  const [form, setForm] = useState({
    email,
    retrieveCode,
    newPassword: '',
  })

  useEffect(() => {
    async function validate() {
      try {
        await authPayload.ForgetPasswdUpdateScheme.validate(form, {abortEarly: false})
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
      const resp = await authService.forgetPasswordUpdate(string_util.trimStringInObject(form))
      if ( resp.status / 100 !== 2 ) {
        console.info("forget_password_verify", resp.status, resp.data)
        return
      }
      router.replace({ pathname: `/login`});
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
        <View className="w-full mt-4 px-2">
          <Text className="text-lg font-pbold text-center m-4">Renew Your Password</Text>
          <Text className="font-pregular text-center m-4">
            Please provide new password for login.
          </Text>
          <FormField
            title="New Password"
            value={form.newPassword}
            handleChangeText={(e) => setForm({
              ...form,
              newPassword:e
            })}
            errorMessage={errors['newPassword']}
            containerStyles="mt-7"
          />
          <CustomButton
            title="Update"
            handlePress={submit}
            containerStyles="mt-7"
            disabled={isSubmitting || isError}
          />
        </View>
      </ScrollView>
    </SafeAreaView>
  )
}

export default ForgetPasswordUpdate