class TransactionService {
  #httpClient
  #httpFormClient
  constructor(formClient, httpClient) {
    this.#httpFormClient = formClient
    this.#httpClient = httpClient
  }

  async getCADToPRKRate(req, {signal}={signal}) {
    return await this.#httpClient.post("/rate", {
      srcCurrency: "CAD",
      destCurrency: "PKR"
    }, {signal})
  }
}

export default TransactionService