import React from 'react'
import { shallow } from 'enzyme'
import Forbidden from './../../src/components/forbidden'

test('forbidden defaults email address to `an unknown email`', () => {
  const forbidden = shallow(<Forbidden />)
  expect(forbidden.html().includes('an unknown email')).toBe(true)
})

test('forbidden renders email address when present', () => {
  const forbidden = shallow(<Forbidden email='sneal@example.com' />)
  expect(forbidden.html().includes('sneal@example.com')).toBe(true)
})
