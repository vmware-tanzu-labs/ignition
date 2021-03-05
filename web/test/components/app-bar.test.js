import React from 'react'
import { shallow } from 'enzyme'
import AppBar from './../../src/components/app-bar'

test('app bar renders email when the profile is present', () => {
  const profile = {
    Email: 'testuser@company.net',
    AccountName: 'corp\tester'
  }
  const appBar = shallow(<AppBar profile={profile} />)
  expect(appBar.html().includes('testuser@company.net')).toBe(true)
})
