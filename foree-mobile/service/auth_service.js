class AuthService {
  #localLogout
  #httpFormClient
  constructor(localLogout, httpFormClient) {
    this.#localLogout = localLogout
    this.#httpFormClient = httpFormClient
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

  async forgetPassword(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/forget_password", req, {signal})
  }

  async forgetPasswordVerify(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/forget_password_verify", req, {signal})
  }

  async forgetPasswordUpdate(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/forget_password_update", req, {signal})
  }

  async onboard(req, {signal}={signal}) {
    return await this.#httpFormClient.post("/onboard", req, {signal})
  }

  async resendCode({signal}={signal}) {
    return await this.#httpFormClient.get("/resend_code", {signal})
  }

  async getUser({signal}={signal}) {
    return await this.#httpFormClient.get("/user", {signal})
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