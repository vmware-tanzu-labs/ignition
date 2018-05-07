import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from 'material-ui/styles'
import Button from 'material-ui/Button'

import bkgd404 from './../../images/bkgd_404_v2.svg'

const styles = theme => ({
  page: {
    backgroundImage: `url("${bkgd404}")`,
    backgroundRepeat: 'no-repeat',
    backgroundPosition: 'center',
    backgroundSize: 'cover',
    height: '100vh',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'flex-start'
  },
  button: {
    backgroundColor: '#007D69',
    color: 'white',
    width: '200px',
    margin: theme.spacing.unit,
    marginTop: '15vh',
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
        margin: 0.5 * theme.spacing.unit,
        marginTop: '10vh',
        minHeight: '32px'
      }
    }
  }
})

class NotFound extends React.Component {
  handleTryAgainButtonClick = async () => {
    window.location.replace('/')
  }

  render () {
    const { classes } = this.props
    return (
      <div className={classes.page}>
        <Button
          size="large"
          variant="raised"
          className={this.props.classes.button}
          onClick={this.handleTryAgainButtonClick}
        >
          Phone Home
        </Button>
      </div>
    )
  }
}

NotFound.propTypes = {
  classes: PropTypes.object.isRequired
}

export default withStyles(styles)(NotFound)
