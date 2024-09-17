import { number, object, string } from "yup"
const QuoteTransactoinScheme = object({
  cinAccId: number().required("required"),
  coutAccId: number().integer().required("required").min(1, "required"),
  srcAmount: number().required("required").min(20, "Minimum limit $ 20.00 CAD").max(1000, "Maximum limit $1000.00 CAD"),
  transactionPurpose: string().required("required")
})

export default {
  QuoteTransactoinScheme
}