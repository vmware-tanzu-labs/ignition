import React from 'react'
import PropTypes from 'prop-types'
import Home from './home'
import Forbidden from './forbidden'
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
        }
      })
      .then(info => { this.setState({info: info}) })
    window
      .fetch('/api/v1/profile', {
        credentials: 'same-origin'
      })
      .then(response => {
        if (response.ok) {
          return response.json()
        }
      })
      .then(profile => { this.setState({ profile: profile }) })
  }

  render () {
    if (this.state.info && this.state.profile) {
      return (<Home info={this.state.info} profile={this.state.profile} />)
    } else if (this.state.profile) {
      return (<Forbidden profile={this.state.profile} />)
    } else {
      return (<h1>Loading...</h1>)
    }
  }
}

App.propTypes = {
  testing: PropTypes.bool
}

export default withRoot(App)
