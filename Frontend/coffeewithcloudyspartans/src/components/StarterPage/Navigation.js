import React,{Component} from 'react';
import {Link} from 'react-router-dom';
import cookie from 'react-cookies';
import {Redirect} from 'react-router';
import './Navigation.css'

class Navigation extends Component {
    constructor(props){
        super(props);
        this.state = {
            LogOut : false
        }
    }
    handleLogout = () => {
        localStorage.clear()
        this.setState({LogOut : true})
    }
    render(){
        let navLogin = null;
        if(this.state.LogOut)
        {
            var LogOut = <Redirect to="/login" />
        }
        if(localStorage.getItem("ACCOUNTTYPE") === "user"){
           navLogin = (
               <div>
               <li class="nav-item dropdown ">
                <a class="nav-link dropdown-toggle lower " href="Dashboard" id="navbarDropdownMenuLink" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                    {localStorage.getItem("usernamey")}
                </a>
                <div class="dropdown-menu" aria-labelledby="navbarDropdownMenuLink">
                     <a class="dropdown-item" ><Link to="/mycart">My Cart</Link></a>
                     <a class="dropdown-item" ><Link to="/payment">Payments</Link></a>
                     <a class="dropdown-item" onClick={this.handleLogout}>LogOut</a>
                </div>
                </li>
                </div>
            );
        }else{
            //Else display login button
           navLogin = (
                   <li class="nav-item dropdown">
                    <a class="nav-link dropdown-toggle lower" href="http://example.com" id="navbarDropdownMenuLink" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                     <p class="backwhite"> Login </p>
                    </a>
                    <div class="dropdown-menu" aria-labelledby="navbarDropdownMenuLink">
                        <a class="dropdown-item" ><Link to="/login">Login</Link></a>
                        <a class="dropdown-item" ><Link to="/admin/login">Admin Login</Link></a>
                    </div>
                </li>
            )

        }// end of else
        
        return(
            <div > 
            {LogOut}        
            <nav class="navbar navbar-expand-lg">
                <a class="navbar-brand" href="/home"><img alt="HomeAway birdhouse" class="site-header-birdhouse__image " role="presentation" src="https://img.icons8.com/ios/50/000000/roman-soldier.png" height="50" width="50"/><p class="Design">&nbsp;&nbsp;CCS</p></a>
                <div className="spacebetweentwo"></div> 
                <a class="btn btn-outline-primary bgwhite roundy bg-fill" href="/home" role="button">Menu</a>
                <div className="spacebetweenone"></div> 
                <div className="spacebetweenone"></div> 
                <a class="btn btn-outline-primary bgwhite roundy bg-fill" href="/mycart" role="button">Cart</a> 
                <div className="spacebetweenone"></div> 
                <div className="spacebetweenone"></div> 
                <a class="btn btn-outline-primary bgwhite roundy bg-fill" href="/payment" role="button">Payments</a> 
                <div className="spacebetweenone"></div> 
                <div className="spacebetweenone"></div> 
                <div id="navbarNavDropdown" class="navbar-collapse collapse">
                    <ul class="navbar-nav mr-auto">
                    
                    </ul>
                    <ul class="navbar-nav">
                       {navLogin}
                        &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;
                        <li class="nav-item">
                        <a class="btn btn-outline-primary bgwhite roundy" href="/admin/login" role="button">Add a Drink</a> 
                        &nbsp; &nbsp; &nbsp; 
                        <a class="site-header-birdhouse lower" href="/home" title="Learn more"><img alt="" class="" role="presentation" src="https://img.icons8.com/doodle/50/000000/coffee-to-go.png"/></a>          
                        </li>
                    </ul>
                </div>
            </nav>
        </div>
        )
    }
}

export default Navigation;