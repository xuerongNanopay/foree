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

export const formatContactMethod = ({transferMethod, bankName, accountNoOrIBAN}, max=20) => {
  if ( transferMethod === ContactTransferCashPickup ) return "Cash Pickup"
  return `${!!bankName ? bankName.slice(0, 14) + (bankName.length > 14 ? "..." : "") : ""}(${!!accountNoOrIBAN ? accountNoOrIBAN.slice(0, 14) + (accountNoOrIBAN.length > 14 ? "..." : "") : ""})`
}