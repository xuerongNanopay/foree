import axios from 'axios'

import AuthService from "./auth_service"
import authPayload from "./auth_payload"
import { Alert } from 'react-native'
import { router } from 'expo-router'
import AsyncStorage from '@react-native-async-storage/async-storage'

const SessionIdKey = "SESSION_ID"
// Config axios
axios.defaults.baseURL = 'http://192.168.2.30:8080/app/v1'
axios.interceptors.request.use(
  async (config) => {
    try {
      const session = await AsyncStorage.getItem(SessionIdKey)
      if ( !!session ) {
        console.log("Request Session: ", session)
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
axios.interceptors.response.use(
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
      Alert.alert("Error", "please try later", [
        {text: 'OK', onPress: () => {}},
      ])

      //TODO: send error
      return Promise.reject(error)
    }

    resp = error.response
    console.log("url", error.request.responseURL)
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
    }

    Alert.alert("Error", "please try later", [
      {text: 'OK', onPress: () => {}},
    ])
  
    //TODO: send error log, error.message, error, error.request, error.response, user
    return Promise.reject(error)
  }
)

const authService = new AuthService()

export {
  authService,
  authPayload
}