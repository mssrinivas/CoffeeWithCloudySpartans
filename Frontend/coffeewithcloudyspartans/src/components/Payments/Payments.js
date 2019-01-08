import React from 'react';
import {Redirect} from 'react-router-dom';
import cookie from 'react-cookies';
import "../Payments/Payments.css"
import axios from 'axios'
import Navigation from '../StarterPage/Navigation';
var swal = require('sweetalert')
const URL ="http://localhost:4004"

class Payments extends React.Component {
    constructor(props){
        super(props);

        this.state={
            name : localStorage.getItem("usernamey"),  //this.props.data.user
            orderCount : 0,
            orderPrice : 0,
            orderId:null,
            amount:0,
            isBillGenerated:false,
            errors : "",
            orders:[]
        }
        this.getBill = this.getBill.bind(this)
        this.getOrderDetails=this.getOrderDetails.bind(this);
        this.amountHandler = this.amountHandler.bind(this);
        this.onSubmitHandler = this.onSubmitHandler.bind(this);
    }


    getOrderDetails(){

        if(!this.state.isBillGenerated){
            alert("You have to create payment to get order details");
            return;
        }

        const url = URL+ "/order/"+this.state.orderId
        axios.get(url).then(response => {
            console.log("RES",response.data);
            this.setState({orderPrice : response.data.amount});    
        }).catch(error=>{
                alert("Error in getting order details")
        })
    }

    componentDidMount()
    {
        const url = URL+ "/checkout"
        console.log("URL IS", url)
        axios.get(url,{
            params: {
                name:localStorage.getItem("usernamey")
              }}
            ).then(response => {
                localStorage.setItem("finalordercount", response.data.orderCount)
                this.setState({orderCount: response.data.orderCount})
            console.log("RES",response.data);
        }).catch(error=>{
                alert("Error in getting order details")
        })
    }

    getBill(e){
         const data ={
            name : this.state.name,
            count: this.state.orderCount
         }

         axios.post(URL+'/amount',data).then((response)=>{
            
            if(response.status === 200){
                console.log(response.status);
                console.log(response.data.id);
                alert("Your payment has been created")
                this.setState({isBillGenerated:true,orderId:response.data.id})
            }else{
                alert("Sorry we could not create your order")
            }
            
         }).catch((error)=>{
                console.log(error);
         })
      }

      amountHandler =(e)=>{
          console.log("Inside amount"+e.target.value)
          this.setState({amount:e.target.value});
          
      }

      onSubmitHandler= (e)=>{
        console.log("inside request submit handler");
          const data ={
                name : this.state.name,
                price:this.state.orderPrice
          }
          const x = this.state.amount 
          const y = this.state.orderPrice
          if(x==y)
          {
            console.log(this.state.amount + " ---" + this.state.orderPrice)
          const id = this.state.orderId;
          const url = URL+"/processOrders/"+id;
          console.log(url);
          axios.post(url,data).then((response)=>{
            console.log(response)
          if(response.status === 200){
              console.log("Successfully processes your order");
              swal("We have successfully processed your order", " ", "success");
              

              const url2 = URL+"/deletecart";
              console.log(url2);
              axios.delete(url2,
                {
                    params:
                    {
                        username : localStorage.getItem("usernamey")
                    }
                }).then((response)=>{
                console.log(response)
              if(response.status === 200){
                  console.log("Emptied Cart");
              }
              else{
                console.log("Emptied Cart Failed");
              }
            })
          }

          }).catch((error)=>{
              console.log(error);
              swal("Sorry we could not process your order", " ", "failure");
          })
        }
        else{
            swal("Sorry we could not process your order", " ", "failure");
            alert("Enter right amount")
        }
      }

  render(){ 
      const {isBillGenerated,orderPrice} = this.state;

      var USERTYPE = localStorage.getItem("ACCOUNTTYPE")
      var Redirecttologin = null;
      if(USERTYPE!="user"){
        Redirecttologin = (<Redirect to="/login"/>)
      }


    return(
        <div>
            <Navigation />
            {Redirecttologin}
            <div className="row justify-content-center" style={{marginTop:'10%'}}>

                <div className="col-md-6" style={{ border: '1px solid grey' }}>
                    <p>Customer Name : {this.state.name}</p>
                    <p>Order Count : {this.state.orderCount}</p>
                    {orderPrice!=null ? <p>Order Price : {this.state.orderPrice}</p> : ""}
                    <button type="button" onClick={this.getBill} className="btn btn-primary float-left mb-1">Create Payment</button>
                    <button type="button" onClick={this.getOrderDetails} className="btn btn-primary float-right mb-1">Order Details</button>
                </div>
            </div>
            <div className="row justify-content-center mt-3">

                <div className="col-md-6" style={{ border: '1px solid grey',display: isBillGenerated ? "" : "none" }}>
                    <form>


                        <div className="form-group">
                            <label htmlFor="cardnumber">Card No</label>
                            <input type="text" className="form-control" id="cardnumber" name="cardnumber" placeholder="Enter Card Number" />
                        </div>

                        <div className="form-group">
                            <label htmlFor="name">Name</label>
                            <input type="text" className="form-control" id="name" name="name" placeholder="Enter Your Name" />
                        </div>

                        <div className="form-group">
                            <label htmlFor="amount">Amount</label>
                            <input type="text" onChange={this.amountHandler} className="form-control" id="amount" name="amount"  placeholder="Enter Amount" value={this.state.amount} />
                        </div>

                        <div className="form-group">
                            <label htmlFor="expirydate">Expiry Date</label>
                            <input type="text" className="form-control" id="expirydate" name="expirydate" placeholder="MM/YY" />
                        </div>

                        <div className="form-group">
                            <label htmlFor="cvv">CVV</label>
                            <input type="password" className="form-control" id="cvv" name="cvv" placeholder="CVV" />
                        </div>

                        <button type="button" onClick={this.onSubmitHandler} className="btn btn-primary float-center mb-2" >Submit</button>
                    </form>
                    

                </div>

            </div>

        </div>
    );
}
}
export default Payments;
