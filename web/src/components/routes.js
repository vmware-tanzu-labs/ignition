import React from 'react'
import App from './app'
import Forbidden from './forbidden'
import NotFound from './notFound'
import { BrowserRouter, Switch, Route, Redirect } from 'react-router-dom'
import withRoot from '../withRoot'

const Routes = () => (
  <BrowserRouter>
    <Switch>
      <Route exact path="/" component={App} />
      <Route path="/404" component={NotFound} />
      <Route path="/forbidden" component={Forbidden} />
      <Redirect from="/*" to="/404" />
    </Switch>
  </BrowserRouter>
)

export default withRoot(Routes)
