import { View, Text } from 'react-native'
import { SafeAreaView } from 'react-native'
import React, { useEffect, useState } from 'react'

import MultiStepForm from '../../components/MultiStepForm'
import FormField from '../../components/FormField'
import { accountPayload } from '../../service'
import Countries from '../../constants/country'
import Regions from '../../constants/region'
import ModalSelect from '../../components/ModalSelect'
import { ContactTransferBank, ContactTransferMethods, PersonalRelationships } from '../../constants/contacts'

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

const CreateContact = () => {
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
        console.log(e)
      }
    }
    validate()
  }, [form])

  const submit = async () => {
    setIsSubmitting(true)
    try {
    } catch (err) {
      console.error("create contact", err)
    } finally {
      setIsSubmitting(false)
    }
  }

  const ContactNameTitle = () => (
    <View>
      <Text className="text-lg font-pbold text-center">Contact Name</Text>
    </View>
  )

  const ContactName = () => (
    <View>
      <Text className="font-pregular text-center mb-4">
        Please provide contact name.
      </Text>
  
      <FieldItem title="First Name" value={form.firstName}
        errorMessage={errors['firstName']}
        handleChangeText={(e) => setForm({
          ...form,
          firstName:e
        })}
      />
      <FieldItem title="Middle Name" value={form.middleName}
        errorMessage={errors['middleName']}
        handleChangeText={(e) => setForm({
          ...form,
          middleName:e
        })}
      />
      <FieldItem title="Last Name" value={form.lastName}
        errorMessage={errors['lastName']}
        handleChangeText={(e) => setForm({
          ...form,
          lastName:e
        })}
      />
    </View>
  )

  const ContactAddressTitle = () => (
    <View>
      <Text className="text-lg font-pbold text-center">Contact Address</Text>
    </View>
  )

  const ContactAddress = () => (
    <View>
      <Text className="font-pregular text-center mb-4">
        Please provide contact address details.
      </Text>

      <FieldItem title="Address Line 1" value={form.address1}
        errorMessage={errors['address1']}
        handleChangeText={(e) => setForm({
          ...form,
          address1:e
        })}
      />
      <FieldItem title="Address Line 2" value={form.address2}
        errorMessage={errors['address2']}
        handleChangeText={(e) => setForm({
          ...form,
          address2:e
        })}
      />
      <FieldItem title="City" value={form.city}
        errorMessage={errors['city']}
        handleChangeText={(e) => setForm({
          ...form,
          city:e
        })}
      />
      <ModalSelect
        title="Province"
        modalTitle="select a province"
        errorMessage={errors['province']}
        containerStyles="mt-2"
        allowSearch={false}
        allowAdd={false}
        value={Regions[form.country]?.[form.province]?.name}
        variant='flat'
        searchKey="name"
        keyExtractor="code"
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
        handleChangeText={(e) => setForm({
          ...form,
          country:e
        })}
        editable={false}
      />
      <FieldItem title="Postal Code" value={form.postalCode}
        errorMessage={errors['postalCode']}
        handleChangeText={(e) => setForm({
          ...form,
          postalCode:e
        })}
      />
      <FieldItem title="Phone Number" value={form.phoneNumber}
        errorMessage={errors['phoneNumber']}
        handleChangeText={(e) => setForm({
          ...form,
          phoneNumber:e
        })}
        keyboardType="name-phone-pad"
      />
    </View>
  )

  const ContactBankInfoTitle = () => (
    <View>
      <Text className="text-lg font-pbold text-center">Contact Bank Info</Text>
    </View>
  )

  const ContactBankInfo = () => (
    <View>
      <Text className="font-pregular text-center mb-4">
        Please provide contact banking information.
      </Text>
  
      <ModalSelect
        title="Relationship to Contact"
        errorMessage={errors['relationshipToContact']}
        modalTitle="select relationship"
        containerStyles="mt-2"
        allowAdd={false}
        allowSearch={false}
        value={PersonalRelationships[form.relationshipToContact]?.name}
        variant='flat'
        searchKey="name"
        keyExtractor="name"
        valueExtractor="name"
        listView={SelectPersonalRelationshipItem}
        list={Object.values(PersonalRelationships)}
        onPress={(o) => {
          console.log(o)
          setForm({
            ...form,
            relationshipToContact: o
          })
        }}
        placeholder="choose relationship..."
      />
      <ModalSelect
        title="Transfer Method"
        errorMessage={errors['transferMethod']}
        modalTitle="select transfer method"
        containerStyles="mt-2"
        allowAdd={false}
        allowSearch={false}
        value={ContactTransferMethods[form.transferMethod]?.name}
        variant='flat'
        keyExtractor="name"
        valueExtractor="value"
        listView={SelectTransferMethodItem}
        list={Object.values(ContactTransferMethods)}
        onPress={(o) => {
          console.log(o)
          setForm({
            ...form,
            transferMethod: o,
            bankName: "",
            accountNoOrIBAN: ""
          })
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
            allowAdd={false}
            allowSearch={true}
            value={form.bankName}
            variant='flat'
            searchKey="bankName"
            keyExtractor="bankName"
            valueExtractor="bankName"
            listView={SelectBankItem}
            list={ContactTransferBank[form.transferMethod]}
            onPress={(o) => {
              console.log(o)
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
            handleChangeText={(e) => setForm({
              ...form,
              accountNoOrIBAN:e
            })}
          />
        </> : null
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

export default CreateContact