import React, { Component } from 'react';
import {Route} from 'react-router-dom'
import DrinksCatalog from './DrinksCatalog/DrinksCatalog';
import Login from './Login/Login';
import SignUp from './SignUp/SignUp';
import AddDrink from './AddDrink/AddDrink';
import CartCatalog from './Cart/CartCatalog';
import AdminLogin from './AdminLogin/AdminLogin';
import Payments from "../components/Payments/Payments";
import StarterPage from "./StarterPage/StarterPage";

class Main extends Component {
  constructor (props) {
  super(props)
  this.state ={
    user: {
    name: ''
  } 
  }
}

  componentDidMount()
  {
    var userabc = localStorage.getItem("usernamey")
    console.log("LOCAL STORAGE value is " + userabc)
    if(userabc) {
    this.setState({user: {
    name: userabc
    }});
  }
}

  loadUser = (data) => {
    console.log("DATA IS " + JSON.stringify(data));
    this.setState({user: {
    name: data
  }})

    localStorage.setItem("usernamey", data)
    console.log("SET localStorage ITEM AS " +  localStorage.getItem("usernamey"))

}
render() {
console.log("STATE IS  " + this.state.user.name);
  return (
        <div>
        <Route exact path="/" render={()=>(<StarterPage loadUser={this.loadUser} />)} />    
        <Route exact path="/home" render={()=>(<DrinksCatalog value={this.state.user.name} />)} />  
        <Route exact path="/login" render={()=>(<Login loadUser={this.loadUser} />)} />   
        <Route exact path="/admin/login" render={()=>(<AdminLogin loadUser={this.loadUser} />)} />   
        <Route exact path="/signup" render={()=>(<SignUp loadUser={this.loadUser} />)} />  
        <Route exact path="/addadrink" render={()=>(<AddDrink value={this.state.user.name} />)} /> 
        <Route path="/mycart" render={()=>(<CartCatalog value={this.state.user.name} />)} />  
        <Route exact path="/payment" render={()=>(<Payments value={this.state.user.name} />)} /> 
       </div>
  );
}
}

export default Main;



