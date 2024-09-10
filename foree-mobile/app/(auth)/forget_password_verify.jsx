import { View, Text, ScrollView } from 'react-native'
import React, { useState } from 'react'
import { router, useLocalSearchParams } from 'expo-router'
import { SafeAreaView } from 'react-native'
import FormField from '../../components/FormField'
import CustomButton from '../../components/CustomButton'
import { authService, authPayload } from '../../service'

const ForgetPasswordVerify = () => {
  const { preForm } = useLocalSearchParams()

  const [errors, setErrors] = useState({});
  const [form, setForm] = useState({
    email: preForm.email,
    retrieveCode,
  })

  useEffect(() => {
    async function validate() {
      try {
        await authPayload.ForgetPasswordVerifycheme.validate(form, {abortEarly: false})
        setErrors({})
      } catch (err) {
        let e = {}
        for ( let i of err.inner ) {
          e[i.path] =  e[i.path] ?? i.errors[0]
        }
        setErrors(e)
      }
    }
    validate()
  }, [form])

  const [isSubmitting, setIsSubmitting] = useState(false)

  const submit = async () => {
    setIsSubmitting(true)
    try {
      const resp = await authService.forgetPasswordVerify(form)
      if ( resp.status / 100 !== 2 ) {
        console.info("forget_password_verify", resp.status, resp.data)
        return
      }
      router.replace({ pathname: `/forget_password_update`, params: form });
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
          <Text className="text-lg font-pbold text-center m-4">Verify Your Account</Text>
          <Text className="font-pregular text-center m-4">
          We have sent a verification code to your email. Please enter the code below to confirm that this account belongs to you.
          </Text>
          <FormField
              value={form.retrieveCode}
              handleChangeText={(e) => setForm({
                ...form,
                retrieveCoderetrieveCode:e
              })}
              errorMessage={errors['retrieveCode']}
              inputStyles="text-center"
              variant="flat"
              containerStyles="mt-1"
              keyboardType="numeric"

            />
          <CustomButton
            title="Verify"
            handlePress={submit}
            containerStyles="mt-7"
            disabled={isSubmitting}
          />
        </View>
      </ScrollView>
    </SafeAreaView>
  )
}

export default ForgetPasswordVerify