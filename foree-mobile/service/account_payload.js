import { object, string } from 'yup'
import fieldScheme from '../constants/validation_sheme'

const CreateContactScheme = object({
  firstName: fieldScheme.NameScheme(),
  middleName: fieldScheme.NameScheme({required: false}),
  lastName: fieldScheme.NameScheme(),
  address1: fieldScheme.String(),
  address2: fieldScheme.String({required: false}),
  city: fieldScheme.String(),
  province: fieldScheme.ProvinceISOScheme(),
  country: 'PK',
  postalCode: fieldScheme.PostalCodeScheme({countryCode:"CA", required: false}),
  phoneNumber: fieldScheme.PhoneNumber({countryCode:"CA", required: false}),
  relationshipToContact: fieldScheme.String(),
  transferMethod: fieldScheme.String(),
  bankName: string().when(["transferMethod"], {
    is: (transferMethod) => transferMethod !== "CASH_PICKUP",
    then: fieldScheme.String()
  }),
  accountNoOrIBAN: string().when(["transferMethod"], {
    is: (transferMethod) => transferMethod !== "CASH_PICKUP",
    then: fieldScheme.String()
  })
})

export default {
  CreateContactScheme
}