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
  country: fieldScheme.String(),
  postalCode: fieldScheme.PostalCodeScheme({countryCode:"PK", required: false}),
  // phoneNumber: fieldScheme.PhoneNumber({countryCode:"PK", required: false}),
  phoneNumber: fieldScheme.NumberOnlyScheme({required:false}),
  relationshipToContact: fieldScheme.String(),
  transferMethod: fieldScheme.String(),
  bankName: string().when(["transferMethod"], ([transferMethod]) => {
    switch (transferMethod) {
      case "":
        return fieldScheme.String({required: false})
      case "CASH_PICKUP":
        return fieldScheme.String({required: false})
      default:
        return fieldScheme.String()
    }
  }),
  accountNoOrIBAN: string().when(["transferMethod"], ([transferMethod]) => {
    switch (transferMethod) {
      case "":
        return fieldScheme.String({required: false})
      case "CASH_PICKUP":
        return fieldScheme.String({required: false})
      default:
        return fieldScheme.String()
    }
  }),
})

export default {
  CreateContactScheme
}