import React from 'react'
import PropTypes from 'prop-types'
import classNames from 'classnames';
import { withStyles } from 'material-ui/styles'
import Button from 'material-ui/Button'

import milkyWay from './../../images/bkgd_milky-way_full.svg'


const styles = theme => ({
  body: {
    fontFamily: 'Roboto, Helvetica, Arial, sans-serif',
    fontWeight: 'lighter',
    marginTop: '68px'
  },
  button: {
    margin: theme.spacing.unit,
    '&:hover': {
      backgroundColor: '#007363'
    },
  }
})

class Forbidden extends React.Component {
  constructor (props) {
    super(props)
  }

  render () {
    const { classes } = this.props
    return (
      <div className={classes.body}>
        403
      </div>
    )
  }
}

Forbidden.propTypes = {
  classes: PropTypes.object.isRequired
}

export default withStyles(styles)(Forbidden)
