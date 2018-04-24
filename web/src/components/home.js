import React from 'react'
import PropTypes from 'prop-types'
import AppBar from './app-bar'
import Body from './body'

const Home = (props) => {
  const {profile, info} = props
  return (
    <div>
      <AppBar profile={profile} />
      <Body info={info} />
    </div>
  )
}

Home.propTypes = {
  profile: PropTypes.object,
  info: PropTypes.object
}

export default Home
