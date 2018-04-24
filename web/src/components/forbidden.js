import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from 'material-ui/styles'
import Button from 'material-ui/Button'

import bkgd403 from './../../images/bkgd_403.svg'

const styles = theme => ({
  button: {
    backgroundColor: '#007D69',
    color: 'white',
    width: '27vh',
    margin: theme.spacing.unit,
    '&:hover': {
      backgroundColor: '#007363'
    }
  },
  forbidden: {
    backgroundImage: `url("${bkgd403}")`,
    backgroundRepeat: 'no-repeat',
    backgroundPosition: 'center',
    backgroundSize: 'cover',
    height: '100vh',
    alignItems: 'center',
    display: 'flex',
    justifyContent: 'center'
  },
  text: {
    fontFamily: 'Roboto, Helvetica, Arial, sans-serif',
    fontSize: '1.5rem',
    fontWeight: 'lighter',
    height: '100vh',
    width: '80vh',
    paddingTop: '40px',
    color: 'white',
    textAlign: 'center'
  }
})

class Forbidden extends React.Component {
  handleTryAgainButtonClick = async () => {
    window.location.replace('/logout')
  }

  render () {
    const { classes } = this.props
    const email = this.props.profile.Email
    return (
      <div className={classes.forbidden}>
        <div className={classes.text}>
          <p>
            <div>
              You&apos;ve attempted to sign in with {email} which
              does not grant you access.
            </div>
            <div>
              Please sign in with your company email account.
            </div>
          </p>
          <Button
            size="large"
            variant="raised"
            className={this.props.classes.button}
            onClick={this.handleTryAgainButtonClick}
          >
            Try Again
          </Button>
        </div>
      </div>
    )
  }
}

Forbidden.defaultProps = {
  profile: {
    Email: 'an unknown email'
  }
}

Forbidden.propTypes = {
  classes: PropTypes.object.isRequired,
  testing: PropTypes.bool,
  profile: PropTypes.object
}

export default withStyles(styles)(Forbidden)
