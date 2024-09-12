import { View, Text } from 'react-native'
import React, { useEffect, useState } from 'react'

const transferMethods = [
  {
    id: "cash",
    name: "Cash Pickup",
    value: "CASH"
  },
  {
    id: "bankAccount",
    name: "Bank Account",
    value: "ACCOUNT_TRANSFERS"
  },
  {
    id: "mobileWallet",
    name: "Mobile Wallet",
    value: "THIRD_PARTY_PAYMENTS"
  },
  {
    id: "Roshan Digital Account",
    name: "Roshan Digit Account",
    value: "THIRD_PARTY_PAYMENTS"
  },
]
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
    country: 'CA',
    postalCode: '',
    phoneNumber: '',
    dob: '',
    pob: '',
    nationality: '',
    identificationType: '',
    identificationValue: '',
  })

  return (
    <View>
    <Text>CreateContact</Text>
    </View>
  )
}

export default CreateContact