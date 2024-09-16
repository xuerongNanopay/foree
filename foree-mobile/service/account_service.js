class AccountService {
  #httpClient
  #httpFormClient
  #allContactCache
  #allContactCacheExpiry

  constructor(formClient, httpClient) {
    this.#httpFormClient = formClient
    this.#httpClient = httpClient
    //TODO: configue the property
    this.#allContactCacheExpiry = 30000
  }

  async verifyContact(req, {signal}={signal}) {
    return await this.#httpClient.post("/verify_contact_account", req, {signal})
  }

  async createContact(req, {signal}={signal}) {
    this.#allContactCache = null
    return await this.#httpFormClient.post("/create_contact_account", req, {signal})
  }

  async deleteContact(contactId, {signal}={signal}) {
    this.#allContactCache = null
    return await this.#httpFormClient.post("/delete_contact_account", {contactId}, {signal})
  }

  async getInteracAccounts({signal}={signal}) {
    return await this.#httpClient.get("/interac_accounts", {signal})
  }

  async getAllContactAccounts({signal}={signal}) {
    if ( this.#allContactCache != null &&  this.#allContactCache.expiryAt.getTime() > new Date().getTime() )  {
      return this.#allContactCache.contacts
    }

    try {
      resp = await this.#httpClient.get("/contact_accounts", {signal})
      if ( resp.status / 100 == 2 && !!resp?.data?.data ) {
        this.#allContactCache = {
          contacts: resp,
          expiryAt: new Date((new Date()).getTime() + this.#allContactCacheExpiry)
        }
      }
      return resp
    } catch (e) {
      this.#allContactCache = null
      throw e
    }
  }

  async getContactAccount(contactId, {signal}={signal}) {
    return await this.#httpClient.get(`/contact_accounts/${contactId}`, {signal})
  }
}

export default AccountService