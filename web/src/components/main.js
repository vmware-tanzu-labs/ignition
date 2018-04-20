import React from 'react'
import ReactDOM from 'react-dom'
import Home from './home'
import Forbidden from './forbidden'
import { BrowserRouter, Switch, Route } from 'react-router-dom'
import withRoot from '../withRoot'

const Main = () => (
  <BrowserRouter>
    <Switch>
      <Route exact path='/403' component={Forbidden} />
      <Route component={Home} />
    </Switch>
  </BrowserRouter>
)

export default withRoot(Main)
