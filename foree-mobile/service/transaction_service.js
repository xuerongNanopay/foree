import axios from 'axios'

class TransactionService {

  async getCADToPRKRate(req, {signal}={signal}) {
    return await axios.post("/rate", {
      srcCurrency: "CAD",
      destCurrency: "PKR"
    }, {signal})
  }
}

export default TransactionService