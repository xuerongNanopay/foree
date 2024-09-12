import { View, Text, ActivityIndicator } from 'react-native'
import React, { useEffect, useState } from 'react'

import { accountPayload } from '../../service'

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

  const submit = async () => {
    setIsSubmitting(true)
    try {
    } catch (err) {
      console.error("create contact", err)
    } finally {
      setIsSubmitting(false)
    }
  }

  const ContactInfoField = () => (
    <View>
      
    </View>
  )

  return (
    <View>
      <Text>CreateContact</Text>
    </View>
  )
}

export default CreateContact