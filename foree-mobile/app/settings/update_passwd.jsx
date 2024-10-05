import { View, Text, SafeAreaView, Alert, TouchableOpacity } from 'react-native'
import React, { useEffect, useState } from 'react'
import FormField from '../../components/FormField'
import CustomButton from '../../components/CustomButton'

const UpdatePasswd = () => {
  const [errors, setErrors] = useState({})
  const [isError, setIsError] = useState(true)
  const [form, setForm] = useState({
    oldPassword: '',
    newPassword: '',
  })
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    async function validate() {
      try {
        // await authPayload.ForgetPasswdUpdateScheme.validate(form, {abortEarly: false})
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


  const submit = async () => {
    setIsSubmitting(true)
    try {

    } catch (err) {
      console.error(err)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <SafeAreaView className="h-full">
      <View className="h-full mt-4 px-2 flex">
        <View className="flex-1"> 
          <Text className="text-lg font-pbold text-center m-4">Update Your Password</Text>
          <Text className="font-pregular text-center m-4">
            Please provide new password for login.
          </Text>
          <FormField
            title="Old Password"
            titleStyles="text-slate-600"
            value={form.oldPassword}
            isPassword={true}
            handleChangeText={(e) => setForm({
              ...form,
              oldPassword:e
            })}
            errorMessage={errors['oldPassword']}
            containerStyles="mt-7"
          />
          <FormField
            title="New Password"
            titleStyles="text-slate-600"
            value={form.newPassword}
            isPassword={true}
            handleChangeText={(e) => setForm({
              ...form,
              newPassword:e
            })}
            errorMessage={errors['newPassword']}
            containerStyles="mt-7"
          />
        </View>
        <TouchableOpacity
          className="mb-10 py-2 border-2 border-[#005a32] bg-[#c7e9c0] rounded-xl"
          onPress={() => {
            Alert.alert("Change password", "Are you sure?", [
              {text: 'Continue', onPress: () => {console.log("TODO: close account")}},
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

export default UpdatePasswd