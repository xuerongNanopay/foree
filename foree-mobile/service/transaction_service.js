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

  async getRewards({signal}={signal}) {
    return await this.#httpClient.get("/transaction_reward", {signal})
  }

  async quote(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/quote", req, {signal})
  }

  async confirmQuote(quoteId, {signal}={signal}) {
    return await this.#httpFormClient.post("/create_transaction", {quoteId}, {signal})
  }

  async getTransactions({status="", offset=0, limit=10}={status, offset, limit}, {signal}={signal}) {
    const searchParams = new URLSearchParams({status, offset, limit})
    console.log(searchParams.toString())
    return await this.#httpClient.get(`/transactions?${searchParams.toString()}`)
  }

  async getTransaction(transactionId, {signal}={signal}) {
    return await this.#httpClient.get(`/transactions/${transactionId}`)
  }
}

export default TransactionService