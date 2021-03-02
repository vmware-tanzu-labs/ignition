import React from 'react'
import PropTypes from 'prop-types'
import classNames from 'classnames'
import { withStyles } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'
import Footer from './footer'
import { getOrgUrl } from '../org'

import milkyWay from './../../images/bkgd_milky-way_full.svg'
import deepSpace from './../../images/bkgd_lvl2_deep-space.svg'
import icePlanet from './../../images/bkgd_lvl3_ice-planet.svg'

import rocketMan from './../../images/frgd_rocket-man.svg'
import moonMan from './../../images/frgd_moon-man.svg'
import pewPew from './../../images/frgd_pewpew-man2.svg'

import step1 from './../../images/step-1.svg'
import step2 from './../../images/step-2.svg'
import step3 from './../../images/step-3.svg'
import pivotalLogo from './../../images/pivotal.png'

const makeSpeechBubbleClass = (theme, bgColor, fgColor) => ({
  position: 'relative', // so we can overlap the button

  padding: '24px',
  borderRadius: '15px',
  backgroundColor: bgColor,
  color: fgColor,

  fontSize: '1.75rem',
  height: 'auto',
  width: '40vw',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',

  '&:before': {
    content: '""',
    width: '0px',
    height: '0px',
    position: 'absolute',
    borderLeft: `100px solid ${bgColor}`,
    borderRight: '100px solid transparent',
    borderTop: `25px solid ${bgColor}`,
    borderBottom: '25px solid transparent',
    right: '-175px',
    top: '75px',
    [theme.breakpoints.down('sm')]: {
      display: 'none'
    }
  },
  // screen over 1920
  [theme.breakpoints.up('xl')]: {
    fontSize: '2.75rem'
  },
  // screen less 600
  '@media screen and (orientation:portrait)': {
    [theme.breakpoints.down('sm')]: {
      width: '90vw',
      fontSize: '1.50rem'
    }
  },
  '@media screen and (orientation:landscape)': {
    [theme.breakpoints.down('sm')]: {
      width: '95vw',
      maxHeight: '60vh',
      fontSize: '1.25rem',
      '& p': {
        marginTop: '0'
      }
    }
  }
})

const speech1Background = '#083c61'
const speech2Background = '#9bd2d2'
const greenButton = '#007D69'

const footerLinks = [
  { text: 'Copyright', url: '' },
  { text: 'Terms', url: '' },
  { text: 'Contact', url: '' }
]

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
    }
  },
  // for buttons that overlap the bottom of a speech bubble
  speechButton: {
    position: 'absolute',
    bottom: -4 * theme.spacing.unit,
    backgroundColor: greenButton,
    color: 'white',
    height: 7 * theme.spacing.unit,
    fontWeight: 'bold',
    letterSpacing: '3px',
    boxShadow: '-5px 5px 3px rgba(0, 0, 0, 0.29)'
  },
  welcomeSpeech: makeSpeechBubbleClass(theme, speech1Background, 'white'),
  spacesSpeech: makeSpeechBubbleClass(theme, speech2Background, 'black'),

  // CTA 1: Welcome
  // container
  ctaWelcome: {
    backgroundImage: `url("${milkyWay}")`,
    backgroundRepeat: 'no-repeat',
    backgroundPosition: 'center',
    backgroundColor: '#00253e',
    backgroundSize: 'cover',
    height: '80vh',
    padding: 6 * theme.spacing.unit,
    alignItems: 'center',

    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'center',
    // screen less 600
    '@media screen and (orientation:portrait)': {
      [theme.breakpoints.down('sm')]: {
        padding: 2 * theme.spacing.unit,
        height: '80vh',
        alignItems: 'center'
      }
    },
    '@media screen and (orientation:landscape)': {
      [theme.breakpoints.down('sm')]: {
        height: '65vh'
      }
    },
    // screen over 1920
    [theme.breakpoints.up('xl')]: {
      height: '89vh'
    }
  },
  // character
  rocketMan: {
    backgroundImage: `url("${rocketMan}")`,
    height: '450px',
    width: '450px',
    backgroundRepeat: 'no-repeat',
    // screen over 1920
    [theme.breakpoints.up('xl')]: {
      height: '600px',
      width: '600px'
    },
    // screen less 600
    [theme.breakpoints.down('sm')]: {
      display: 'none'
    }
  },

  // CTA 2: three steps
  // container
  ctaSteps: {
    backgroundImage: `url("${deepSpace}")`,
    backgroundRepeat: 'no-repeat',
    backgroundPosition: 'center',
    backgroundSize: 'cover',
    backgroundColor: '#00253e',
    color: 'white',
    padding: 6 * theme.spacing.unit,
    height: '77vh',
    fontSize: '32px',
    display: 'flex',
    flexDirection: 'row',
    '@media screen and (orientation:portrait)': {
      [theme.breakpoints.down('sm')]: {
        width: '100vw',
        height: '120vh',
        padding: 2 * theme.spacing.unit,
        fontSize: '2rem',
        paddingBottom: '20vh',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'space-around'
      }
    },
    '@media screen and (orientation:landscape)': {
      [theme.breakpoints.down('sm')]: {
        height: '100vh',
        flexDirection: 'row',
        justifyContent: 'space-around'
      }
    }
  },
  // character
  pewPew: {
    backgroundImage: `url("${pewPew}")`,
    backgroundRepeat: 'no-repeat',
    backgroundPosition: 'left',
    backgroundSize: 'auto',
    height: '450px',
    width: '450px',
    flexShrink: 0,
    // screen over 1920
    [theme.breakpoints.up('xl')]: {
      height: '600px',
      width: '600px'
    },
    // screen less 600
    [theme.breakpoints.down('sm')]: {
      display: 'none'
    }
  },
  step: {
    textAlign: 'center',
    borderTop: '9px solid #FFC712',
    marginTop: '64px',
    flexGrow: '1',
    // screen over 1920
    [theme.breakpoints.up('xl')]: {
      marginTop: '87px'
    },
    // screen less 600
    [theme.breakpoints.down('sm')]: {
      borderTop: 'none',
      flexGrow: '0'
    },
    '& p': {
      maxWidth: '258px',
      margin: 'auto',
      // screen less 600
      [theme.breakpoints.down('sm')]: {
        fontSize: '1.5rem'
      }
    }
  },
  stepImage: {
    height: '146px',
    marginTop: '-82px',
    // screen over 1920
    [theme.breakpoints.up('xl')]: {
      height: '240px',
      marginTop: '-148px'
    },
    // screen less 600
    '@media screen and (orientation:portrait)': {
      [theme.breakpoints.down('sm')]: {
        height: '25vw'
      }
    },
    '@media screen and (orientation:landscape)': {
      [theme.breakpoints.down('sm')]: {
        height: '36vh'
      }
    }
  },

  // CTA 3: spaces overview
  ctaSpaces: {
    backgroundImage: `url("${icePlanet}")`,
    backgroundRepeat: 'no-repeat',
    backgroundPosition: 'center',
    backgroundSize: 'cover',
    height: '80vh',
    padding: 6 * theme.spacing.unit,
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-around',
    alignItems: 'center',
    // screen less 600
    [theme.breakpoints.down('sm')]: {
      height: '85vh',
      alignItems: 'center'
    }
  },
  moonMan: {
    backgroundImage: `url("${moonMan}")`,
    height: '450px',
    width: '400px',
    backgroundRepeat: 'no-repeat',
    // screen over 1920
    [theme.breakpoints.up('xl')]: {
      height: '600px',
      width: '600px'
    }
  },
  temporary: {
    alignItems: 'center'
  },
  emphasis: {
    fontWeight: 600
  }
})

class Body extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      orgUrl: ''
    }
  }

  handleOrgButtonClick = async () => {
    // TODO: show spinner
    const url = await getOrgUrl()
    if (url) {
      this.setState({ orgUrl: url })
      window.location = url
    }
  }

  renderWelcomeInfo () {
    const { classes } = this.props

    return (
      <div className={classes.ctaWelcome}>
        <div>
          <div className={classes.welcomeSpeech}>
            <p>
              <span className={classes.emphasis}>
                {this.props.info.CompanyName}
              </span>{' '}
              is giving you a playground to push (deploy) apps and experiment.
              Tanzu Application Service (TAS) uses{' '}
              <span className={classes.emphasis}>orgs</span> to organize things.
            </p>
            <p>
              Orgs contain <span className={classes.emphasis}>spaces</span>, and
              each space can host <span className={classes.emphasis}>apps</span>.
              You will get your very own org and can create as many spaces as
              you like. {this.props.info.IgnitionOrgCount} people in your company are
              already using TAS!
            </p>
            {this.renderButton('Take me to my Org!', classes.speechButton)}
          </div>
        </div>
        <div className={classes.rocketMan} />
      </div>
    )
  }

  renderGettingStartedSteps () {
    const { classes } = this.props
    return (
      <div className={classes.ctaSteps}>
        <div className={classes.pewPew} />
        <div className={classes.step}>
          <div>
            <img className={classes.stepImage} src={step1} alt="step 1" />
          </div>
          <p>
            Get the<br />
            <a href="https://docs.pivotal.io/platform/cf-cli/">
              Cloud Foundry CLI
            </a>
            <br />
            from VMware
          </p>
        </div>
        <div className={classes.step}>
          <div>
            <img className={classes.stepImage} src={step2} alt="step 2" />
          </div>
          <p>
            Download the <br />
            <a href="https://github.com/cloudfoundry-samples/spring-music">
              sample app
            </a>
            <br />
            from Github
          </p>
        </div>
        <div className={classes.step}>
          <div>
            <img className={classes.stepImage} src={step3} alt="step 3" />
          </div>
          <p>
            Learn to<br />
            <a href="https://docs.pivotal.io/pivotalcf/devguide/deploy-apps/deploy-app.html">
              deploy an app
            </a>
          </p>
        </div>
      </div>
    )
  }

  renderSpacesInfo () {
    const { classes } = this.props
    return (
      <div className={classes.ctaSpaces}>
        <div>
          <div className={classes.spacesSpeech}>
            <p>
              <span className={classes.emphasis}>Spaces</span> can act like
              environments, and your first space is called{' '}
              {'"' + this.props.info.ExperimentationSpaceName + '"'}.
            </p>
            <p>
              Once apps are pushed to a space, you can bind them to{' '}
              <span className={classes.emphasis}>services</span> like MySQL and
              RabbitMQ by visiting the &quot;Marketplace&quot; link in TAS.
            </p>
            {this.renderButton(
              'I\'m ready. Go to my org!',
              classes.speechButton
            )}
          </div>
        </div>
        <div className={classes.moonMan} />
      </div>
    )
  }

  renderButton (text, extraClasses) {
    return (
      <Button
        size="large"
        variant="raised"
        className={classNames(this.props.classes.button, extraClasses)}
        onClick={this.handleOrgButtonClick}
      >
        {text}
      </Button>
    )
  }

  render () {
    const { classes } = this.props
    return (
      <div className={classes.body}>
        {this.renderWelcomeInfo()}
        {this.renderGettingStartedSteps()}
        {this.renderSpacesInfo()}
        <Footer links={footerLinks} logoURL={pivotalLogo} />
      </div>
    )
  }
}

Body.defaultProps = {
  info: {
    CompanyName: 'VMware',
    ExperimentationSpaceName: 'development',
    IgnitionOrgCount: 0
  }
}

Body.propTypes = {
  classes: PropTypes.object.isRequired,
  testing: PropTypes.bool,
  info: PropTypes.object
}

export default withStyles(styles)(Body)
