import { View, Text, SafeAreaView } from 'react-native'
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
      <View className="w-full mt-4 px-2">
        <Text className="text-lg font-pbold text-center m-4">Update Your Password</Text>
        <Text className="font-pregular text-center m-4">
          Please provide new password for login.
        </Text>
        <FormField
          title="Old Password"
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
          value={form.newPassword}
          isPassword={true}
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
    </SafeAreaView>
  )
}

export default UpdatePasswd