import { View, Text, Alert } from 'react-native'
import { SafeAreaView } from 'react-native'
import React, { useCallback, useEffect, useState } from 'react'

import MultiStepForm from '../../components/MultiStepForm'
import FormField from '../../components/FormField'
import { accountPayload, accountService } from '../../service'
import Countries from '../../constants/country'
import Regions from '../../constants/region'
import ModalSelect from '../../components/ModalSelect'
import { ContactTransferBank, ContactTransferMethods, PersonalRelationships } from '../../constants/contacts'
import string_util from '../../util/string_util'
import { router } from 'expo-router'

const ContactCreate = () => {
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
    country: 'PK',
    postalCode: '',
    phoneNumber: '',
    relationshipToContact: '',
    transferMethod: '',
    bankName: '',
    accountNoOrIBAN: ''
  })

  useEffect(() => {
    async function validate() {
      try {
        await accountPayload.CreateContactScheme.validate(form, {abortEarly: false})
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

  const submitNewContact = async () => {
    setIsSubmitting(true)
    try {
      const resp = await accountService.createContact(string_util.trimStringInObject(form))
      if ( resp.status / 100 !== 2 ) {
        console.warn("create contact", resp.status, resp.data)
        return
      }
    } catch (err) {
      console.error("create contact", err, err.response.data)
    } finally {
      if ( router.canGoBack() ) {
        router.back()
      } else {
        router.replace('/contact_tab')
      }
      setIsSubmitting(false)
    }
  }

  const submit = async () => {
    setIsSubmitting(true)
    try {
      if ( form.transferMethod === "CASH_PICKUP" ) {
        submitNewContact()
      } else {
        const resp = await accountService.verifyContact(string_util.trimStringInObject(form))
        if ( resp.status / 100 === 2 && 
          (resp?.data?.data?.accountStatus === "Active" || resp?.data?.data?.accountStatus === "Dormant or Inoperative")
        ) {
          submitNewContact()
        } else if ( resp.status / 100 === 2 && !!resp?.data?.data?.accountStatus ) {
          Alert.alert("Error", `We can't send fund to a ${resp?.data?.data?.accountStatus} account.`, [
            {text: 'OK', onPress: () => {}},
          ])
        }  else {
          Alert.alert("Warming", "We can't verify your contact bank details through NBP. if you are sure it is correct, please click 'continue'", [
            {text: 'Cancel', onPress: () => {}},
            {text: 'Continue', onPress: () => {submitNewContact()}},
          ])
        }
      }
    } catch (err) {
      console.error("create contact", err)
    } finally {
      setIsSubmitting(false)
    }
  }

  const ContactNameTitle = useCallback(() => (
    <View>
      <Text className="text-lg font-pbold text-center">Contact Name</Text>
    </View>
  ), [])

  const ContactName = useCallback(() => (
    <View>
      <Text className="font-pregular text-center mb-4">
        Please provide contact name.
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

  const ContactAddressTitle = useCallback(() => (
    <View>
      <Text className="text-lg font-pbold text-center">Contact Address</Text>
    </View>
  ), [])

  const ContactAddress = useCallback(() => (
    <View>
      <Text className="font-pregular text-center mb-4">
        Please provide contact address details.
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
        keyExtractor="code"
        showExtractor="name"
        valueExtractor="isoCode"
        listView={SelectProvinceItem}
        list={Object.values(Regions[form.country])}
        onPress={(o) => {
          setForm({
            ...form,
            province: o
          })
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
        keyboardType="name-phone-pad"
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

  const ContactBankInfoTitle = useCallback(() => (
    <View>
      <Text className="text-lg font-pbold text-center">Contact Bank Info</Text>
    </View>
  ), [])

  const ContactBankInfo = useCallback(() => (
    <View>
      <Text className="font-pregular text-center mb-4">
        Please provide contact banking information.
      </Text>
  
      <ModalSelect
        title="Relationship to Contact"
        errorMessage={errors['relationshipToContact']}
        modalTitle="select relationship"
        containerStyles="mt-2"
        value={form.relationshipToContact}
        variant='flat'
        searchKey="name"
        keyExtractor="name"
        showExtractor="name"
        valueExtractor="name"
        listView={SelectPersonalRelationshipItem}
        list={PersonalRelationships}
        onPress={(o) => {
          setForm((form) => ({
            ...form,
            relationshipToContact: o
          }))
        }}
        placeholder="choose relationship..."
      />
      <ModalSelect
        title="Transfer Method"
        errorMessage={errors['transferMethod']}
        modalTitle="select transfer method"
        containerStyles="mt-2"
        value={form.transferMethod}
        variant='flat'
        keyExtractor="name"
        showExtractor="name"
        valueExtractor="value"
        listView={SelectTransferMethodItem}
        list={ContactTransferMethods}
        onPress={(o) => {
          setForm((form) => ({
            ...form,
            transferMethod: o,
            bankName: "",
            accountNoOrIBAN: ""
          }))
        }}
        placeholder="choose transfer method"
      />
      {
        !!form.transferMethod && form.transferMethod !== "CASH_PICKUP" ? 
        <>
          <ModalSelect
            key={form.transferMethod}
            title="Bank Name"
            errorMessage={errors['bankName']}
            modalTitle="select bank"
            containerStyles="mt-2"
            allowSearch={true}
            value={form.bankName}
            variant='flat'
            searchKey="bankName"
            keyExtractor="bankAbbr"
            showExtractor="bankName"
            valueExtractor="bankAbbr"
            listView={SelectBankItem}
            list={ContactTransferBank[form.transferMethod]}
            onPress={(o) => {
              setForm({
                ...form,
                bankName: o,
                accountNoOrIBAN: ""
              })
            }}
            placeholder="choose transfer method"
          />
          <FieldItem title="Account No or IBAN" value={form.accountNoOrIBAN}
            errorMessage={errors['accountNoOrIBAN']}
            handleChangeText={(e) => setForm((form) => ({
              ...form,
              accountNoOrIBAN:e
            }))}
          />
        </> : null
      }
    </View>
  ), [
    form.relationshipToContact, errors['relationshipToContact'],
    form.transferMethod, errors['transferMethod'],
    form.bankName, errors['bankName'],
    form.accountNoOrIBAN, errors['accountNoOrIBAN']
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
      <ReviewItem title="Relationship to Contact" value={PersonalRelationships.find(r => r.name === form.PersonalRelationships)?.["name"]}/>
      <ReviewItem title="Transfer Method" value={ContactTransferMethods.find(r => r.value === form.transferMethod)?.["name"]}/>
      {
        !!form.transferMethod && form.transferMethod !== "CASH_PICKUP" ? 
        (
          <>
            <ReviewItem title="Bank Name" value={ContactTransferBank[form.transferMethod]?.find(r => r.bankAbbr === form.bankName)?.["bankName"]}/>
            <ReviewItem title="Account No or IBAN" value={form.accountNoOrIBAN}/>
          </>
        ) : null
      }
    </View>
  )

  const CreateContactFlow = [
    {
      titleView: ContactNameTitle,
      formView: ContactName,
      canGoNext: () => {
        return !errors.firstName && 
          !errors.middleName && 
          !errors.lastName
      }
    },
    {
      titleView: ContactAddressTitle,
      formView: ContactAddress,
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
      titleView: ContactBankInfoTitle,
      formView: ContactBankInfo,
      canGoNext: () => {
        return !errors.relationshipToContact && 
          !errors.transferMethod && 
          !errors.bankName &&
          !errors.accountNoOrIBAN
      }
    },
    {
      titleView: ReviewTitle,
      formView: Review,
      canGoNext: () => {
        return true
      }
    },
  ]

  return (
    <SafeAreaView className="bg-slate-100">
      <MultiStepForm
        steps={() => CreateContactFlow}
        onSumbit={submit}
        containerStyle=""
        submitDisabled={isSubmitting}
      />
    </SafeAreaView>
  )
}

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

const SelectProvinceItem = (province) => (
  <Text className="font-pregular py-3 text-xl">
    {province["name"]}
  </Text>
)

const SelectPersonalRelationshipItem = (relationship) => (
  <Text className="font-pregular py-3 text-xl">
    {relationship["name"]}
  </Text>
)

const SelectTransferMethodItem = (transferMethod) => (
  <Text className="font-pregular py-3 text-xl">
    {transferMethod["name"]}
  </Text>
)

const SelectBankItem = (bank) => (
  <Text className="font-pregular py-3 text-xl">
    {bank["bankName"]}
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

export default ContactCreate