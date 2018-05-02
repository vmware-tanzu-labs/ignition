import React from 'react'
import PropTypes from 'prop-types'
import Home from './home'
import withRoot from '../withRoot'
import { Redirect } from 'react-router-dom'

class App extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      forbidden: false,
      redirect: false
    }
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
          this.setState({ forbidden: true })
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
        if (this.state.forbidden) {
          this.setState({ redirect: true })
        }
      })
  }

  render () {
    if (this.state.redirect) {
      return (
        <Redirect
          to={{
            pathname: '/403',
            state: { profile: this.state.profile }
          }}
        />
      )
    }
    if (this.state.info && this.state.profile) {
      return <Home info={this.state.info} profile={this.state.profile} />
    } else {
      return <div />
    }
  }
}

App.propTypes = {
  testing: PropTypes.bool,
  router: PropTypes.func
}

export default withRoot(App)
