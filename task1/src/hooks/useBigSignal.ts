import { useState } from 'react'

type BigSignal<T extends Record<string, any>> = T & BigSignalClass<T>

class BigSignalClass<T extends Record<string, any>> {
  [key: string]: any
  set(newState: T) {
    for (const key in this) {
      if (this.hasOwnProperty(key) && newState.hasOwnProperty(key)) {
        if (typeof this[key] !== typeof newState[key]) throw new Error('Object structure does not match')
        if (isObj(this[key])) {
          this[key].set(newState[key])
          continue
        }
        this[key] = newState[key] as any
      }
    }
  }
}

function isObj(obj: object): boolean {
  return typeof obj === 'object' && !Array.isArray(obj) && obj !== null
}

export default function useBigSignal<T extends Record<string, any>>(initialState: T): BigSignal<T> {
  if (!isObj(initialState)) throw new Error('useBigSignal only accepts objects as an argument')
  const obj = new BigSignalClass<T>()
  for (const key in initialState) {
    if (initialState.hasOwnProperty(key)) {
      if (isObj(initialState[key])) {
        obj[key] = useBigSignal(initialState[key])
        continue
      }
      const [state, setState] = useState(initialState[key])
      const privateKey = `#${key}`
      Object.defineProperty(obj, privateKey, {
        value: state,
        writable: true,
        configurable: true,
        enumerable: false,
      })
      Object.defineProperty(obj, key, {
        get() {
          return this[privateKey]
        },
        set(value) {
          setState(value)
        },
        enumerable: true,
      })
    }
  }
  return obj as BigSignal<T>
}
