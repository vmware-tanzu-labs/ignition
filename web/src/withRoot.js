import React from 'react'
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles'
import CssBaseline from '@material-ui/core/CssBaseline'

const theme = createMuiTheme({
  palette: {
    primary: {
      light: '#63aaf7',
      main: '#1f7bc4',
      dark: '#005093',
      contrastText: '#000000'
    },
    secondary: {
      light: '#9be965',
      main: '#68b634',
      dark: '#348600',
      contrastText: '#000000'
    }
  }
})

function withRoot(Component) {
  function WithRoot(props) {
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
