import { ContactTransferCashPickup } from "../constants/contacts"

export const formatName = ({firstName, middleName, lastName}, max=30) => {
  let fullName = !!middleName ? `${firstName} ${middleName} ${lastName}` : `${firstName} ${lastName}`
  if ( fullName.length < max ) return fullName
  fullName = `${firstName} ${lastName}`
  if ( fullName.length < max ) return fullName
  fullName = `${firstName}`
  if ( fullName.length < max ) return fullName
  fullName = `${lastName}`
  if ( fullName.length < max ) return fullName
  fullName = `${firstName}`
  fullName = fullName.slice(0, max-3) + "..."
  return fullName
}

export const formatContactMethod = ({transferMethod, bankName, accountNoOrIBAN}, max=14) => {
  if ( transferMethod === ContactTransferCashPickup ) return "Cash Pickup"
  return `${!!bankName ? bankName.slice(0, max) + (bankName.length > max ? "..." : "") : ""}(${!!accountNoOrIBAN ? accountNoOrIBAN.slice(0, max) + (accountNoOrIBAN.length > max ? "..." : "") : ""})`
}

export const currencyFormatter = (amount, currency, isNegative=false) => {
  return `${isNegative?"-":""}$${new Intl.NumberFormat("en", {minimumFractionDigits: 2}).format(amount)}${!!currency ? ` ${currency}` : ''}`
}