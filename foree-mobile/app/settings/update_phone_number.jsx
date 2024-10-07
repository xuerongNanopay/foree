import { View, Text, SafeAreaView, TouchableOpacity, Alert } from 'react-native'
import React, { useCallback, useEffect, useState } from 'react'
import FormField from '../../components/FormField'
import { router, useFocusEffect } from 'expo-router'
import { authPayload, authService } from '../../service'
import string_util from '../../util/string_util'

const FieldItem = ({
  title,
  value,
  handleChangeText,
  editable=true,
  keyboardType='ascii-capable',
  errorMessage,
}) => (
  <FormField
    title={title}
    value={value}
    handleChangeText={handleChangeText}
    keyboardType={keyboardType}
    containerStyles="mt-2"
    titleStyles="text-slate-600"
    inputContainerStyles={"h-6 pb-1"}
    variant='flat'
    editable={editable}
    errorMessage={errorMessage}
  />
)

const UpdatePhoneNumber = () => {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isError, setIsError] = useState(true)
  const [errors, setErrors] = useState({})
  const [form, setForm] = useState({
    phoneNumber: ''
  })

  useFocusEffect(useCallback(() => {
    const controller = new AbortController()
    const getUserDetail = async() => {
      try {
        const resp = await authService.getUserDetail({signal: controller.signal})
        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("get userDetail", resp.status, resp.data)
          router.replace("/personal_settings")
        } else {
          const userDetail = resp.data.data
          setForm({
            phoneNumber: userDetail.phoneNumber,
          })
        }
      } catch (e) {
        console.error("get userDetail", e, e.response, e.response?.status, e.response?.data)
        router.replace("/personal_settings")
      }
    }
    getUserDetail()
    return () => {
      controller.abort()
    }
  }, []))

  useEffect(() => {
    async function validate() {
      try {
        await authPayload.UpdatePhoneNumberScheme.validate(form, {abortEarly: false})
        setErrors({})
        setIsError(false)
      } catch (err) {
        let e = {}
        for ( let i of err.inner ) {
          e[i.path] =  e[i.path] ?? i.errors[0]
        }
        setIsError(true)
        setErrors(e)
      }
    }
    validate()
  }, [form])

  const submit = async () => {
    setIsSubmitting(true)
    try {
      const resp = await authService.updatePhone(string_util.trimStringInObject(form))
      if ( resp.status / 100 !== 2 ) {
        console.warn("update address", resp.status, resp.data)
        return
      }
      if ( router.canGoBack ) {
        router.back()
      } else {
        router.replace("/personal_settings")
      }
    } catch (err) {
      console.error("update address", err)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <SafeAreaView classname="h-full">
      <View className="h-full bg-slate-100 pt-2 px-3 flex">
        <View className="flex-1 mt-3">
          <FieldItem title="Phone Number" value={form.phoneNumber}
            errorMessage={errors['address1']}
            handleChangeText={(e) => setForm((form) => ({
              ...form,
              phoneNumber:e
            }))}
            keyboardType="phone-pad"
          />
        </View>
        <TouchableOpacity
          className={`mb-6 py-2 border-2 border-[#005a32] bg-[#c7e9c0] rounded-xl ${isSubmitting||isError ? 'opacity-50' : ''}`}
          onPress={() => {
            Alert.alert("Update Phone Number", "Are you sure?", [
              {text: 'Continue', onPress: () => {submit()}},
              {text: 'Cancel', onPress: () => {}},
            ])
          }}
        >
          <Text className="font-pbold text-lg text-[#005a32] text-center">Update</Text>
        </TouchableOpacity>
      </View>
    </SafeAreaView>
  )
}

export default UpdatePhoneNumber