import React from 'react'
import ReactDOM from 'react-dom'
import Home from './components/home'
import Forbidden from './components/forbidden'
import { BrowserRouter, Switch, Route } from 'react-router-dom'

ReactDOM.render((
  <BrowserRouter>
    <Switch>
      <Route exact path='/403' component={Forbidden} />
      <Route component={Home} />
    </Switch>
  </BrowserRouter>
), document.getElementById('root'))
