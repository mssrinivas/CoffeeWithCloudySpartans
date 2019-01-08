import React from 'react';
import {Redirect} from 'react-router-dom';
import Navigation from './Navigation';
import cookie from 'react-cookies'
import './StarterPage.css';
import {Link} from 'react-router-dom'

class StarterPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      gototlogin : null
    }

  }

  GotoLogin = () => 
  {
    this.setState({gototlogin : true})
  }
  
  render()
  {
    let Redirect_to_home = null;
    let gototloginbutt = null;

    if(this.state.gototlogin === true)
    {
      gototloginbutt = (<Redirect to="/login" />)
    }
    return (
        
      <div class="backgroundcontainer" >
      {Redirect_to_home}
      {gototloginbutt}
      <img class="bg" src='http://eskipaper.com/images/cool-coffee-wallpaper-1.jpg?h=630&la=en&w=900'/>
      <div class="centeringone row centeredone">

                              <div class="info-form">
                                  <form action="" class="form-inline justify-content-center">
                                      <div class="form-group">
                                      <button type="button" class="bluebutton btn btn-primary btn-lg roundy" onClick = {this.onSubmit}>Welcome to Coffee with Cloudy Spartans</button>
                                      </div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
                                      
                                      &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
                              
                                      <button type="button" class="bluebutton btn btn-primary btn-lg roundy" onClick = {this.GotoLogin}>Join Us</button>
                                      <br />  
                                      <br /> 
                                      <br />   
                                  </form>
                                  <h1 class="whitefont">This place serves love via coffee </h1> 
                              </div>
                        <br />  
           </div>
      </div>

    )
  }
}
export default StarterPage;

