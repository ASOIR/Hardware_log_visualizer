import { describe, it, expect } from 'vitest'
import { add } from '../math'

describe('math utilities', () => {
  it('should add two numbers properly', () => {
    expect(add(1, 2)).toBe(3)
  })
})
