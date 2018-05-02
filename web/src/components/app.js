import React from 'react'
import PropTypes from 'prop-types'
import Home from './home'
import Forbidden from './forbidden'
import NotFound from './notfound'
import withRoot from '../withRoot'

class App extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      forbidden: false,
      notFound: false
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
      })
    if (window.location.pathname !== '/') {
      this.setState({ notFound: true })
    }
  }

  render () {
    if (this.state.forbidden) {
      return <Forbidden profile={this.state.profile} />
    } else if (this.state.notFound) {
      return <NotFound />
    } else if (this.state.info && this.state.profile) {
      return <Home info={this.state.info} profile={this.state.profile} />
    } else {
      return <div>&nbsp;</div>
    }
  }
}

App.propTypes = {
  testing: PropTypes.bool
}

export default withRoot(App)
