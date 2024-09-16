import { View, Text, TouchableOpacity, TextInput, Modal, SafeAreaView, ScrollView } from 'react-native'
import React, { useState } from 'react'
import Currencies from '../constants/currency'

const variants = {
  bordered: "border-2 rounded-2xl",
  flat: "border-b-2"
}

const CurrencyInputField = ({
  title,
  containerStyles,
  variant='bordered',
  inputContainerStyles,
  supportCurrency=Object.values(Currencies),
  value,
  isEditable,
  onCurrencyChange=()=>{}
}) => {
  const [visible, setVisible] = useState(false)
  const [selectedCurrency, setSelectedCurrency] = useState(null)
  const [amount, setAmount] = useState("")
  const [amt, setAmt] = useState({amount: 0.0, currency: ''})
  console.log(amt)

  useState(() => {
    onCurrencyChange(amt)
  }, amt)
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
            <Text className="font-semibold">{!!selectedCurrency ? `${selectedCurrency.unicodeIcon} ${selectedCurrency.isoCode}` : " ------- "} 🔽</Text>
          </TouchableOpacity>
          <TextInput
            value={amount}
            keyboardType='decimal-pad'
            placeholder="0.00"
            isEditable={isEditable}
            onChangeText={(e) => {
              if ( e.match(/(\.\d\d\d)|.*\..*\..*/) ) return
              setAmount(e)
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
            className="flex-1 h-full text-right px-2 font-semibold"
          />
        </View>
      </View>
      {/* Currency picker */}
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
              supportCurrency.map(v => (
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
    </>
  )
}

export default CurrencyInputField