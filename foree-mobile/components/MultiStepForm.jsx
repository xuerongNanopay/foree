import { View, Text, ScrollView } from 'react-native'
import React, { useState, useEffect, useRef } from 'react'
import CustomButton from './CustomButton'

const MultiStepForm = ({
  containerStyle,
  showProgress=true,
  submitTintTitle = 'Submit',
  onSumbit = () => {},
  steps = () => [],
  isSubmitting=false
}) => {
  const scrollRef = useRef()
  const formStep = steps()
  const [curIdx, setCurIdx] = useState(0)
  const [isFirst, setIsFirst] = useState(curIdx===0)
  const [isLast, setIsLast] =  useState(curIdx===formStep.length-1)
  const [progress, setProgress] = useState((curIdx+1)/formStep.length)
  const [progressCss, setProgressCss] = useState('-100%')

  useEffect(() => {
    setIsFirst(curIdx===0)
    setIsLast(curIdx===formStep.length-1)
    setProgress((curIdx+1)/formStep.length)
    scrollRef.current?.scrollTo({
      y: 0,
      animated: false,
    })
  }, [curIdx])

  useEffect(()=>{
    setProgressCss('-' + Math.round((1-progress)*100) + '%')
  }, [progress])

  return (
    <View className={`${containerStyle}`}>
      { showProgress ? 
        (
          <View
            className="h-[3px] bg-slate-300 relative"
          >
            <View 
              className={`border-t-[3px] border-secondary top-[0px]`}
              style={{
                left: progressCss
              }}
            />
          </View>
        ) : null
      }
      <View className="h-full px-2 pb-4">
        <ScrollView
          className=""
          automaticallyAdjustKeyboardInsets
          ref={scrollRef}
        >
          {
            formStep[curIdx].formView()
          }
        </ScrollView>
        <View className="mt-2 w-full flex flex-row justify-between">
          {
            !isFirst ?
            <CustomButton
              title={"< Previous"}
              isLoading={isSubmitting}
              handlePress={()=>{setCurIdx(curIdx-1)}}
              containerStyles={"w-[49%]"}
              variant="bordered"
                // isLoading={isSubmitting}
            /> : null
          }
          {
            !isLast ? 
            <CustomButton
              title={"Next >"}
              isLoading={isSubmitting}
              handlePress={()=>{setCurIdx(curIdx+1)}}
              containerStyles={isFirst ? "w-[100%]" : "w-[49%]"}
                // isLoading={isSubmitting}
            /> : null
          }
          {
            isLast ? 
            <CustomButton
              title={submitTintTitle}
              isLoading={isSubmitting}
              handlePress={onSumbit}
              containerStyles={isFirst ? "w-[100%]" : "w-[49%]"}
                // isLoading={isSubmitting}
            /> : null
          }
        </View>
      </View>
    </View>
  )
}

export default MultiStepForm