import React from 'react'
import { shallow } from 'enzyme'
import AppBar from './../../src/components/app-bar'

test('app bar renders name when the profile is present', () => {
  const profile = {
    Name: 'Test User',
    Email: 'testuser@company.net',
    AccountName: 'corp\tester'
  }
  const appBar = shallow(<AppBar profile={profile} />)
  expect(appBar.html().includes('Test User')).toBe(true)
})
