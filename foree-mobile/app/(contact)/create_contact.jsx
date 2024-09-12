import { View, Text } from 'react-native'
import React, { useEffect, useState } from 'react'

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
    relationshipToContact: '',
    identificationType: '',
    transferMethod: '',
    bankName: '',
    accountNoOrIBAN: ''
  })

  return (
    <View>
    <Text>CreateContact</Text>
    </View>
  )
}

export default CreateContact