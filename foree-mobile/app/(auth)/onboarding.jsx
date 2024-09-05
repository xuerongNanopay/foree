import { View, Text, ScrollView } from 'react-native'
import { SafeAreaView } from 'react-native'
import React, { useState } from 'react'
import MultiStepForm from '../../components/MultiStepForm'
import FormField from '../../components/FormField'

const FieldItem = ({
  title,
  value,
  handleChangeText,
  keyboardType='ascii-capable'
}) => (
  <FormField
    title={title}
    value={value}
    handleChangeText={handleChangeText}
    keyboardType={keyboardType}
    otherStyles="mt-2"
  />
)

const ReviewItem = ({
  title,
  value,
}) => (
  <FormField
    title={title}
    value={value}
    handleChangeText={() => {}}
    keyboardType="ascii-capable"
    otherStyles="mt-2"
    variant="flat"
    inputStyles="text-slate-500"
    inputContainerStyles="border-slate-700 h-7"
    editable={false}
  />
)

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
    dob: '',
    pob: '',
    nationality: '',
    identificationType: '',
    identificationValue: '',
  })

  const submit = () => {
    setIsSubmitting(true)
    setTimeout(() => {
      setIsSubmitting(false)
      console.log(form)
      // router.replace("/verify_email")
    }, 1000);
  }

  const [isSubmitting, setIsSubmitting] = useState(false)

  const NameField = () => (
    <View>
      <Text className="text-lg font-pbold text-center m-4">Let's Get to Know You!</Text>
      <Text className="font-pregular text-center mb-4">
        Please enter your full legal name so we can begin setting up your account
      </Text>

      <FieldItem title="First Name" value={form.firstName}
        handleChangeText={(e) => setForm({
          ...form,
          firstName:e
        })}
      />
      <FieldItem title="Middle Name" value={form.middleName}
        handleChangeText={(e) => setForm({
          ...form,
          middleName:e
        })}
      />
      <FieldItem title="Last Name" value={form.lastName}
        handleChangeText={(e) => setForm({
          ...form,
          lastName:e
        })}
      />
    </View>
  )

  const AddressField = () => (
    <View>
      <Text className="text-lg font-pbold text-center m-4">Your Residential Address and Phone Number</Text>
      <Text className="font-pregular text-center mb-4">
        We require this information to continue setting up your Foree Remittance account
      </Text>

      <FieldItem title="Address Line 1" value={form.addressLine1}
        handleChangeText={(e) => setForm({
          ...form,
          addressLine1:e
        })}
      />
      <FieldItem title="Address Line 2" value={form.addressLine2}
        handleChangeText={(e) => setForm({
          ...form,
          addressLine2:e
        })}
      />
      <FieldItem title="City" value={form.city}
        handleChangeText={(e) => setForm({
          ...form,
          city:e
        })}
      />
      <FieldItem title="Province" value={form.province}
        handleChangeText={(e) => setForm({
          ...form,
          province:e
        })}
      />
      <FieldItem title="Country" value={form.country}
        handleChangeText={(e) => setForm({
          ...form,
          country:e
        })}
      />
      <FieldItem title="Postal Code" value={form.postalCode}
        handleChangeText={(e) => setForm({
          ...form,
          postalCode:e
        })}
      />
      <FieldItem title="Phone Number" value={form.phoneNumber}
        handleChangeText={(e) => setForm({
          ...form,
          phoneNumber:e
        })}
        keyboardType="phone-pad"
      />
    </View>
  )

  const PersonalDetailField = () => (
    <View>
      <Text className="text-lg font-pbold text-center m-4">Personal Details</Text>
      <Text className="font-pregular text-center mb-4">
        Almost done! Infomation below is requested by xxxxx xxxxx of xxxxxxx, our Foree Remittance payout parter, inorder to process your transfers under ...... regulatory guidelines
      </Text>
      <FieldItem title="Date of Birth" value={form.dob}
        handleChangeText={(e) => setForm({
          ...form,
          dob:e
        })}
      />
      <FieldItem title="Place of Birth" value={form.pob}
        handleChangeText={(e) => setForm({
          ...form,
          pob:e
        })}
      />
      <FieldItem title="Nationality" value={form.nationality}
        handleChangeText={(e) => setForm({
          ...form,
          nationality:e
        })}
      />
      <FieldItem title="Identification Document Type" value={form.identificationType}
        handleChangeText={(e) => setForm({
          ...form,
          identificationType:e
        })}
      />
      <FieldItem title="Identification Number" value={form.identificationValue}
        handleChangeText={(e) => setForm({
          ...form,
          identificationValue:e
        })}
      />
    </View>
  )

  const Review = () => (
    <View>
      <Text className="text-lg font-pbold text-center m-4">Review</Text>
      <ReviewItem title="First Name" value={form.firstName}/>
      <ReviewItem title="Middle Name" value={form.middleName}/>
      <ReviewItem title="Last Name" value={form.lastName}/>
      <ReviewItem title="Address Line 1" value={form.addressLine1}/>
      <ReviewItem title="Address Line 2" value={form.addressLine2}/>
      <ReviewItem title="City" value={form.city}/>
      <ReviewItem title="Province" value={form.province}/>
      <ReviewItem title="Country" value={form.country}/>
      <ReviewItem title="Postal Code" value={form.postalCode}/>
      <ReviewItem title="Phone Number" value={form.phoneNumber}/>
      <ReviewItem title="Date of Birth" value={form.dob}/>
      <ReviewItem title="Place of Birth" value={form.pob}/>
      <ReviewItem title="Nationality" value={form.nationality}/>
      <ReviewItem title="Identification Document Type" value={form.identificationType}/>
      <ReviewItem title="Identification Number" value={form.identificationValue}/>
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
    },
    {
      formView: PersonalDetailField,
      canGoNext: () => {
        return true
      }
    },
    {
      formView: Review,
      canGoNext: () => {
        return true
      }
    }
  ]

  return (
    <SafeAreaView className="bg-slate-100">
      <MultiStepForm
        steps={() => OnboardingFlow}
        onSumbit={submit}
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