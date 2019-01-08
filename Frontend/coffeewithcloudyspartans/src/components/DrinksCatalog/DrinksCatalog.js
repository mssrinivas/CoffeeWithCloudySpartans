import React from 'react';
 import './DrinksCatalog.css'
import {Redirect} from 'react-router';
import {Link} from 'react-router-dom';
import Drinks from '../Drinks/Drinks'
import Navigation from '../StarterPage/Navigation'
import { Container, Row, Col, Button, Fa, Card, CardBody, ModalFooter } from 'mdbreact';
var swal = require('sweetalert')

const URL="http://localhost:4004"
class DrinksCatalog extends React.Component {

 constructor(props) {
    super(props);
    this.state = {
      Drinks: [],
      drinkid: ""
    };
    this.handleClickPage = this.handleClickPage.bind(this);
    this.handleClickDrink = this.handleClickDrink.bind(this);
  }

  handleClickPage(event) {
    this.setState({
      currentPage: Number(event.target.id)
    });
  }

  handleClick(key){
    console.log("KEY IS " +key);
    localStorage.setItem("activekey" , key)
    this.setState(
      {propId:key})
    console.log(this.state)
}

handleClickDrink(key){
  console.log("Drink ID " + JSON.stringify(key));
  localStorage.setItem("DrinkItem" , key)
  fetch(URL+'/addtocart', {
    method: 'post',
    headers: {'Content-Type': 'application/json'},
    credentials : 'include',
    body: JSON.stringify({
        userid : localStorage.getItem("usernamey"),
        cartItems : {
          productid : key.id,
          name : key.name,
          price : key.price,
          size : key.size,
          count : 1
        }
    })
  })
  .then(response => {
    if(response.status === 400)
      {
        this.setState({errors : true})
        swal("Drink Cannot be Added!"," ", "failure");
      }
    else
      {
        swal("Drink Added successfully!", " ", "success");
     }
    })

}

  componentDidMount() {
    var result = []
    fetch(URL+'/getalldrinks')
    .then((response) => {
    response.json()
    .then(drinks => {
            this.setState({Drinks : drinks})
           })
          })
  }

  render() {

    let Redirect_to_login = null;
        let redirecty_value = null;
        var USERTYPE = localStorage.getItem("ACCOUNTTYPE")
            if(USERTYPE==="user"){
                  redirecty_value  = (
                    <div class="middle">
                    <div className="float-right">
                      </div>
                      <table class="tabledef">
                      <tbody>
                      {  
                        this.state.Drinks.map((drink, index) => {
                          console.log("TRIPS IS ", drink)
                            return ( 
                              <Drinks
                              key={index}
                              drinkinfo={drink}
                              id = {drink.productid}
                              name={drink.name}
                              sizes={drink.size}
                              price={drink.price}
                              clicked={this.handleClick}
                              drinkclicked= {this.handleClickDrink}
                              />  
                            );
                          })
             }
             </tbody>
             </table>
           </div>
         );
      }
      else{
        Redirect_to_login = <Redirect to="/login" />
      }

    return (
      <div> 
      {Redirect_to_login}
      <Navigation />
      <div class="divide">
      </div>
      <div id="bodydiv" >
      {redirecty_value}
      </div>
    </div>
    );
  }
}

export default DrinksCatalog;