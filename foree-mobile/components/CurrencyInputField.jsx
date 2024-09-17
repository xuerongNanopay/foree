import { View, Text, TouchableOpacity, TextInput, Modal, SafeAreaView, ScrollView } from 'react-native'
import React, { useEffect, useState } from 'react'
import Currencies from '../constants/currency'

const variants = {
  bordered: "border-2 rounded-2xl",
  flat: "border-b-2"
}

// Will introduce cycle, if mix source and dest.
// if use string as amount, the problem can solove
const CurrencyInputField = ({
  title,
  containerStyles,
  variant='bordered',
  inputContainerStyles,
  placeholder="0.00",
  inputStyles,
  supportCurrencies=Object.values(Currencies),
  value=0.0,
  editable=true,
  onCurrencyChange=(e)=>{},
  errorMessage
}) => {
  const [visible, setVisible] = useState(false)
  const [selectedCurrency, setSelectedCurrency] = useState(supportCurrencies.length === 1 ? supportCurrencies[0] : null)
  const [amtString, setAmtString] = useState(!value ? "" : value.toFixed(2))
  const [amt, setAmt] = useState(supportCurrencies.length === 1 ? {amount: 0.0, currency: supportCurrencies[0].isoCode} : {amount: 0.0, currency: ''})

  useEffect(() => {
    onCurrencyChange(amt)
  }, [amt])

  useEffect(() => {
    setAmt({
      ...amt,
      amount: value
    })
    setAmtString(!value ? "" : value.toFixed(2))
  }, [value])

  return (
    <>
      <View className={`${containerStyles}`}>
        <Text className="test-base test-gray-100 font-pmedium mb-2">{title}</Text>
        <View
          className={`
            flex flex-row items-center
            h-12 ${variants[variant] ?? variants.bordered}
            border-slate-400 focus:border-secondary-200 ${inputContainerStyles}
          `}
        >
          <TouchableOpacity
            onPress={() => {
              setVisible(true)
            }}
            activeOpacity={0.7}
            className="bg-slate-200 h-full px-2 rounded-l-2xl flex justify-center"
          >
            <Text className="font-semibold">
              {!!selectedCurrency ? `${selectedCurrency.unicodeIcon} ${selectedCurrency.isoCode}` : " ------- "} {supportCurrencies.length === 1 ? "" : "🔽"}
            </Text>
          </TouchableOpacity>
          <TextInput
            className={`flex-1 h-full text-right px-2 font-semibold ${inputStyles}`}
            value={amtString}
            keyboardType='decimal-pad'
            placeholder={placeholder}
            editable={editable}
            onChangeText={(e) => {
              if ( !!e.match(/(\.\d\d\d)|.*\..*\..*/) ) {
                setAmt({...amt})
                return
              }
              setAmtString(e)
              n = parseFloat(e)
              if ( isNaN(n) ) {
                setAmt({
                  ...amt,
                  amount: 0.0
                })
              } else {
                setAmt({
                  ...amt,
                  amount: n
                })
              }
            }}
          />
        </View>
        {
          !!errorMessage ?
          <View>
            <Text className="mt-2 text-red-600">{errorMessage}</Text>
          </View> : null
        }
      </View>
      {/* Currency picker */}
      {
        supportCurrencies.length === 1 ? <></> :
        <Modal
          visible={visible}
          animationType='slide'
        >
          <SafeAreaView className="h-full flex flex-col">
            <View className="flex flex-row items-center border-b-[1px] border-slate-400">
              <Text
                onPress={() => {
                  setVisible(false)
                }}
                className="py-2 px-4 text-2xl font-bold text-slate-600"
              >
                &#8592;
              </Text>
              <Text
                className="font-psemibold text-xl text-slate-600"
              >
                Choose a currency
              </Text>
            </View>
            <ScrollView
              className="flex-1"
              showsVerticalScrollIndicator={false}
            >
              {
                supportCurrencies.map(v => (
                  <TouchableOpacity
                    onPress={() => {
                      
                      setSelectedCurrency(v)
                      setAmount("")
                      setAmt({
                        currency: v.isoCode,
                        amount: 0.0
                      })
                      setVisible(false)
                    }}
                    key={v.isoCode}
                    className="mx-2 py-2 border-b-[1px] border-slate-400"
                  >
                    <Text className="font-semibold text-2xl">{v.unicodeIcon} {v.isoCode}</Text>
                  </TouchableOpacity>
                ))
              }
            </ScrollView>
          </SafeAreaView>
        </Modal>
      }
    </>
  )
}

export default CurrencyInputField