import React from 'react';
import { Container, Row, Col, Button, Fa, Card, CardBody, ModalFooter } from 'mdbreact';
import {Redirect} from 'react-router-dom';
import cookie from 'react-cookies';
import jwtDecode from 'jwt-decode';

const URL="http://localhost:4004"
class AdminLogin extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      signInEmail: '',
      signInPassword: '',
      Redirection_Value : false,
      errors : false
    }
  }

  onEmailChange = (event) => {
    this.setState({signInEmail: event.target.value})
  }

  onPasswordChange = (event) => {
    this.setState({signInPassword: event.target.value})
  }

  onSubmitSignIn = () => {
    fetch(URL + '/admin/login', {
      method: 'post',
      headers: {'Content-Type': 'application/json'},
      credentials : 'include',
      body: JSON.stringify({
      username: this.state.signInEmail,
      password: this.state.signInPassword
      })
    })
    .then(response => {
      if(response.status === 400)
        {
          this.setState({errors : true})
        }
      else
        {
          response.json()
          .then(user => {
         var decoded = jwtDecode(user);
              var accounttype = decoded.user.type;
              var userone = decoded.user._id;
              console.log("ACC - " + accounttype)
              console.log("USER", userone)
              localStorage.setItem("ACCOUNTTYPE", accounttype);
                console.log("NAME - " + userone)
                this.props.loadUser(userone);
                this.setState({Redirection_Value : true})
        })
        }
      })
  }
  
  render()
  {
    let Redirecty = null;
    let Errors = null;
    if(this.state.Redirection_Value === true || localStorage.getItem("ACCOUNTTYPE") === "admin")
    {
     Redirecty =  <Redirect to="/addadrink" />
    }
    if(this.state.errors === true)
    {
      Errors = <p class="error">Username or Password doesn't exist </p>
    }
    return(
      <div>
    {Redirecty}
    <br />
      <Container>
       <h1 class="makeitceneter"> Log in to Coffee with Cloudy Spartans</h1>
       <br />
       <br />
        <section className="form-elegant">
          <Row >
            <Col md="4" className="mx-auto">
              <Card>
                <CardBody className="mx-4">
                  <div className="text-center">
                    <h3 className="dark-grey-text mb-5">Admin Account Login</h3>
                    <hr></hr>
                  </div>
                  <input type="email" class="form-control" id="exampleInputEmail1" aria-describedby="emailHelp" placeholder="Username"  onChange={this.onEmailChange} required/>
                  <br>
                  </br>
                  <input type="password" class="form-control" id="exampleInputPassword1" placeholder="Password"  onChange={this.onPasswordChange} required/>
                   <br>
                  </br>
                  <div className="text-center mb-3">
                    <Button type="button" gradient="blue" className="btn btn-primary btn-lg btn-block" onClick = {this.onSubmitSignIn}>Log In</Button>
                      <hr></hr>
                  </div>
                </CardBody>     
              </Card>
            </Col>
          </Row>
        </section>
         {Errors}
      </Container>
      </div>
      );
  }
}
export default AdminLogin;