import axios from 'axios'

class AccountService {
  #httpClient
  #httpFormClient
  constructor(formClient, httpClient) {
    this.#httpFormClient = formClient
    this.#httpClient = httpClient
  }

  async verifyContact(req, {signal}={signal}) {
    return await this.#httpClient.post("/verify_contact_account", req, {signal})
  }

  async createContact(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/create_contact_account", req, {signal})
  }
}

export default AccountService