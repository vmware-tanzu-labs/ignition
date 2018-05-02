import React from 'react'
import PropTypes from 'prop-types'
import Home from './home'
import withRoot from '../withRoot'

class App extends React.Component {
  constructor (props) {
    super(props)
    this.state = {}
  }

  componentDidMount () {
    if (this.props && this.props.testing) {
      return
    }
    window
      .fetch('/api/v1/info', {
        credentials: 'same-origin'
      })
      .then(response => {
        if (response.ok) {
          return response.json()
        } else if (response.status === 401) {
          window.location.replace('/login')
        } else if (response.status === 403) {
          window.location.replace('/403')
          // TODO:
          // need to send profile={this.state.profile} to /403
          // in order to see the email address that was attempted
        }
      })
      .then(info => {
        this.setState({ info: info })
      })
    window
      .fetch('/api/v1/profile', {
        credentials: 'same-origin'
      })
      .then(response => {
        if (response.ok) {
          return response.json()
        }
      })
      .then(profile => {
        this.setState({ profile: profile })
      })
  }

  render () {
    if (this.state.info && this.state.profile) {
      return <Home info={this.state.info} profile={this.state.profile} />
    } else {
      return <div />
    }
  }
}

App.propTypes = {
  testing: PropTypes.bool
}

export default withRoot(App)
