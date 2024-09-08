import { View, Text, ScrollView } from 'react-native'
import { SafeAreaView } from 'react-native'
import React, { useEffect, useState } from 'react'
import MultiStepForm from '../../components/MultiStepForm'
import FormField from '../../components/FormField'
import Countries from '../../constants/country'
import Regions from '../../constants/region'
import payload from '../../service/payload'
import ModalSelect, { SelectCountryItem } from '../../components/ModalSelect'

const FieldItem = ({
  title,
  value,
  handleChangeText,
  editable=true,
  keyboardType='ascii-capable'
}) => (
  <FormField
    title={title}
    value={value}
    handleChangeText={handleChangeText}
    keyboardType={keyboardType}
    containerStyles="mt-2"
    variant='flat'
    editable={editable}
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
    containerStyles="mt-2"
    variant="flat"
    inputStyles="text-slate-500"
    inputContainerStyles="border-slate-700 h-7"
    editable={false}
  />
)

const Onboarding = () => {

  const [errors, setErrors] = useState({});

  const [form, setForm] = useState({
    firstName: '',
    middleName: '',
    lastName: '',
    address1: '',
    address2: '',
    city: '',
    province: '',
    country: 'CA',
    postalCode: '',
    phoneNumber: '',
    dob: '',
    pob: '',
    nationality: '',
    identificationType: '',
    identificationValue: '',
  })

  useEffect(() => {
    async function validate() {
      try {
        await payload.OnboardingScheme.validate(form, {abortEarly: false})
        setErrors({})
      } catch (err) {
        let e = {}
        for ( let i of err.inner ) {
          e[i.path] =  e[i.path] ?? i.errors[0]
        }
        // setErrors(e)
      }
    }
    validate()
  }, [form])

  const submit = () => {
    setIsSubmitting(true)
    setTimeout(() => {
      setIsSubmitting(false)
      console.log(form)
      // router.replace("/verify_email")
    }, 4000);
  }

  const [isSubmitting, setIsSubmitting] = useState(false)

  const NameFieldTitle = () => (
    <View>
      <Text className="text-lg font-pbold text-center">Let's Get to Know You!</Text>
    </View>
  )
  const NameField = () => (
    <View>
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

  const AddressFieldTitle = () => (
    <View>
      <Text className="text-lg font-pbold text-center">Your Residential Address and Phone Number</Text>
    </View>
  )
  const AddressField = () => (
    <View>
      <Text className="font-pregular text-center mb-4">
        We require this information to continue setting up your Foree Remittance account
      </Text>

      <FieldItem title="Address Line 1" value={form.address1}
        handleChangeText={(e) => setForm({
          ...form,
          address1:e
        })}
      />
      <FieldItem title="Address Line 2" value={form.address2}
        handleChangeText={(e) => setForm({
          ...form,
          address2:e
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
      <FieldItem title="Country" value={Countries[form.country]?.name}
        handleChangeText={(e) => setForm({
          ...form,
          country:e
        })}
        editable={false}
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

  const PersonalDetailFieldTitle = () => (
    <View>
      <Text className="text-lg font-pbold text-center">Personal Details</Text>
    </View>
  )

  const PersonalDetailField = () => (
    <View>
      <Text className="font-pregular text-center mb-4">
        Almost done! Infomation below is requested by xxxxx xxxxx of xxxxxxx, our Foree Remittance payout parter, inorder to process your transfers under ...... regulatory guidelines
      </Text>
      <FieldItem title="Date of Birth(YYYY-MM-DD)" value={form.dob}
        handleChangeText={(e) => setForm({
          ...form,
          dob:e
        })}
        keyboardType="numbers-and-punctuation"
      />
      <ModalSelect
        title="Place of Birth"
        allowAdd={false}
        value={Countries[form.pob]?.name}
        variant='flat'
        listView={SelectCountryItem}
        list={Object.values(Countries)}
        onPress={(o) => {
          setForm({
            ...form,
            pob: o.isoCode
          })
        }}
        placeholder="select a country"
      />
      <ModalSelect
        title="Nationality"
        allowAdd={false}
        value={Countries[form.nationality]?.name}
        variant='flat'
        listView={SelectCountryItem}
        list={Object.values(Countries)}
        onPress={(o) => {
          setForm({
            ...form,
            nationality: o.isoCode
          })
        }}
        placeholder="select a country"
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

  const ReviewTitle = () => (
    <View>
      <Text className="text-lg font-pbold text-center">Review</Text>
    </View>
  )

  const Review = () => (
    <View>
      <Text className="font-pregular text-center mb-4">
        Please review your information.
      </Text>
      <ReviewItem title="First Name" value={form.firstName}/>
      <ReviewItem title="Middle Name" value={form.middleName}/>
      <ReviewItem title="Last Name" value={form.lastName}/>
      <ReviewItem title="Address Line 1" value={form.address1}/>
      <ReviewItem title="Address Line 2" value={form.address2}/>
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
      titleView: NameFieldTitle,
      formView: NameField,
      canGoNext: () => {
        return !errors.firstName && 
          !errors.middleName && 
          !errors.lastName
      }
    },
    {
      titleView: AddressFieldTitle,
      formView: AddressField,
      canGoNext: () => {
        return !errors.address1 && 
          !errors.address2 && 
          !errors.city &&
          !errors.province &&
          !errors.country &&
          !errors.postalCode &&
          !errors.phoneNumber
      }
    },
    {
      titleView: PersonalDetailFieldTitle,
      formView: PersonalDetailField,
      canGoNext: () => {
        return !errors.dob && 
          !errors.pob && 
          !errors.nationality &&
          !errors.identificationType &&
          !errors.identificationValue
      }
    },
    {
      titleView: ReviewTitle,
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
        containerStyle=""
        isSubmitting={isSubmitting}
      />
    </SafeAreaView>
  )
}

export default Onboarding