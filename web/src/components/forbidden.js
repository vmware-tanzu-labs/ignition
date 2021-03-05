import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'
import { Redirect } from 'react-router-dom'

import bkgd403 from './../../images/bkgd_403.svg'

const styles = theme => ({
  forbidden: {
    backgroundImage: `url("${bkgd403}")`,
    backgroundRepeat: 'no-repeat',
    backgroundPosition: 'center',
    backgroundSize: 'cover',
    height: '100vh',
    display: 'flex',
    justifyContent: 'center'
  },
  info: {
    paddingTop: '5vh',
    textAlign: 'center',
    '@media screen and (orientation:landscape)': {
      // iphone 5/SE, 6/7/8 & galaxy S5
      '@media only screen and (max-height: 375px) and (max-width: 667px)': {
        paddingTop: '1vh'
      }
    }
  },
  button: {
    backgroundColor: '#007D69',
    color: 'white',
    width: '200px',
    margin: theme.spacing(1),
    fontWeight: 'bold',
    letterSpacing: '3px',
    boxShadow: '-5px 5px 3px rgba(0, 0, 0, 0.20)',

    '&:hover': {
      backgroundColor: '#007363'
    },
    '@media screen and (orientation:landscape)': {
      // iphone 5/SE, 6/7/8 & galaxy S5
      '@media only screen and (max-height: 375px) and (max-width: 667px)': {
        padding: '0',
        width: '150px',
        margin: theme.spacing(0.5),
        minHeight: '32px'
      }
    }
  },
  text: {
    fontFamily: 'Roboto, Helvetica, Arial, sans-serif',
    fontSize: '2rem',
    fontWeight: 'lighter',
    width: '90vw',
    margin: 'auto',
    color: 'white',
    '@media screen and (orientation:portrait)': {
      [theme.breakpoints.down('sm')]: {
        fontSize: '1.2rem',
        paddingBottom: '5px'
      },
      // iphone 5/SE, 6/7/8 & galaxy S5
      '@media only screen and (max-height: 668px) and (max-width: 375px)': {
        fontSize: '1rem'
      },
      [theme.breakpoints.up('sm')]: {
        fontSize: '2rem'
      }
    },
    '@media screen and (orientation:landscape)': {
      [theme.breakpoints.down('sm')]: {
        fontSize: '1rem',
        paddingBottom: '3px'
      },
      // iphone 5/SE, 6/7/8 & galaxy S5
      '@media only screen and (max-height: 375px) and (max-width: 667px)': {
        padding: '0'
      },
      [theme.breakpoints.up('lg')]: {
        fontSize: '2.5rem'
      }
    }
  }
})

class Forbidden extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      redirect: false
    }
  }

  handleTryAgainButtonClick = async () => {
    window.location.replace('/logout')
  }

  render() {
    const { classes } = this.props
    let email
    if (this.props.location.state !== undefined) {
      email = this.props.location.state.profile.Email
    } else {
      this.setState({ redirect: true })
    }
    if (this.state.redirect) {
      return (
        <Redirect
          to={{
            pathname: '/'
          }}
        />
      )
    } else {
      return (
        <div className={classes.forbidden}>
          <div className={classes.info}>
            <div>
              <div className={classes.text}>
                You&apos;ve attempted to sign in with {email} which does not
                grant you access.
              </div>
              <div className={classes.text}>
                Please sign in with your company email account.
              </div>
            </div>
            <Button
              size="large"
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
}

Forbidden.defaultProps = {
  location: {
    state: {
      profile: {
        Email: 'an unknown email'
      }
    }
  }
}

Forbidden.propTypes = {
  classes: PropTypes.object.isRequired,
  testing: PropTypes.bool,
  location: PropTypes.object
}

export default withStyles(styles)(Forbidden)
