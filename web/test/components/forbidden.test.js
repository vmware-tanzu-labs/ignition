import React from 'react'
import { shallow } from 'enzyme'
import Forbidden from './../../src/components/forbidden'

test('forbidden defaults email address to `an unknown email`', () => {
  const forbidden = shallow(<Forbidden />)
  expect(forbidden.html().includes('an unknown email')).toBe(true)
})

test('forbidden renders email address when present', () => {
  const location = {
    state: {
      profile: {
        Email: 'sneal@example.com'
      }
    }
  }
  // this.props.location.state.profile.Email
  const forbidden = shallow(<Forbidden />)
  forbidden.setProps({ location: location })
  expect(forbidden.html().includes('sneal@example.com')).toBe(true)
})
