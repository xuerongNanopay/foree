class GeneralService {
  #httpClient

  constructor(httpClient) {
    this.#httpClient = httpClient
  }

  async cusomterSupport({signal}={signal}) {
    return await this.#httpClient.get("/customer_support", {signal})
  }
}

export default GeneralService