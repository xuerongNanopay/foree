import axios from 'axios'

import AuthService from "./auth_service"
import authPayload from "./auth_payload"
import { Alert } from 'react-native'
import { router } from 'expo-router'
import AsyncStorage from '@react-native-async-storage/async-storage'
import TransactionService from './transaction_service'

const SessionIdKey = "SESSION_ID"
// Config axios
axios.defaults.baseURL = 'http://192.168.2.30:8080/app/v1'

const httpClient = axios.create({
  baseURL: 'http://192.168.2.30:8080/app/v1',
  // timeout: 1000,
  // headers: {'X-Custom-Header': 'foobar'}
})

httpClient.interceptors.request.use(
  async (config) => {
    try {
      const session = await AsyncStorage.getItem(SessionIdKey)
      if ( !!session ) {
        console.debug("Request Session: ", session)
        config.headers[SessionIdKey] = session
      }
    } catch (e) {
      console.error("get session error", e)
      //TODO: send error
    } finally {
      return config
    }
  }
)

httpClient.interceptors.response.use(
  async (response) => {
    try {
      const body = response.data
      if ( !!body?.data?.sessionId ) {
        await AsyncStorage.setItem(SessionIdKey, body.data.sessionId)
      }
    } catch (e) {
      console.error("update session error", e)
      //TODO: send error
    } finally {
      return Promise.resolve(response)
    }
  },
  (error) => {
    if (!error.response) {
      //TODO: send error
      return Promise.reject(error)
    }

    resp = error.response
    // console.debug("url", error.request.responseURL)
    if (resp.status == 401 ) {
      router.replace("/login")
      return Promise.reject(error)
    } else if ( resp.status = 428 && resp?.data?.require === "TO_MAIN" ) {
      router.replace("/home")
    } else if ( resp.status = 428 && resp?.data?.require === "LOGIN" ) {
      router.replace("/login")
      return Promise.reject(error)
    } else if ( resp.status = 428 && resp?.data?.require === "VERIFY_EMAIL" ) {
      router.replace("/verify_email")
      return Promise.reject(error)
    } else if ( resp.status = 428 && resp?.data?.require === "CREATE_USER" ) {
      router.replace("/onboarding")
      return Promise.reject(error)
    }
  
    //TODO: send error log, error.message, error, error.request, error.response, user
    return Promise.reject(error) 
  }
)

const httpFormClient = axios.create({
  baseURL: 'http://192.168.2.30:8080/app/v1',
  // timeout: 1000,
  // headers: {'X-Custom-Header': 'foobar'}
})

httpFormClient.interceptors.request.use(
  async (config) => {
    try {
      const session = await AsyncStorage.getItem(SessionIdKey)
      if ( !!session ) {
        console.debug("Request Session: ", session)
        config.headers[SessionIdKey] = session
      }
    } catch (e) {
      console.error("get session error", e)
      //TODO: send error
    } finally {
      return config
    }
  }
)

httpFormClient.interceptors.response.use(
  async (response) => {
    try {
      const body = response.data
      if ( !!body?.data?.sessionId ) {
        await AsyncStorage.setItem(SessionIdKey, body.data.sessionId)
      }
    } catch (e) {
      console.error("update session error", e)
      //TODO: send error
    } finally {
      return Promise.resolve(response)
    }
  },
  (error) => {
    if (!error.response) {
      Alert.alert("Error", "please try again later or contact us", [
        {text: 'OK', onPress: () => {}},
      ])

      //TODO: send error
      return Promise.reject(error)
    }

    resp = error.response
    // console.debug("url", error.request.responseURL)
    if (
      resp.status == 400 && !!resp?.data?.details && resp?.data?.details.length > 0
    ) {
      Alert.alert("Failed", resp.data.details[0].message, [
        {text: 'OK', onPress: () => {}},
      ]);
      return Promise.resolve(resp)
    } else if (resp.status == 401 ) {
      router.replace("/login")
      return Promise.resolve(resp)
    } else if ( resp.status = 428 && resp?.data?.require === "TO_MAIN" ) {
      router.replace("/home")
    } else if ( resp.status = 428 && resp?.data?.require === "LOGIN" ) {
      router.replace("/login")
      return Promise.resolve(resp)
    } else if ( resp.status = 428 && resp?.data?.require === "VERIFY_EMAIL" ) {
      router.replace("/verify_email")
      return Promise.resolve(resp)
    } else if ( resp.status = 428 && resp?.data?.require === "CREATE_USER" ) {
      router.replace("/onboarding")
      return Promise.resolve(resp)
    }

    Alert.alert("Error", "please try again later or contact us", [
      {text: 'OK', onPress: () => {}},
    ])
  
    //TODO: send error log, error.message, error, error.request, error.response, user
    return resp.status == 500 ? Promise.resolve(resp) : Promise.reject(error)
  }
)

localLogout = async () => {
  await AsyncStorage.removeItem(SessionIdKey)
}

const hasLocalSession = async () => {
  try {
    return await AsyncStorage.getItem(SessionIdKey)
  } catch (e) {
    return null
  }
}

const authService = new AuthService(localLogout, httpFormClient)
const transactionService = new TransactionService(httpFormClient, httpClient)
export {
  authService,
  authPayload,
  hasLocalSession,
  transactionService
}