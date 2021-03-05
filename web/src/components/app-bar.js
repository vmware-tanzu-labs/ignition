import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import AppBar from '@material-ui/core/AppBar'
import Button from '@material-ui/core/Button'
import Toolbar from '@material-ui/core/Toolbar'
import Typography from '@material-ui/core/Typography'
import IconButton from '@material-ui/core/IconButton'
import AccountCircle from '@material-ui/icons/AccountCircle'
import Menu from '@material-ui/core/Menu'
import MenuItem from '@material-ui/core/MenuItem'
import { getOrgUrl } from './../org'

import ignitionLogo from './../../images/ignition.svg'

const styles = theme => ({
  root: {
    display: 'flex',
    position: 'sticky',
    top: 0,
    left: 'auto',
    right: 0,
    zIndex: 999
  },
  logoContainer: {
    background: '#F2F0F1',
    padding: 0
  },
  logo: {
    height: '64px',
    padding: '8px 24px'
  },
  userContainer: {
    display: 'flex',
    flexGrow: 1,
    alignItems: 'center',
    '@media screen and (orientation:portrait)': {
      [theme.breakpoints.down('sm')]: {
        justifyContent: 'flex-end'
      }
    }
  },
  name: {
    flexGrow: 1,
    paddingRight: '10px',
    '@media screen and (orientation:portrait)': {
      [theme.breakpoints.down('sm')]: {
        display: 'none'
      }
    }
  },
  icon: {
    flexShrink: 1
  },
  button: {
    margin: theme.spacing.unit,
    backgroundColor: theme.palette.primary.dark,
    '&:hover': {
      backgroundColor: theme.palette.primary.main
    },
    color: 'white',
    letterSpacing: '0.5px'
  },
  menuButton: {
    marginLeft: -12,
    marginRight: 20
  }
})

class MenuAppBar extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      anchorEl: null
    }
  }

  handleButton = async () => {
    const url = await getOrgUrl()
    if (url) {
      window.location = url
    }
  }

  handleMenu = event => {
    this.setState({ anchorEl: event.currentTarget })
  }

  handleClose = () => {
    this.setState({ anchorEl: null })
  }

  handleLogout = (e, location = window.location) => {
    this.setState({ anchorEl: null })
    if (this.props && this.props.testing) {
      return
    }
    location.replace('/logout')
  }

  render () {
    const { classes } = this.props
    const { anchorEl } = this.state
    const open = Boolean(anchorEl)
    const name = this.props.profile.Name

    return (
      <div className={classes.root}>
        <AppBar color="white">
          <Toolbar disableGutters={true}>
            <div className={classes.logoContainer}>
              <img
                className={classes.logo}
                src={ignitionLogo}
                alt="ignition logo"
              />
            </div>
            <div className={classes.userContainer}>
              <Typography
                variant="subheading"
                color="primary"
                align="right"
                className={classes.name}
              >
                {`Welcome, ${name}`}
              </Typography>
              <Button
                className={classes.button}
                size="large"
                variant="raised"
                onClick={this.handleButton}
                alt=""
              >
                My Org
              </Button>
              <IconButton
                aria-owns={open ? 'menu-appbar' : null}
                aria-haspopup="true"
                onClick={this.handleMenu}
                color="primary"
                className={classes.icon}
              >
                <AccountCircle />
              </IconButton>
              <Menu
                id="menu-appbar"
                anchorEl={anchorEl}
                anchorOrigin={{
                  vertical: 'top',
                  horizontal: 'right'
                }}
                transformOrigin={{
                  vertical: 'top',
                  horizontal: 'right'
                }}
                open={open}
                onClose={this.handleClose}
              >
                <MenuItem onClick={this.handleLogout}>Logout</MenuItem>
              </Menu>
            </div>
          </Toolbar>
        </AppBar>
      </div>
    )
  }
}

MenuAppBar.propTypes = {
  classes: PropTypes.object.isRequired,
  testing: PropTypes.bool,
  profile: PropTypes.object
}

export default withStyles(styles)(MenuAppBar)
