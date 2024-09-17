import { View, Text, ScrollView } from 'react-native'
import { SafeAreaView } from 'react-native'
import React, { useCallback, useEffect, useState } from 'react'
import { router } from 'expo-router'

import MultiStepForm from '../../components/MultiStepForm'
import FormField from '../../components/FormField'
import Countries from '../../constants/country'
import Regions from '../../constants/region'
import { authPayload, authService } from '../../service'
import ModalSelect, { SelectCountryItem } from '../../components/ModalSelect'
import { useGlobalContext } from '../../context/GlobalProvider'
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
    variant='flat'
    editable={editable}
    errorMessage={errorMessage}
  />
)

const IDTypes = [
  {
    id: "PASSPORT",
    name: "Passport"
  },
  {
    id: "DRIVER_LICENSE",
    name: "Driver License"
  },
  {
    id: "PROVINCIAL_ID",
    name: "Provincial Id"
  },
  {
    id: "NATIONAL_ID",
    name: "Nation Id"
  }
]


const SelectIDTypesItem = (idType) => (
  <Text className="font-pregular py-3 text-xl">
    {idType["name"]}
  </Text>
)

const SelectProvinceItem = (province) => (
  <Text className="font-pregular py-3 text-xl">
    {province["name"]}
  </Text>
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
  const { setUser } = useGlobalContext()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [errors, setErrors] = useState({})
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
        await authPayload.OnboardingScheme.validate(form, {abortEarly: false})
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

  const submit = async () => {
    setIsSubmitting(true)
    try {
      const resp = await authService.onboard(string_util.trimStringInObject(form))
      if ( resp.status / 100 !== 2 ) {
        console.warn("onboarding", resp.status, resp.data)
        return
      }
      se = resp.data.data
      setUser(se)
      router.replace("/home_tab")
    } catch (err) {
      console.error("onboarding", err)
    } finally {
      setIsSubmitting(false)
    }
  }


  const NameFieldTitle = useCallback(() => (
    <View>
      <Text className="text-lg font-pbold text-center">Let's Get to Know You!</Text>
    </View>
  ), [])

  const NameField = useCallback(() => (
    <View>
      <Text className="font-pregular text-center mb-4">
        Please enter your full legal name so we can begin setting up your account
      </Text>

      <FieldItem title="First Name" value={form.firstName}
        errorMessage={errors['firstName']}
        handleChangeText={(e) => setForm((form) => ({
          ...form,
          firstName:e
        }))}
      />
      <FieldItem title="Middle Name" value={form.middleName}
        errorMessage={errors['middleName']}
        handleChangeText={(e) => setForm((form) => ({
          ...form,
          middleName:e
        }))}
      />
      <FieldItem title="Last Name" value={form.lastName}
        errorMessage={errors['lastName']}
        handleChangeText={(e) => setForm((form) => ({
          ...form,
          lastName:e
        }))}
      />
    </View>
  ), [
    form.firstName, errors['firstName'],
    form.middleName, errors['middleName'],
    form.lastName, errors['lastName']
  ])

  const AddressFieldTitle = useCallback(() => (
    <View>
      <Text className="text-lg font-pbold text-center">Your Residential Address and Phone Number</Text>
    </View>
  ), [])

  const AddressField = useCallback(() => (
    <View>
      <Text className="font-pregular text-center mb-4">
        We require this information to continue setting up your Foree Remittance account
      </Text>

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
        modalTitle="select a province"
        errorMessage={errors['province']}
        containerStyles="mt-2"
        value={form.province}
        variant='flat'
        searchKey="name"
        keyExtractor="name"
        showExtractor="name"
        valueExtractor="isoCode"
        listView={SelectProvinceItem}
        list={Object.values(Regions[form.country])}
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
      <FieldItem title="Phone Number" value={form.phoneNumber}
        errorMessage={errors['phoneNumber']}
        handleChangeText={(e) => setForm((form) => ({
          ...form,
          phoneNumber:e
        }))}
        keyboardType="phone-pad"
      />
    </View>
  ), [
    form.address1, errors['address1'],
    form.address2, errors['address2'],
    form.city, errors['city'],
    form.province, errors['province'],
    form.country, errors['country'],
    form.postalCode, errors['postalCode'],
    form.phoneNumber, errors['phoneNumber'],
  ])

  const PersonalDetailFieldTitle = useCallback(() => (
    <View>
      <Text className="text-lg font-pbold text-center">Personal Details</Text>
    </View>
  ),[])

  const PersonalDetailField = useCallback(() => (
    <View>
      <Text className="font-pregular text-center mb-4">
        Almost done! Infomation below is requested by xxxxx xxxxx of xxxxxxx, our Foree Remittance payout parter, inorder to process your transfers under ...... regulatory guidelines
      </Text>
      <FieldItem title="Date of Birth(YYYY-MM-DD)" value={form.dob}
        errorMessage={errors['dob']}
        handleChangeText={(e) => setForm((form) => ({
          ...form,
          dob:e
        }))}
        keyboardType="numbers-and-punctuation"
      />
      <ModalSelect
        title="Place of Birth"
        modalTitle="select a country"
        errorMessage={errors['pob']}
        containerStyles="mt-2"
        value={() => {
          return Countries[form.pob] ? `${Countries[form.pob]?.unicodeIcon} ${Countries[form.pob]?.name}` : ""
        }}
        variant='flat'
        searchKey="name"
        keyExtractor="name"
        valueExtractor="isoCode"
        listView={SelectCountryItem}
        list={Object.values(Countries)}
        onPress={(o) => {
          setForm((form) => ({
            ...form,
            pob: o
          }))
        }}
        placeholder="select a country"
      />
      <ModalSelect
        title="Nationality"
        modalTitle="select a country"
        containerStyles="mt-2"
        errorMessage={errors['nationality']}
        value={() => {
          return Countries[form.nationality] ? `${Countries[form.nationality]?.unicodeIcon} ${Countries[form.nationality]?.name}` : ""
        }}
        variant='flat'
        searchKey="name"
        keyExtractor="name"
        valueExtractor="isoCode"
        listView={SelectCountryItem}
        list={Object.values(Countries)}
        onPress={(o) => {
          setForm((form) => ({
            ...form,
            nationality: o
          }))
        }}
        placeholder="select a country"
      />
      <ModalSelect
        title="Identification Document Type"
        errorMessage={errors['identificationType']}
        modalTitle="select identification type"
        containerStyles="mt-2"
        value={form.identificationType}
        variant='flat'
        searchKey="id"
        keyExtractor="id"
        showExtractor="name"
        valueExtractor="id"
        listView={SelectIDTypesItem}
        list={Object.values(IDTypes)}
        onPress={(o) => {
          setForm((form) => ({
            ...form,
            identificationType: o
          }))
        }}
        placeholder="select ID type"
      />
      <FieldItem title="Identification Number" value={form.identificationValue}
        errorMessage={errors['identificationValue']}
        handleChangeText={(e) => setForm((form) => ({
          ...form,
          identificationValue:e
        }))}
      />
    </View>
  ), [
    form.dob, errors['dob'],
    form.pob, errors['pob'],
    form.nationality, errors['nationality'],
    form.identificationType, errors['identificationType'],
    form.identificationValue, errors['identificationValue'],
  ])

  const ReviewTitle = useCallback(() => (
    <View>
      <Text className="text-lg font-pbold text-center">Review</Text>
    </View>
  ), [])

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
      <ReviewItem title="Province" value={Regions[form.country]?.[form.province]?.name}/>
      <ReviewItem title="Country" value={Countries[form.country]?`${Countries[form.country]?.unicodeIcon} ${Countries[form.country]?.name}`: ""}/>
      <ReviewItem title="Postal Code" value={form.postalCode}/>
      <ReviewItem title="Phone Number" value={form.phoneNumber}/>
      <ReviewItem title="Date of Birth" value={form.dob}/>
      <ReviewItem title="Place of Birth" value={Countries[form.pob]?`${Countries[form.pob]?.unicodeIcon} ${Countries[form.pob]?.name}`: ""}/>
      <ReviewItem title="Nationality" value={Countries[form.nationality]?`${Countries[form.nationality]?.unicodeIcon} ${Countries[form.nationality]?.name}`: ""}/>
      <ReviewItem title="Identification Document Type" value={IDTypes.find(i => i.id === form.identificationType)?.name}/>
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
        submitDisabled={isSubmitting}
      />
    </SafeAreaView>
  )
}

export default Onboarding