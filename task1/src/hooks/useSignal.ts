import { useState } from 'react'
import React from 'react'

class Signal<T> {
  readonly #state: T
  readonly #setState: React.Dispatch<React.SetStateAction<T>>
  constructor([state, setState]: [T, React.Dispatch<React.SetStateAction<T>>]) {
    this.#state = state
    this.#setState = setState
  }
  get val(): T {
    return this.#state
  }
  set val(newValue: T) {
    this.#setState(newValue)
  }
}

export default function useSignal<T>(initialState: T | (() => T)): Signal<T> {
  return new Signal(useState<T>(initialState))
}
