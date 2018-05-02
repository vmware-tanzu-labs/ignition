import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from 'material-ui/styles'

import bkgd404 from './../../images/bkgd_404_v1.svg'

const styles = theme => ({
  notFound: {
    backgroundImage: `url("${bkgd404}")`,
    backgroundRepeat: 'no-repeat',
    backgroundPosition: 'center',
    backgroundSize: 'cover',
    height: '100vh'
  }
})

class NotFound extends React.Component {
  render () {
    const { classes } = this.props
    return <div className={classes.notFound} />
  }
}

NotFound.propTypes = {
  classes: PropTypes.object.isRequired
}

export default withStyles(styles)(NotFound)
