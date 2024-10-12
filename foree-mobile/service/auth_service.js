class AuthService {
  #localLogout
  #httpClient
  #httpFormClient
  constructor(localLogout, httpFormClient, httpClient) {
    this.#localLogout = localLogout
    this.#httpFormClient = httpFormClient
    this.#httpClient = httpClient
  }

  async login(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/login", req, {signal})
  }

  async signUp(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/sign_up", req, {signal})
  }

  async verifyEmail(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/verify_email", req, {signal})
  }

  async resendCode({signal}={signal}) {
    return await this.#httpFormClient.get("/resend_code", {signal})
  }

  async forgetPasswd(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/forget_passwd", req, {signal})
  }

  async forgetPasswdVerify(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/forget_passwd_verify", req, {signal})
  }

  async forgetPasswdUpdate(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/forget_passwd_update", req, {signal})
  }

  async onboard(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/onboard", req, {signal})
  }

  async resendCode({signal}={signal}) {
    return await this.#httpFormClient.get("/resend_code", {signal})
  }

  async getUser({signal}={signal}) {
    return await this.#httpClient.get("/user", {signal})
  }

  async getUserDetail({signal}={signal}) {
    return await this.#httpClient.get("/user_detail", {signal})
  }

  async getUserSetting({signal}={signal}) {
    return await this.#httpClient.get("/user_setting", {signal})
  }

  async getUserExtra({signal}={signal}) {
    return await this.#httpClient.get("/user_extra", {signal})
  }

  async updatePasswd(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/update_passwd", req, {signal})
  }

  async updateAddress(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/update_address", req, {signal})
  }

  async updatePhone(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/update_phone", req, {signal})
  }

  async updateUserSetting(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/update_user_setting", req, {signal})
  }

  async logout({signal}={signal}) {
    try {
      await this.#httpFormClient.get("/logout", {signal})
      await this.#localLogout()
    } catch(e) {
      console.error(e)
      //TODO: send error
    }
  }
}

export default AuthService