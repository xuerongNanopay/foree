import { View, Text, SafeAreaView, TouchableOpacity, Alert } from 'react-native'
import React, { useCallback, useEffect, useState } from 'react'
import ModalSelect from '../../components/ModalSelect'
import Regions from '../../constants/region'
import Countries from '../../constants/country'
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

const SelectProvinceItem = (province) => (
  <Text className="font-pregular py-3 text-xl">
    {province["name"]}
  </Text>
)

const UpdateAddress = () => {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isError, setIsError] = useState(true)
  const [errors, setErrors] = useState({})
  const [form, setForm] = useState({
    address1: '',
    address2: '',
    city: '',
    province: '',
    country: '',
    postalCode: ''
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
            address1: userDetail.address1,
            address2: userDetail.address2,
            city: userDetail.city,
            province: userDetail.province,
            country: userDetail.country,
            postalCode: userDetail.postalCode,
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
        await authPayload.UpdateAddressScheme.validate(form, {abortEarly: false})
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
      const resp = await authService.updateAddress(string_util.trimStringInObject(form))
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
        <View className="flex-1">
          <FieldItem title="Address Line 1" value={form.address1}
            errorMessage={errors['address1']}
            handleChangeText={(e) => setForm((form) => ({
              ...form,
              address1:e
            }))}
          />
          <FieldItem title="Address Line 2" value={form.address2}
            errorMessage={errors['address2']}
            handleChangeText={(e) => setForm((form) => ({
              ...form,
              address2:e
            }))}
          />
          <FieldItem title="City" value={form.city}
            errorMessage={errors['city']}
            handleChangeText={(e) => setForm((form) => ({
              ...form,
              city:e
            }))}
          />
          <ModalSelect
            title="Province"
            titleStyles="text-slate-600"
            modalTitle="select a province"
            errorMessage={errors['province']}
            containerStyles="mt-2"
            value={form.province}
            variant='flat'
            inputContainerStyles={"h-6 pb-1"}
            searchKey="name"
            keyExtractor="name"
            showExtractor="name"
            valueExtractor="isoCode"
            listView={SelectProvinceItem}
            list={Object.values(Regions[form.country] ?? [])}
            onPress={(o) => {
              setForm((form) => ({
                ...form,
                province: o
              }))
            }}
            placeholder="select a province"
          />
          <FieldItem title="Country" value={`${Countries[form.country]?.unicodeIcon} ${Countries[form.country]?.name}`}
            errorMessage={errors['country']}
            handleChangeText={(e) => setForm((form) => ({
              ...form,
              country:e
            }))}
            editable={false}
          />
          <FieldItem title="Postal Code" value={form.postalCode}
            errorMessage={errors['postalCode']}
            handleChangeText={(e) => setForm((form) => ({
              ...form,
              postalCode:e
            }))}
          />
        </View>
        <TouchableOpacity
          className={`mb-6 py-2 border-2 border-[#005a32] bg-[#c7e9c0] rounded-xl ${isSubmitting||isError34e ? 'opacity-50' : ''}`}
          disabled={isSubmitting||isError}
          onPress={() => {
            Alert.alert("Update Address", "Are you sure?", [
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

export default UpdateAddress