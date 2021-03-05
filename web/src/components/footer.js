import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'

const styles = theme =>  ({
  root: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between'
  },
  linksContainer: {
    display: 'flex',
    alignItems: 'center',
    '& a': {
      color: theme.palette.primary.dark,
      fontSize: '12px',
      fontWeight: '600',
      letterSpacing: '0.5px',
      textDecoration: 'none',
      padding: '8px'
    }
  },
  img: {
    height: '60px',
    padding: '12px'
  }
})

const Footer = props => {
  const { classes, links, logoURL } = props
  return (
    <div className={classes.root}>
      <div className={classes.linksContainer}>
        {links.map(l => (
          <a href={l.url} key={l.text}>
            {l.text}
          </a>
        ))}
      </div>
      <img className={classes.img} src={logoURL} alt="logo" />
    </div>
  )
}

Footer.propTypes = {
  classes: PropTypes.object,
  links: PropTypes.arrayOf(
    PropTypes.shape({
      text: PropTypes.string.isRequired,
      url: PropTypes.string.isRequired
    })
  ),
  logoURL: PropTypes.string
}

export default withStyles(styles)(Footer)
