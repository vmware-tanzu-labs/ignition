import React from 'react'
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles'
import teal from '@material-ui/core/colors/teal'
import blue from '@material-ui/core/colors/blue'
import CssBaseline from '@material-ui/core/CssBaseline'

const theme = createMuiTheme({
  palette: {
    primary: teal,
    secondary: blue,
    type: 'light'
  }
})

function withRoot (Component) {
  function WithRoot (props) {
    return (
      <MuiThemeProvider theme={theme}>
        <CssBaseline />
        <Component {...props} />
      </MuiThemeProvider>
    )
  }

  return WithRoot
}

export default withRoot
