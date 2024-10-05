import { View, Text, SafeAreaView, TouchableOpacity, Alert } from 'react-native'
import React, { useState } from 'react'
import FormField from '../../components/FormField'

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
  const [errors, setErrors] = useState({})
  const [form, setForm] = useState({
    phoneNumber: ''
  })


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
          className="mb-6 py-2 border-2 border-[#005a32] bg-[#c7e9c0] rounded-xl"
          onPress={() => {
            Alert.alert("Update Phone Number", "Are you sure?", [
              {text: 'Continue', onPress: () => {console.log("TODO: update phone number")}},
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