class TransactionService {
  #httpClient
  #httpFormClient
  constructor(formClient, httpClient) {
    this.#httpFormClient = formClient
    this.#httpClient = httpClient
  }

  async getCADToPRKRate({signal}={signal}) {
    return await this.#httpClient.post("/rate", {
      srcCurrency: "CAD",
      destCurrency: "PKR"
    }, {signal})
  }

  async getDailyLimit({signal}={signal}) {
    return await this.#httpClient.get("/transaction_limit", {signal})
  }

  async quote(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/quote", req, {signal})
  }

  async confirmQuote(quoteId, {signal}={signal}) {
    return await this.#httpFormClient.post("/create_transaction", {quoteId}, {signal})
  }
}

export default TransactionService