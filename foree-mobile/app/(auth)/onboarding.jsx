import { View, Text, ScrollView } from 'react-native'
import { SafeAreaView } from 'react-native'
import React, { useState } from 'react'
import MultiStepForm from '../../components/MultiStepForm'
import FormField from '../../components/FormField'


const Onboarding = () => {

  const [form, setForm] = useState({
    firstName: '',
    middleName: '',
    lastName: '',
    addressLine1: '',
    addressLine2: '',
    city: '',
    province: '',
    country: 'Canada',
    postalCode: '',
    phoneNumber: '',
  })

  const [isSubmitting, setIsSubmitting] = useState(false)

  const NameField = () => (
    <View>
      <Text className="text-lg font-pbold text-center m-4">Let's Get to Know You!</Text>
      <Text className="font-pregular text-center mb-4">
        Please enter your full legal name so we can begin setting up your account
      </Text>

      <FormField
        title="First Name"
        value={form.firstName}
        handleChangeText={(e) => setForm({
          ...form,
          firstName:e
        })}
        keyboardType="ascii-capable"
        otherStyles="mt-2"
      />
      <FormField
        title="Middle Name"
        value={form.middleName}
        handleChangeText={(e) => setForm({
          ...form,
          middleName:e
        })}
        keyboardType="ascii-capable"
        otherStyles="mt-2"
      />
      <FormField
        title="Last Name"
        value={form.lastName}
        handleChangeText={(e) => setForm({
          ...form,
          lastName:e
        })}
        otherStyles="mt-2"
        keyboardType="ascii-capable"
      />
    </View>
  )

  const AddressField = () => (
    <View>
      <Text className="text-lg font-pbold text-center m-4">Your Residential Address and Phone Number</Text>
      <Text className="font-pregular text-center mb-4">
        We require this information to continue setting up your Foree Remittance account
      </Text>

      <FormField
        title="Address Line 1"
        value={form.addressLine1}
        handleChangeText={(e) => setForm({
          ...form,
          addressLine1:e
        })}
        keyboardType="ascii-capable"
        otherStyles="mt-2"
      />
      <FormField
        title="Address Line 2"
        value={form.addressLine2}
        handleChangeText={(e) => setForm({
          ...form,
          addressLine2:e
        })}
        keyboardType="ascii-capable"
        otherStyles="mt-2"
      />
      <FormField
        title="City"
        value={form.city}
        handleChangeText={(e) => setForm({
          ...form,
          city:e
        })}
        keyboardType="ascii-capable"
        otherStyles="mt-2"
      />
      <FormField
        title="Province"
        value={form.province}
        handleChangeText={(e) => setForm({
          ...form,
          province:e
        })}
        keyboardType="ascii-capable"
        otherStyles="mt-2"
      />
      <FormField
        title="Country"
        value={form.country}
        handleChangeText={(e) => setForm({
          ...form,
          country:e
        })}
        keyboardType="ascii-capable"
        otherStyles="mt-2"
      />
      <FormField
        title="Postal Code"
        value={form.postalCode}
        handleChangeText={(e) => setForm({
          ...form,
          postalCode:e
        })}
        otherStyles="mt-2"
        keyboardType="ascii-capable"
      />
      <FormField
        title="Phone Number"
        value={form.phoneNumber}
        handleChangeText={(e) => setForm({
          ...form,
          phoneNumber:e
        })}
        otherStyles="mt-2"
        keyboardType="ascii-capable"
      />
    </View>
  )

  const OnboardingFlow = [
    {
      formView: NameField,
      canGoNext: () => {
        return true
      }
    },
    {
      formView: AddressField,
      canGoNext: () => {
        return true
      }
    }
  ]

  return (
    <SafeAreaView className="bg-slate-100">
      <MultiStepForm
        steps={() => OnboardingFlow}
      />
      {/* <ScrollView
        className="bg-slate-100"
        automaticallyAdjustKeyboardInsets
      >
        <MultiStepForm/>
      </ScrollView> */}
    </SafeAreaView>
  )
}

export default Onboarding