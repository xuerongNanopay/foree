import { View, Text, SafeAreaView, TouchableOpacity, FlatList, Image } from 'react-native'
import React, { useCallback, useEffect, useState } from 'react'
import { router, useFocusEffect, useNavigation } from 'expo-router'

import { icons } from '../../constants'
import { transactionService } from '../../service'
import string_util from '../../util/string_util'
import TxSummaryChip from '../../components/TxSummaryChip'

const TransactionTab = () => {
  const [selectedStatus, setSelectedStatus] = useState('All')
  const [page, setPage] = useState(0)
  const [count, setCount] = useState(1000)
  const [txs, setTxs] = useState([])
  const maxSize = 10
  const navigation = useNavigation()

  useFocusEffect(useCallback(() => {
    console.log("aaa")
    const controller = new AbortController()
    if ( page == 0 && selectedStatus == "All") {
      loadTransactions(controller.signal)
    } else {
      setSelectedStatus('All')
      setPage(0)
    }
    return () => {
      controller.abort()
    }
  }, []))
  
  const loadTransactions = (signal) => {
    const getTransactions = async() => {
      try {
        const resp = await transactionService.getTransactions({status: selectedStatus, offset:page*10, limit:maxSize}, {signal: signal})
        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("get transactions", resp.status, resp.data)
        } else {
          setTxs(resp.data.data)
        }
      } catch (e) {
        console.error("get transactions", e, e.response, e.response?.status, e.response?.data)
      }
    }
    const countTransactions = async() => {
      try {
        const resp = await transactionService.countTransactions({status: selectedStatus, offset:page*10, limit:maxSize}, {signal: signal})
        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("count transactions", resp.status, resp.data)
        } else {
          setCount(resp.data.data.count)
        }
      } catch (e) {
        console.error("count transactions", e, e.response, e.response?.status, e.response?.data)
      }
    }

    getTransactions()
    countTransactions()
  }

  useEffect(() => {
    const controller = new AbortController()
    loadTransactions(controller.signal)
    return () => {
      controller.abort()
    }
  },[page, selectedStatus])

  const statusChipItem = useCallback(({item}) => {
    const bgColor = selectedStatus === item.id ? `${item.selectBgColor}` : ""
    return (
      <TouchableOpacity 
        onPress={() => {
          setPage(_ => {
            return 0
          })
          setSelectedStatus(item.id)
        }}
        className={`p-2 border-[1px] ${item.borderColor} rounded-2xl mr-2 ${bgColor}`}
      >
        <Text className={`${item.textColor}`}>{item.id}</Text>
      </TouchableOpacity>
    )
  },[selectedStatus])

  const TxItem = ({index, item}) => {
    const tx = item
    return(
      <TouchableOpacity
        onPress={() => router.push(`/transaction/${tx.id}`)}
        className={`py-1 px-1 ${index%2===1 ? "bg-slate-200": ""}`}
      >
        <View className="mb-1 flex-row items-center justify-between">
          <Text className="font-semibold">{!!tx.destAccSummary ? string_util.formatStringWithLimit(tx.destAccSummary, 14) : ""}</Text>
          <Text className="font-semibold text-slate-600">${new Intl.NumberFormat("en", {minimumFractionDigits: 2}).format(tx.destAmount)}{!!tx.destCurrency ? ` ${tx.destCurrency}` : ''}</Text>
        </View>
        <View className="mb-1 flex-row items-center justify-between">
          <Text className="font-semibold">Total Amount</Text>
          <Text className="font-semibold text-slate-600">${new Intl.NumberFormat("en", {minimumFractionDigits: 2}).format(tx.totalAmount)}{!!tx.totalCurrency ? ` ${tx.totalCurrency}` : ''}</Text>
        </View>
        <View className="flex-row items-center justify-between">
          <Text className="italic text-slate-600">{tx.nbpReference}</Text>
          <TxSummaryChip
            status={tx.status}
          />
        </View>
      </TouchableOpacity>
    )
  }
  return (
    <SafeAreaView>
      <View className="h-full px-4 pt-4">
        <View className="pb-2 mb-4 border-b-[1px] border-slate-300">
          {/* Title */}
          <View className="flex flex-row items-center pb-2 mb-2 border-b-[1px] border-slate-300">
            <Text className="flex-1 font-pbold text-2xl">Transactions</Text>
            <TouchableOpacity
              onPress={()=> {router.push("/transaction/create")}}
              activeOpacity={0.7}
              className="bg-[#1A6B54] py-2 px-4 rounded-full"
              disabled={false}
            >
              <Text className="font-pextrabold text-white">Send</Text>
            </TouchableOpacity>
          </View>
          {/* Status */}
          <View className="flex flex-row items-center">
            <TouchableOpacity
              onPress={() => {setPage((page) => {
                if ( page === 0 ) {
                  loadTransactions()
                  return page
                } else return 0
              })}}
              className="border-[1px] border-slate-400 rounded-lg p-1"
            >
              <Image
                source={icons.renewable}
                className="w-[27px] h-[27px]"
                resizeMode='contain'
              />
            </TouchableOpacity>
            <FlatList
              className="flex-1 mx-2"
              data={transactionStatuses}
              renderItem={statusChipItem}
              keyExtractor={item => item.id}
              showsHorizontalScrollIndicator={false}
              horizontal={true}
            />
            {
              selectedStatus !== "All" ? 
              <TouchableOpacity
                onPress={()=> {setSelectedStatus("All")}}
                className="p-2 h-full flex flex-row items-center"
                disabled={false}
              >
              <Image 
                  source={icons.x}
                  className="w-3 h-3"
                  resizeMode='contain'
              />
              </TouchableOpacity> :
              <></>
            }
          </View>
          {/* Status Pagenation */}
          <View className="mt-2 flex flex-row items-center">
            <View className="flex-1">
              <View className="flex flex-row">
                <View className={`border-[1px] rounded-2xl p-2 ${transactionStatuses.find(x=>x.id===selectedStatus).borderColor} ${transactionStatuses.find(x=>x.id===selectedStatus).selectBgColor}`}>
                  <Text className={`font-psemibold ${transactionStatuses.find(x=>x.id===selectedStatus).textColor}`}>{selectedStatus}</Text>
                </View>
              </View>
            </View>
            <View className="flex flex-row items-center">
              <Text className="mr-2">{`${page*maxSize+1}-${page*(maxSize)+maxSize}`}</Text>
              <TouchableOpacity
                onPress={()=> {
                  setPage((page) => {
                    return page > 0 ? page-1 : 0
                  })
                }}
                activeOpacity={0.7}
                disabled={page==0}
                className="mr-2 p-2"
              >
                <Image
                  source={icons.leftArrowDark}
                  className="w-[15px] h-[15px]"
                  resizeMode='contain'
                />
              </TouchableOpacity>
              <TouchableOpacity
                onPress={()=> {
                  setPage((page) => {
                    return (page+1)*maxSize > count ? page : page+1
                  })
                }}
                activeOpacity={0.7}
                disabled={false}
                className="mr-2 p-2"
              >
                <Image
                  source={icons.rightArrowDark}
                  className="w-[15px] h-[15px]"
                  resizeMode='contain'
                />
              </TouchableOpacity>
            </View>
          </View>
        </View>
        <FlatList
          data={txs}
          keyExtractor={tx => tx.id}
          renderItem={TxItem}
          showsVerticalScrollIndicator={false}
          showsHorizontalScrollIndicator={false}
          ListEmptyComponent={
            <View>
              <Text className="text-center font-pbold text-xl text-slate-600 mt-44">⛔ No Transactions</Text>
            </View>
          }
        />
      </View>
    </SafeAreaView>
  )
}

const transactionStatuses = [
  {
    id:"All",
    borderColor: "border-slate-800",
    textColor: "text-slate-800",
    selectBgColor: "bg-slate-200",
  },
  { 
    id:"Initial",
    borderColor: "border-purple-800",
    textColor: "text-purple-800",
    selectBgColor: "bg-purple-200",
  },
  { 
    id:"Await Payment",
    borderColor: "border-yellow-800",
    textColor: "text-yellow-800",
    selectBgColor: "bg-yellow-200",
  },
  { 
    id:"In Progress",
    borderColor: "border-purple-800",
    textColor: "text-purple-800",
    selectBgColor: "bg-purple-200",
  },
  { 
    id:"Completed",
    borderColor: "border-[#005a32]",
    textColor: "text-[#005a32]",
    selectBgColor: "bg-[#c7e9c0]",
  },
  { 
    id:"Cancelled",
    borderColor: "border-red-800",
    textColor: "text-red-800",
    selectBgColor: "bg-red-200",
  },
  { 
    id:"Ready To Pickup",
    borderColor: "border-yellow-800",
    textColor: "text-yellow-800",
    selectBgColor: "bg-yellow-200",
  },
  { 
    id:"Refunding",
    borderColor: "border-purple-800",
    textColor: "text-purple-800",
    selectBgColor: "bg-purple-200",
  }
]

export default TransactionTab